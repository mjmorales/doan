package agent

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// Untar un-tars a tarball to a destination directory
func Untar(source, destination string) error {
	file, err := os.OpenFile(source, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	// Create the destination directory if it doesn't exist
	err = os.Mkdir(destination, 0755)
	// check if error is because directory already exists
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create destination directory: %s", err)
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %s", err)
	}

	// Open and iterate through the files in the archive.
	tarReader := tar.NewReader(bufio.NewReader(gzipReader))
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return fmt.Errorf("could not read next file in archive: %s", err)
		}

		path := fmt.Sprintf("%s/%s", destination, hdr.Name)

		// Create a directory if file is a directory
		if hdr.FileInfo().IsDir() {
			os.Mkdir(path, hdr.FileInfo().Mode())
			continue
		}

		// Log file being extracted if debug is enabled
		log.Debug().Msgf("Extracting %s", path)

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("could not create file %s: %s", hdr.Name, err)
		}

		defer outFile.Close()

		// Copy over contents
		if _, err := io.Copy(outFile, tarReader); err != nil {
			return fmt.Errorf("could not copy file %s: %s", hdr.Name, err)
		}
	}

	return nil
}

// Relink updates the symlinks to the active ansible repo
func Relink(stagingRepoPath string) error {
	// Remove the current active symlink
	err := os.Remove(DoanActiveDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not remove symlink: %s", err)
	}

	// Point the active symlink to the latest ansible repo
	err = os.Symlink(stagingRepoPath, DoanActiveDir)
	if err != nil {
		return fmt.Errorf("could not create symlink: %s", err)
	}

	return nil
}

// RemoveOldestStagingRepo removes the oldest staging repo
// if there are more than maxStagingRepos
func RemoveOldestStagingRepo(maxStagingRepos int) error {
	// Get the list of staging repos
	stagingRepos, err := os.ReadDir(DoanStagingDir)
	if err != nil {
		return fmt.Errorf("could not read staging directory: %s", err)
	}

	// Remove the oldest staging repo if there are more than maxStagingRepos
	if len(stagingRepos) > maxStagingRepos {
		oldestStagingRepo := stagingRepos[0]
		oldestStagingRepoPath := fmt.Sprintf("%s/%s", DoanStagingDir, oldestStagingRepo.Name())
		err = os.RemoveAll(oldestStagingRepoPath)
		if err != nil {
			return fmt.Errorf("could not remove oldest staging repo: %s", err)
		}

		// Recursively call RemoveOldestStagingRepo
		// until there are no more than maxStagingRepos
		RemoveOldestStagingRepo(maxStagingRepos)
	}

	return nil
}

func GetLocalMD5Sum(agentConfig AgentConfig) (string, error) {
	md5SumPath := fmt.Sprintf(
		"%s/%s/%s",
		DoanTarBallDir,
		agentConfig.AnsibleNameSpace,
		agentConfig.AnsibleTarballName,
	)

	f, err := os.Open(md5SumPath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %s", err)
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("could not copy file: %s", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// CompareMD5Sums checls if md5sum in artifactory matches
// the md5sum in the latest tarball
// CompareMD5Sums returns false if the md5sums do not match
// and returns an error if there is an issue getting the md5sums
func CompareMD5Sums(agentConfig AgentConfig) (bool, error) {
	remoteMD5Sum, err := GetRemoteMD5Sum(agentConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get remote md5sum: %s", err)
	}

	localMD5Sum, err := GetLocalMD5Sum(agentConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get local md5sum: %s", err)
	}

	log.Debug().Msgf("remote md5sum: %s ,local sum: %s", remoteMD5Sum, localMD5Sum)
	return remoteMD5Sum == localMD5Sum, nil
}

// Untar the latest ansible repo
// and updates symlinks to the active ansible repo.
// DeployRepo returns an error if the relinking fails.
func DeployRepo(agentConfig AgentConfig) error {
	checksumMatch, err := CompareMD5Sums(agentConfig)
	if err != nil {
		log.Error().Msgf("failed to compare md5sums: %s", err)
	}

	if checksumMatch {
		log.Info().Msgf("md5sums match, skipping deploy")
		return nil
	}

	// Download the latest ansible repo
	err = DownloadRepo(agentConfig)
	if err != nil {
		return fmt.Errorf("failed to download ansible repo: %s", err)
	}

	// Untar the latest ansible repo to a staging directory
	latestTarballPath := fmt.Sprintf(
		"%s/%s/%s",
		DoanTarBallDir,
		agentConfig.AnsibleNameSpace,
		agentConfig.AnsibleTarballName,
	)

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	stagingRepoPath := fmt.Sprintf("%s/%s", DoanStagingDir, timestamp)
	err = Untar(latestTarballPath, stagingRepoPath)
	if err != nil {
		return fmt.Errorf("failed to untar ansible repo: %s", err)
	}

	// Remove the oldest staging repo if there are more than maxStagingRepos
	err = RemoveOldestStagingRepo(agentConfig.MaxStagingRepos)
	if err != nil {
		return fmt.Errorf("failed to remove oldest staging repo: %s", err)
	}

	// Relink the active ansible repo with the latest staging repo
	err = Relink(stagingRepoPath)
	if err != nil {
		return fmt.Errorf("failed to relink ansible repo: %s", err)
	}

	return nil
}

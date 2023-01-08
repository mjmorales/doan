package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/rs/zerolog/log"

	aritfactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
)

type JFrogServers struct {
	AccessToken       string `json:"accessToken"`
	ArtifactoryURL    string `json:"artifactoryUrl"`
	DistributionURL   string `json:"distributionUrl"`
	IsDefault         bool   `json:"isDefault"`
	MissionControlURL string `json:"missionControlUrl"`
	PipelinesURL      string `json:"pipelinesUrl"`
	ServerID          string `json:"serverId"`
	URL               string `json:"url"`
	User              string `json:"user"`
	XrayURL           string `json:"xrayUrl"`
}

// JFrogCLIConfig is the JSON structure of the JFrog CLI config file
// usually located at ~/.jfrog/jfrog-cli.conf
type JFrogCLIConfig struct {
	Servers []JFrogServers `json:"servers"`
	Version string         `json:"version"`
}

// NewRtDetailsFromConfig returns auth.ServiceDetails using the JFrog CLI config.
// Requires the JFrog CLI to be installed and configured.
// Returns an error if the Artifactory details cannot be set.
func NewRtDetailsFromConfig(jFrogCLIConfigPath string) (auth.ServiceDetails, error) {
	rtDetails := aritfactoryAuth.NewArtifactoryDetails()
	// Read JFrog CLI Config JSON file from the JFrogCLIConfigPath parm
	content, err := os.ReadFile(jFrogCLIConfigPath)
	if err != nil {
		log.Error().Msgf("failed to read JFrogCLIConfig: %v", err)
		return rtDetails, err
	}

	// Unmarshal the JSON file into a JFrogCLIConfig struct
	var jfrogCLIConfig JFrogCLIConfig
	err = json.Unmarshal(content, &jfrogCLIConfig)
	if err != nil {
		log.Error().Msgf("failed to unmarshal JFrogCLIConfig: %v", err)
		return rtDetails, err
	}

	// Set the Artifactory details from the JFrogCLIConfig struct
	server := jfrogCLIConfig.Servers[0]
	rtDetails.SetUrl(server.ArtifactoryURL)
	rtDetails.SetUser(server.User)
	rtDetails.SetAccessToken(server.AccessToken)

	return rtDetails, nil
}

// CreateArtifactoryServicesManager retuns an ArtifactoryServicesManager struct
// that is used to interact with Artifactory. It returns an error if
// the ArtifactoryServicesManager cannot be created.
func CreateArtifactoryServicesManager(jfrogCliConfigPath string) (artifactory.ArtifactoryServicesManager, error) {
	var accessManager artifactory.ArtifactoryServicesManager
	rtDetails, err := NewRtDetailsFromConfig(jfrogCliConfigPath)
	if err != nil {
		log.Error().Msgf("failed to create Auth Details: %v", err)
		return accessManager, err
	}

	ctx := context.TODO()
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetContext(ctx).
		Build()

	if err != nil {
		log.Error().Msgf("failed to create Artifactory Service Config: %v", err)
		return accessManager, err
	}

	accessManager, err = artifactory.New(serviceConfig)
	if err != nil {
		log.Error().Msgf("failed to create Artifactory Services Manager: %v", err)
		return accessManager, err
	}

	return accessManager, nil
}

// DownloadRepo downloads the latest Ansible Repo tarball from Artifactory
// and returns an error if the download fails.
func DownloadRepo(agentConfig AgentConfig) error {
	rtManager, err := CreateArtifactoryServicesManager(agentConfig.JFrogCLIConfigPath)
	if err != nil {
		return err
	}

	// Download Ansible Tarball from Artifactory via JFrog CLI
	params := services.NewDownloadParams()
	params.Pattern = agentConfig.AnsibleRepoPath
	params.Target = fmt.Sprintf("%s/%s", DoanTarBallDir, agentConfig.AnsibleTarballName)

	totalDownloaded, totalFailed, err := rtManager.DownloadFiles(params)
	if err != nil {
		return err
	}

	if totalFailed > 0 {
		log.Error().Msgf("failed to download %d files", totalFailed)
		return nil
	}

	log.Info().Msgf("downloaded %d files", totalDownloaded)
	return nil
}

// GetRemoteMD5Sum checks the MD5 sum of the tarball in Artifactory
func GetRemoteMD5Sum(agentConfig AgentConfig) (string, error) {
	rtManager, err := CreateArtifactoryServicesManager(agentConfig.JFrogCLIConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to create Artifactory Services Manager: %v", err)
	}

	params := services.NewSearchParams()
	params.Pattern = agentConfig.AnsibleRepoPath

	reader, err := rtManager.SearchFiles(params)
	if err != nil {
		return "", err
	}

	defer reader.Close()

	err = reader.GetError()
	if err != nil {
		return "", err
	}

	var md5sum string
	for currentResult := new(utils.ResultItem); reader.NextRecord(currentResult) == nil; currentResult = new(utils.ResultItem) {
		log.Debug().Msgf("Found artifact: %s of type: %s\n md5:%s", currentResult.Name, currentResult.Type, currentResult.Actual_Md5)
		md5sum = currentResult.Actual_Md5
		break
	}

	return md5sum, nil
}

package agent

import (
	"fmt"
	"os"
)

const (
	// DoanWorkingDir is the working directory for the agent
	DoanWorkingDir = "/var/lib/doan"
	DoanTarBallDir = DoanWorkingDir + "/tarballs"
	DoanStagingDir = DoanWorkingDir + "/staging"
	DoanActiveDir  = DoanWorkingDir + "/active"
)

// Init creates the directories needed for the agent
func initDir() error {
	doanDirectories := []string{
		DoanWorkingDir,
		DoanTarBallDir,
		DoanStagingDir,
	}

	for _, dir := range doanDirectories {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("could not create directory %s: %s", dir, err)
		}
	}

	return nil
}

func Init(agentConfig AgentConfig) error {
	err := initDir()
	if err != nil {
		return err
	}

	err = DeployRepo(agentConfig)
	if err != nil {
		return err
	}

	err = RunActiveAnsiblePlaybook(agentConfig)
	if err != nil {
		return err
	}

	return nil
}

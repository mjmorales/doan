package agent

import (
	"os"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

// CLIFlags are the command line flags for the agent
type CLIFlags struct {
	UpdateAnsibleRepo  bool
	JFrogCLIConfigPath string
	AnsibleRepoPath    string
	Init               bool
	MaxStagingRepos    int
	AnsibleTarballName string
	AnsibleNameSpace   string
	Daemon             bool
	DaemonInterval     string
	ConfigFilePath     string
	Version            bool
}

// AgentConfig is the configuration for the agent
type AgentConfig struct {
	JFrogCLIConfigPath string `yaml:"jfrog_cli_config_path"`
	AnsibleRepoPath    string `yaml:"ansible_repo_path"`
	MaxStagingRepos    int    `yaml:"max_staging_repos"`
	AnsibleTarballName string `yaml:"ansible_tarball_name"`
	AnsibleNameSpace   string `yaml:"ansible_namespace"`
	DaemonInterval     string `yaml:"daemon_interval"`
	Daemon             bool   `yaml:"daemon"`
}

func (c *AgentConfig) WithConfigFromFile(configFilePath string) *AgentConfig {
	var ac AgentConfig

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Error().Msgf("failed to open config file: %s", err)
		return c
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&ac)
	if err != nil {
		log.Error().Msgf("failed to decode config file: %s", err)
	}

	return &ac
}

// NewConfig returns a new Config struct from the command line flags
func (c *AgentConfig) WithConfigFromCLI(agentFlags CLIFlags) *AgentConfig {
	c.JFrogCLIConfigPath = agentFlags.JFrogCLIConfigPath
	c.AnsibleRepoPath = agentFlags.AnsibleRepoPath
	c.MaxStagingRepos = agentFlags.MaxStagingRepos
	c.AnsibleTarballName = agentFlags.AnsibleTarballName
	c.AnsibleNameSpace = agentFlags.AnsibleNameSpace
	c.DaemonInterval = agentFlags.DaemonInterval
	c.Daemon = agentFlags.Daemon
	return c
}

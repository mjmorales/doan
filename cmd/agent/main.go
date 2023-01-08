package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/mjmorales/doan/pkg/agent"
	logger "github.com/mjmorales/doan/pkg/logger"
)

// ParseFlags parses the command line flags
func ParseFlags() agent.CLIFlags {
	updateAnsibleRepoPtr := flag.Bool("update-ansible-repo", false, "update the ansible repo and Exit (default: false)")
	jfrogCLIConfigPathPtr := flag.String("jfrog-cli-config-path", "$HOME/.jfrog/jfrog-cli.conf", "path to the JFrog CLI config file (default: $HOME/.jfrog/jfrog-cli.conf)")
	ansibleRepoPath := flag.String("ansible-repo-path", "generic-repo/path/to/tar", "path to the ansible repo (default: generic-repo/path/to/tar)")
	initPtr := flag.Bool("init", false, "initialize the agent (default: false)")
	maxStagingReposPtr := flag.Int("max-staging-repos", 10, "maximum number of staging repos to keep (default: 10)")
	ansibleTarballNamePtr := flag.String("ansible-tarball-name", "ansible.tar.gz", "name of the ansible tarball (default: ansible.tar.gz)")
	ansibleNameSpacePtr := flag.String("ansible-namespace", "ansible", "name of the ansible namespace (default: ansible)")
	daemonPtr := flag.Bool("daemon", false, "run binary in daemon mode (default: false)")
	daemonIntervalPtr := flag.String("daemon-interval", "1m", "interval string to run the daemon (default: 1m)")
	configFilePathPtr := flag.String("config-file-path", "", "path to the config file (default: \"\")")
	versionPtr := flag.Bool("version", false, "print the version and exit (default: false)")

	flag.Parse()
	agentFlags := agent.CLIFlags{
		UpdateAnsibleRepo:  *updateAnsibleRepoPtr,
		JFrogCLIConfigPath: *jfrogCLIConfigPathPtr,
		AnsibleRepoPath:    *ansibleRepoPath,
		Init:               *initPtr,
		MaxStagingRepos:    *maxStagingReposPtr,
		AnsibleTarballName: *ansibleTarballNamePtr,
		AnsibleNameSpace:   *ansibleNameSpacePtr,
		Daemon:             *daemonPtr,
		DaemonInterval:     *daemonIntervalPtr,
		ConfigFilePath:     *configFilePathPtr,
		Version:            *versionPtr,
	}

	return agentFlags
}

func init() {
	// set the global log configuration
	logger.SetGlobalLogConfig()
}

func main() {
	// parse the command line flags
	agentFlags := ParseFlags()
	agentConfig := &agent.AgentConfig{}

	if agentFlags.Version {
		fmt.Print(agent.GetVersion())
		os.Exit(0)
	}

	if agentFlags.ConfigFilePath != "" {
		log.Info().Msg("loading agent config from file")
		agentConfig = agentConfig.WithConfigFromFile(agentFlags.ConfigFilePath)
	} else {
		log.Info().Msg("loading agent config from CLI")
		agentConfig = agentConfig.WithConfigFromCLI(agentFlags)
	}

	// initialize the agent if the init flag is set
	if agentFlags.Init {
		log.Info().Msg("initializing agent")
		err := agent.Init(*agentConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize agent")
		}

		os.Exit(0)
	}

	// update the ansible repo if the update-ansible-repo flag is set
	if agentFlags.UpdateAnsibleRepo {
		err := agent.DeployRepo(*agentConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to update ansible repo")
		}

		os.Exit(0)
	}

	// start the agent if the agent flag is set
	if agentConfig.Daemon {
		log.Info().Msg("starting agent in daemon mode")
		agent.Daemon(*agentConfig)
	}
}

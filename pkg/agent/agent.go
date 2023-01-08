package agent

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// Runner is a struct that holds a mutex lock
// to prevent multiple ansible runs at the same time
type Runner struct {
	mutex sync.Mutex
}

// Run runs the DeployRepo and RunActiveAnsiblePlaybook functions
// with a mutex lock to prevent multiple runs at the same time
// Run returns early if the runner is already running
func (r *Runner) Run(agentConfig AgentConfig) {
	locked := r.mutex.TryLock()
	if locked {
		log.Error().Msg("this runner is already running")
		return
	}

	defer r.mutex.Unlock()

	err := DeployRepo(agentConfig)
	if err != nil {
		log.Error().Msgf("failed to deploy repo: %s", err)
	}

	err = RunActiveAnsiblePlaybook(agentConfig)
	if err != nil {
		log.Error().Msgf("failed to run ansible: %s", err)
	}
}

// RunActiveAnsiblePlaybook runs localhost inventory
// on the base.yaml playbook in the active ansible repo
// RunActiveAnsiblePlaybook returns an error if the ansible run fails
func RunActiveAnsiblePlaybook(agentConfig AgentConfig) error {
	ansiblePlayBookCommand := "ansible-playbook"
	ansiblePlayBookCommandParams := []string{
		"/var/lib/doan/active/ansible/base.yaml",
		"-i",
		"/var/lib/doan/active/ansible/inventory.yaml",
		"--vault-password-file",
		"~/.vault_pass.txt",
	}

	cmd := exec.Command(
		ansiblePlayBookCommand,
		ansiblePlayBookCommandParams...,
	)
	cmd.Stdout = log.Logger
	cmd.Stderr = log.Logger

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run ansible: %s", err)
	}

	log.Info().Msg("ansible run complete")
	return nil
}

// Daemon creates the Runner struct and starts the Run function
// in a scheduled interval
func Daemon(agentConfig AgentConfig) {
	runner := &Runner{}

	s := gocron.NewScheduler(time.UTC)
	s.Every(agentConfig.DaemonInterval).Do(func() { runner.Run(agentConfig) })
	s.StartBlocking()
}

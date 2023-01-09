package agent

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
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

func GetDropletTags() ([]string, error) {
	dropletTags := []string{}

	// get the droplet tags through the DO HTTP API
	doTagsUrl := "http://169.254.169.254/metadata/v1/tags"
	resp, err := http.Get(doTagsUrl)
	if err != nil {
		return dropletTags, fmt.Errorf("could not get droplet tags: %s", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dropletTags, fmt.Errorf("could not read droplet tags: %b - %s", resp.StatusCode, err)
	}

	playbookTags := []string{"base"}
	dropletTags = strings.Split(string(body), "\n")

	for _, tag := range dropletTags {
		if strings.Contains(tag, "ansible-") {
			playbookTags = append(playbookTags, tag)
		}
	}

	return playbookTags, nil
}

// RunActiveAnsiblePlaybook runs localhost inventory
// on the base.yaml playbook in the active ansible repo
// RunActiveAnsiblePlaybook returns an error if the ansible run fails
func RunActiveAnsiblePlaybook(agentConfig AgentConfig) error {
	tags, err := GetDropletTags()
	if err != nil {
		return err
	}

	ansiblePlayBookCommand := "ansible-playbook"
	ansiblePlayBookCommandParams := []string{
		"/var/lib/doan/active/ansible/base.yaml",
		"-i",
		"/var/lib/doan/active/ansible/inventory.yaml",
		"--vault-password-file",
		"~/.vault_pass.txt",
		"--tags",
		strings.Join(tags, ","),
	}

	cmd := exec.Command(
		ansiblePlayBookCommand,
		ansiblePlayBookCommandParams...,
	)
	cmd.Stdout = log.Logger
	cmd.Stderr = log.Logger

	err = cmd.Run()
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

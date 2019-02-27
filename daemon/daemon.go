package daemon

import (
	"log"

	"github.com/takama/daemon"
)

const (

	// name of the service
	name = "github-desktop-notifications"
)

// Description of the daemon
const Description = `Integrates with GitHub Notifications API to provide Desktop Notifications`

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon daemon.Daemon
}

// Install the daemon
func (service *Service) Install(configFile string, token string) (string, error) {
	args := []string{"daemon"}
	if token != "" {
		args = append(args, "--token", token)
	}

	if configFile != "" {
		args = append(args, "--config", configFile)
	}

	return service.daemon.Install(args...)
}

// Start the daemon
func (service *Service) Start() (string, error) {
	return service.daemon.Start()
}

// Stop the daemon
func (service *Service) Stop() (string, error) {
	return service.daemon.Stop()
}

// Status - check the service status
func (service *Service) Status() (string, error) {
	return service.daemon.Status()
}

// New creates a new Service to interact
func New() (*Service, error) {
	daemon, err := daemon.New(name, Description)
	if err != nil {
		return nil, err
	}

	s := &Service{
		daemon: daemon,
	}
	return s, nil
}

package daemon

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

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

// Manage by daemon commands or run the daemon
func (service *Service) Manage(token string) error {

	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	errorc := make(chan error, 1)
	go lookForNotifications(token, stdlog, errorc)

	select {
	case err := <-errorc:
		return err
	case killSignal := <-interrupt:

		stdlog.Println("Got signal:", killSignal)
		if killSignal == os.Interrupt {
			return errors.New("Daemon was interruped by system signal")
		}
		return errors.New("Daemon was killed")
	}
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

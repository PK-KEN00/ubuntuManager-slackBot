package systemctl

import (
	"../lib"
	"fmt"
	"os/exec"
	"strings"
)

func RunServiceStatus(service string) string {

	if !isValidService(service) {
		return fmt.Sprintf("Invalid service name. Available services are: %s", strings.Join(lib.ServiceNames, ", "))
	}

	cmd := exec.Command("systemctl", "status", service)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Failed to get service status: %v", err)
	}

	return fmt.Sprintf("Service status for %s:\n```%s```", service, string(output))
}

func RunServiceRestart(service string) string {

	if !isValidService(service) {
		return fmt.Sprintf("Invalid service name. Available services are: %s", strings.Join(lib.ServiceNames, ", "))
	}

	cmd := exec.Command("systemctl", "restart", service)
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Failed to restart service: %v", err)
	}

	return fmt.Sprintf("Service %s restarted successfully.", service)
}

func isValidService(service string) bool {
	for _, s := range lib.ServiceNames {
		if s == service {
			return true
		}
	}
	return false
}

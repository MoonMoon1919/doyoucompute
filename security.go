package doyoucompute

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	ErrShellNotAllowed   = errors.New("shell not allowed")
	ErrCommandNotAllowed = errors.New("command not allowed")
	ErrDangerousCommand  = errors.New("dangerous command blocked")
	ErrInvalidWorkingDir = errors.New("invalid working directory")
)

var dangerousCommands = []string{
	"format", "fdisk", "mkfs", // Disk formatting
	"shutdown", "reboot", "halt", // System control
	"sudo", "su", // Privilege escalation
	"iptables", "ufw", "firewall-cmd", // Firewall changes
	"crontab", "at", "batch", "atq", "atrm", // Task scheduling
	"systemctl", "service", "launchctl", // Do allow starting/stopping services
	"dd", // prevent people from creating huge files
}

var dangerousPatterns = []string{
	"rm -rf /", "rm -fr /", // Root deletion
	"rm -rf /*", "rm -fr /*", // Root contents
	"> /dev/sd", "> /dev/hd", "> /dev/nvme", // Device writing
	"dd of=/dev/",                   // DD to devices
	":(){ :|:& };:",                 // Fork bomb
	"chmod 777 /", "chmod -R 777 /", // Dangerous permissions on root
}

// ValidateCommandPlan validates that a command plan is safe to execute
func ValidateCommandPlan(plan CommandPlan, config ExecutionConfig) error {
	if err := validatePlanArgs(plan, config); err != nil {
		return err
	}

	if err := validateShell(plan.Shell, config); err != nil {
		return err
	}

	return nil
}

// validatePlanArgs checks basic plan structure and all command-related validation
func validatePlanArgs(plan CommandPlan, config ExecutionConfig) error {
	if len(plan.Args) == 0 {
		return errors.New("command plan has no arguments")
	}

	if strings.TrimSpace(plan.Args[0]) == "" {
		return errors.New("command cannot be empty")
	}

	baseCommand := filepath.Base(plan.Args[0])

	// Check command allow-list if configured
	if len(config.AllowedCommands) > 0 {
		allowed := false
		for _, cmd := range config.AllowedCommands {
			if baseCommand == cmd {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("%w: %s (allowed: %v)", ErrCommandNotAllowed, baseCommand, config.AllowedCommands)
		}
	}

	// Check for dangerous commands and patterns if enabled
	if config.BlockDangerousCommands {
		for _, dangerous := range dangerousCommands {
			if baseCommand == dangerous {
				return fmt.Errorf("%w: %s", ErrDangerousCommand, dangerous)
			}
		}

		fullCommand := strings.Join(plan.Args, " ")
		for _, pattern := range dangerousPatterns {
			if strings.Contains(fullCommand, pattern) {
				return fmt.Errorf("%w: contains '%s'", ErrDangerousCommand, pattern)
			}
		}
	}

	return nil
}

// validateShell checks if the shell is allowed
func validateShell(shell string, config ExecutionConfig) error {
	if len(config.AllowedShells) == 0 {
		return nil // No restrictions
	}

	for _, allowed := range config.AllowedShells {
		if shell == allowed {
			return nil
		}
	}

	return fmt.Errorf("%w: %s (allowed: %v)", ErrShellNotAllowed, shell, config.AllowedShells)
}

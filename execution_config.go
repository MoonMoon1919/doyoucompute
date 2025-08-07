package doyoucompute

import "time"

type ExecutionConfig struct {
	// Timeout for command execution (0 means no timeout)
	Timeout time.Duration

	// AllowedShells restricts which shells/interpreters can be used
	AllowedShells []string

	// AllowedCommands restricts which commands can be executed (nil means allow all)
	AllowedCommands []string

	// BlockDangerousCommands prevents obviously dangerous operations
	BlockDangerousCommands bool
}

func DefaultSecureConfig() ExecutionConfig {
	return ExecutionConfig{
		Timeout:                30 * time.Second,
		AllowedShells:          []string{"bash", "sh", "python3", "python", "node", "go"},
		BlockDangerousCommands: true,
	}
}

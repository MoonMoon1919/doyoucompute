package doyoucompute

import "testing"

func TestValidateCommandPlan(t *testing.T) {
	tests := []struct {
		name         string
		plan         CommandPlan
		config       ExecutionConfig
		errorMessage string
	}{
		{
			name: "Passing-AllowAll",
			config: ExecutionConfig{
				BlockDangerousCommands: false,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"echo", "hello", "world"},
			},
			errorMessage: "",
		},
		{
			name: "Passing-InAllowedCommandsAndShell",
			config: ExecutionConfig{
				AllowedShells:          []string{"sh"},
				AllowedCommands:        []string{"echo"},
				BlockDangerousCommands: false,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"echo", "hello", "world"},
			},
			errorMessage: "",
		},
		{
			name: "Fail-DisallowedShell",
			config: ExecutionConfig{
				AllowedShells:          []string{"sh"},
				AllowedCommands:        []string{"echo"},
				BlockDangerousCommands: false,
			},
			plan: CommandPlan{
				Shell: "bash",
				Args:  []string{"echo", "hello", "world"},
			},
			errorMessage: "shell not allowed: bash (allowed: [sh])",
		},
		{
			name: "Fail-EmptyArgs",
			config: ExecutionConfig{
				BlockDangerousCommands: false,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{},
			},
			errorMessage: "command plan has no arguments",
		},
		{
			name: "Fail-EmptyCommand",
			config: ExecutionConfig{
				BlockDangerousCommands: false,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"", "hello", "world"},
			},
			errorMessage: "command cannot be empty",
		},
		{
			name: "Fail-DangerousCommand",
			config: ExecutionConfig{
				BlockDangerousCommands: true,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"sudo", "su"},
			},
			errorMessage: "dangerous command blocked: sudo",
		},
		{
			name: "Fail-DangerousPattern",
			config: ExecutionConfig{
				BlockDangerousCommands: true,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"chmod", "777", "/"},
			},
			errorMessage: "dangerous command blocked: contains 'chmod 777 /'",
		},
		{
			name: "Fail-NotInAllowedCommands",
			config: ExecutionConfig{
				AllowedCommands:        []string{"echo"},
				BlockDangerousCommands: true,
			},
			plan: CommandPlan{
				Shell: "sh",
				Args:  []string{"chmod", "+x", "script.sh"},
			},
			errorMessage: "command not allowed: chmod (allowed: [echo])",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCommandPlan(tc.plan, tc.config)

			var errMessage string
			if err != nil {
				errMessage = err.Error()
			}

			if errMessage != tc.errorMessage {
				t.Errorf("Got error %s, expected %s", errMessage, tc.errorMessage)
			}
		})
	}
}

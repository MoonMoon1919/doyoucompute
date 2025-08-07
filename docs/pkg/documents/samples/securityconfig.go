package samples

import (
	"fmt"
	"time"

	"github.com/MoonMoon1919/doyoucompute"
)

func securityconfig() {
	// Default secure configuration
	config := doyoucompute.DefaultSecureConfig()

	// or, a custom configuration!
	config = doyoucompute.ExecutionConfig{
		Timeout:                30 * time.Second,
		AllowedShells:          []string{"bash", "python3"},
		BlockDangerousCommands: true,
	}

	runner := doyoucompute.NewTaskRunner(config)
	repo := doyoucompute.NewFileRepository()
	markdownRenderer := doyoucompute.NewMarkdownRenderer()
	execRenderer := doyoucompute.NewExecutionRenderer()
	service := doyoucompute.NewService(repo, runner, markdownRenderer, execRenderer)

	// do something with service, probably not print!
	fmt.Printf("service: %v\n", service)
}

# DOYOUCOMPUTE

A lightweight framework for creating runnable documentation. Write your documentation once, then render it as markdown or execute it as a script. Ideal for tutorials, setup guides, and operational runbooks that need to stay up-to-date.

## Features

- 📝 Write documentation using a fluent, type-safe Go API
- 🚀 Execute embedded commands to validate your docs stay current
- 📋 Generate clean markdown output for GitHub, GitLab, etc.
- 🔧 Compare generated docs with existing files for CI/CD validation
- ⚡ Section-based execution for targeted testing


## Quick Start

### Installation

```bash
go get github.com/MoonMoon1919/doyoucompute
```

### Basic Usage

Create a simple document with executable commands:

```go
package main

import "github.com/MoonMoon1919/doyoucompute"

func main() {
    doc, err := doyoucompute.NewDocument("My Project")
    if err != nil {
        return err
    }

    // Add an introduction
    doc.WriteIntro().
        Text("Welcome to my project! ").
        Text("Follow these steps to get started.")

    // Add a setup section with executable commands
    setup := doc.NewSection("Setup")
    setup.WriteParagraph().
        Text("First, install dependencies:")

    setup.WriteCodeBlock("bash", []string{"npm install"}, doyoucompute.Exec)

    setup.WriteParagraph().
        Text("Then start the development server:")

    setup.WriteCodeBlock("bash", []string{"npm run dev"}, doyoucompute.Exec)
}
```

### CLI Usage

Create a CLI wrapper for your documents:

```go
package main

import (
    "os"
    "github.com/MoonMoon1919/doyoucompute"
    "github.com/MoonMoon1919/doyoucompute/app"
)

func createReadmeDoc() doyoucompute.Document {
	doc := doyoucompute.NewDocument("My Project")

	// Add your content

	return doc
}

func main() {
    // Create service with file repository and task runner
    repo := doyoucompute.NewFileRepository()
    runner := doyoucompute.NewTaskRunner(doyoucompute.DefaultSecureConfig())
    markdownRenderer := doyoucompute.NewMarkdownRenderer()
    executionRenderer := doyoucompute.NewExecutionRenderer()

    service := doyoucompute.NewService(repo, runner, markdownRenderer, executionRenderer)

    // Create and run CLI app
    app := app.New(service)
    app.Register(createReadmeDoc())

    if err := app.Run(os.Args); err != nil {
        panic(err)
    }
}
```

#### Available Commands

| Command | Description | Example |
| ---- | ---- | ---- |
| render | Generate markdown from document | ./cli render --doc-name=readme --path=README.md |
| compare | Compare document with existing file | ./cli compare --doc-name=readme --path=README.md |
| run | Execute all commands in document | ./cli run --doc-name=setup |
| plan | Show execution plan without running | ./cli plan --doc-name=setup --section="Database Setup" |
| list | List all available documents | ./cli list |

## Security Features

DOYOUCOMPUTE includes built-in security features to prevent dangerous command execution:

- 🛡️ Dangerous command blocking (rm -rf, sudo, etc.)
- ⏱️ Configurable execution timeouts
- 🐚 Shell allow-listing
- 🔒 Command validation and sanitization
- 🌍 Environment variable validation


### Configuration

Customize execution behavior with security configurations:

```go
// Default secure configuration
config := doyoucompute.DefaultSecureConfig()

// Custom configuration
config := doyoucompute.ExecutionConfig{
	Timeout: 30 * time.Second,
	AllowedShells: []string{"bash", "python3"},
	BlockDangerousCommands: true,
}

runner, err := doyoucompute.NewTaskRunner(config)
service := doyoucompute.NewService(repo, runner, markdownRenderer, executionRenderer)
```

## Environment Variables

Commands can specify required environment variables:

```go
// Command that requires API_KEY to be set
setup := doyoucompute.NewSection("Setup")

setup.WriteExecutable(
    "bash",
	[]string{"curl", "-H", "Authorization: Bearer $API_KEY", "api.example.com"},
	[]string{"API_KEY"})
```

## Recommendations

- 🔄 Run 'compare' in CI to ensure docs stay current
- 🧪 Use 'plan' to preview commands before execution
- 📂 Organize related commands into logical sections


## Contributing

See [CONTRIBUTING](./CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Disclaimers

This work does not represent the interests or technologies of any employer, past or present. It is a personal project only.

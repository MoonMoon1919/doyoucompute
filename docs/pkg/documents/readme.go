package documents

import "github.com/MoonMoon1919/doyoucompute"

func cliSection() doyoucompute.Section {
	cliSection := doyoucompute.NewSection("CLI Usage")

	cliSection.WriteIntro().
		Text("Create a CLI wrapper for your documents:")

	cliSection.WriteCodeBlock("go", []string{`package main

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
    // Define your documents
    documents := []*doyoucompute.Document{
        createReadmeDoc(),
    }

    // Create service with file repository and task runner
    repo := doyoucompute.NewFileRepository()
    runner := doyoucompute.NewTaskRunner()
    markdownRenderer := doyoucompute.NewMarkdownRenderer()
    executionRenderer := doyoucompute.NewExecutionRenderer()

    service := doyoucompute.NewService(repo, runner, markdownRenderer, executionRenderer)

    // Create and run CLI app
    app := app.New(documents, service)
    if err := app.Run(os.Args); err != nil {
        panic(err)
    }
}`}, doyoucompute.Static)

	availableCommandsSection := cliSection.CreateSection("Available Commands")
	commandsTable := availableCommandsSection.CreateTable([]string{"Command", "Description", "Example"})

	commandsTable.AddRow(
		"render",
		"Generate markdown from document",
		"./cli render --doc-name=readme --path=README.md",
	)
	commandsTable.AddRow(
		"compare",
		"Compare document with existing file",
		"./cli compare --doc-name=readme --path=README.md",
	)
	commandsTable.AddRow(
		"run",
		"Execute all commands in document",
		"./cli run --doc-name=setup",
	)
	commandsTable.AddRow(
		"plan",
		"Show execution plan without running",
		"./cli plan --doc-name=setup --section=\"Database Setup\"",
	)
	commandsTable.AddRow(
		"list",
		"List all available documents",
		"./cli list",
	)

	return cliSection
}

func basicUsageSection() doyoucompute.Section {
	basicUsageSection := doyoucompute.NewSection("Basic Usage")
	basicUsageSection.WriteIntro().
		Text("Create a simple document with executable commands:")

	basicUsageSection.WriteCodeBlock("go", []string{`package main

import "github.com/MoonMoon1919/doyoucompute"

func main() {
    doc := doyoucompute.NewDocument("My Project")

    // Add an introduction
    doc.WriteIntro().
        Text("Welcome to my project! ").
        Text("Follow these steps to get started.")

    // Add a setup section with executable commands
    setup := doc.NewSection("Setup")
    setup.NewParagraph().
        Text("First, install dependencies:")

    setup.WriteCodeBlock("bash", []string{"npm install"}, doyoucompute.Exec)

    setup.NewParagraph().
        Text("Then start the development server:")

    setup.WriteCodeBlock("bash", []string{"npm run dev"}, doyoucompute.Exec)
}`}, doyoucompute.Static)

	return basicUsageSection
}

func quickstartSection() doyoucompute.Section {
	quickStartSection := doyoucompute.NewSection("Quick Start")

	installationSection := quickStartSection.CreateSection("Installation")
	installationSection.WriteCodeBlock("bash", []string{"go get github.com/MoonMoon1919/doyoucompute"}, doyoucompute.Static)
	quickStartSection.AddSection(basicUsageSection())
	quickStartSection.AddSection(cliSection())

	return quickStartSection
}

func Readme() doyoucompute.Document {
	document := doyoucompute.NewDocument("DOYOUCOMPUTE")
	document.WriteIntro().
		Text("A lightweight framework for creating runnable documentation.").
		Text("Write your documentation once, then render it as markdown or execute it as a script.").
		Text("Ideal for tutorials, setup guides, and operational runbooks that need to stay up-to-date.")

	// Features
	featuresSection := document.CreateSection("Features")
	featureList := featuresSection.CreateList(doyoucompute.BULLET)
	featureList.Append("📝 Write documentation using a fluent, type-safe Go API")
	featureList.Append("🚀 Execute embedded commands to validate your docs stay current")
	featureList.Append("📋 Generate clean markdown output for GitHub, GitLab, etc.")
	featureList.Append("🔧 Compare generated docs with existing files for CI/CD validation")
	featureList.Append("⚡ Section-based execution for targeted testing")

	// Quick Start
	document.AddSection(quickstartSection())

	// Contributing
	contributing := document.CreateSection("Contributing")
	contributing.WriteIntro().
		Text("See").
		Link("CONTRIBUTING", "./CONTRIBUTING.md").
		Text("for details.")

	// License section
	licenseSection := document.CreateSection("License")
	licenseSection.WriteIntro().
		Text("MIT License - see").
		Link("LICENSE", "./LICENSE").
		Text("for details.")

	return document
}

package documents

import (
	"os"

	"github.com/MoonMoon1919/doyoucompute"
)

func recommendations() doyoucompute.Section {
	recommendationsSection := doyoucompute.NewSection("Recommendations")
	practicesList := recommendationsSection.CreateList(doyoucompute.BULLET)
	practicesList.Append("🔄 Run 'compare' in CI to ensure docs stay current")
	practicesList.Append("🧪 Use 'plan' to preview commands before execution")
	practicesList.Append("📂 Organize related commands into logical sections")

	return recommendationsSection
}

func environmentVariables() (doyoucompute.Section, error) {
	envSection := doyoucompute.NewSection("Environment Variables")
	envSection.WriteIntro().
		Text("Commands can specify required environment variables:")

	sample, err := os.ReadFile("./docs/pkg/documents/samples/envvars.go")
	if err != nil {
		return doyoucompute.Section{}, err
	}

	envSection.WriteCodeBlock("go", []string{string(sample)}, doyoucompute.Static)

	envSection.WriteParagraph().
		Text("The command will fail to run if the required environment variables are not set and report which are missing.")

	return envSection, nil
}

func configurationSecurity() (doyoucompute.Section, error) {
	configSection := doyoucompute.NewSection("Configuration")

	configSection.WriteIntro().
		Text("Customize execution behavior with security configurations:")

	sample, err := os.ReadFile("./docs/pkg/documents/samples/securityconfig.go")
	if err != nil {
		return doyoucompute.Section{}, err
	}

	configSection.WriteCodeBlock("go", []string{string(sample)}, doyoucompute.Static)

	return configSection, nil
}

func securitySection() doyoucompute.Section {
	securitySection := doyoucompute.NewSection("Security Features")

	securitySection.WriteIntro().
		Text("DOYOUCOMPUTE includes built-in security features to prevent dangerous command execution:")

	securityList := securitySection.CreateList(doyoucompute.BULLET)
	securityList.Append("🛡️ Dangerous command blocking (rm -rf, sudo, etc.)")
	securityList.Append("⏱️ Configurable execution timeouts")
	securityList.Append("🐚 Shell allow-listing")
	securityList.Append("🔒 Command validation and sanitization")
	securityList.Append("🌍 Environment variable validation")

	return securitySection
}

func cliSection() (doyoucompute.Section, error) {
	cliSection := doyoucompute.NewSection("CLI Usage")

	cliSection.WriteIntro().
		Text("Create a CLI wrapper for your documents:")

	sample, err := os.ReadFile("./docs/pkg/documents/samples/app.go")
	if err != nil {
		return doyoucompute.Section{}, err
	}

	cliSection.WriteCodeBlock("go", []string{string(sample)}, doyoucompute.Static)

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

	return cliSection, nil
}

func basicUsageSection() (doyoucompute.Section, error) {
	basicUsageSection := doyoucompute.NewSection("Basic Usage")
	basicUsageSection.WriteIntro().
		Text("Create a simple document with executable commands:")

	sample, err := os.ReadFile("./docs/pkg/documents/samples/basics.go")
	if err != nil {
		return doyoucompute.Section{}, err
	}

	basicUsageSection.WriteCodeBlock("go", []string{string(sample)}, doyoucompute.Static)

	return basicUsageSection, nil
}

func quickstartSection() (doyoucompute.Section, error) {
	quickStartSection := doyoucompute.NewSection("Quick Start")

	installationSection := quickStartSection.CreateSection("Installation")
	installationSection.WriteCodeBlock("bash", []string{"go get github.com/MoonMoon1919/doyoucompute"}, doyoucompute.Static)

	basicUsage, err := basicUsageSection()
	if err != nil {
		return doyoucompute.Section{}, err
	}

	quickStartSection.AddSection(basicUsage)

	cliSection, err := cliSection()
	if err != nil {
		return doyoucompute.Section{}, err
	}

	quickStartSection.AddSection(cliSection)

	return quickStartSection, nil
}

func Readme() (doyoucompute.Document, error) {
	document, err := doyoucompute.NewDocument("DOYOUCOMPUTE")
	if err != nil {
		return doyoucompute.Document{}, err
	}

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
	quickStart, err := quickstartSection()
	if err != nil {
		return doyoucompute.Document{}, err
	}

	document.AddSection(quickStart)

	// Security
	securitySection := securitySection()

	config, err := configurationSecurity()
	if err != nil {
		return doyoucompute.Document{}, err
	}

	securitySection.AddSection(config)
	document.AddSection(securitySection)

	// Env vars
	envvars, err := environmentVariables()
	if err != nil {
		return doyoucompute.Document{}, err
	}
	document.AddSection(envvars)

	// Recs
	document.AddSection(recommendations())

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

	// Disclaimer
	disclaimerSection := document.CreateSection("Disclaimers")
	disclaimerSection.WriteIntro().
		Text("This work does not represent the interests or technologies of any employer, past or present.").
		Text("It is a personal project only.")

	return document, nil
}

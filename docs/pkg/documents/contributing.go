package documents

import (
	"github.com/MoonMoon1919/doyoucompute"
)

func gettingStarted() doyoucompute.Section {
	gettingStartedSection := doyoucompute.NewSection("Getting started")

	// Get familiar
	getFamiliar := gettingStartedSection.CreateSection("Get familiar with the project")
	getFamiliar.WriteParagraph().
		Text("Read the").
		Link("README", "README.md").
		Text("to understand the project's scope and purpose.")

	getFamiliar.WriteParagraph().
		Text("Look at the projects").
		Link("own documentation", "https://github.com/MoonMoon1919/doyoucompute/tree/main/docs").
		Text("to see real world usage.")

	// Chose a task
	choseATask := gettingStartedSection.CreateSection("Find a task")
	choseATask.WriteParagraph().
		Text("Browse the ").
		Link("issue tracker", "https://github.com/MoonMoon1919/doyoucompute/issues").
		Text(" to see what's being worked on and what needs attention.")

	choseATask.WriteParagraph().
		Text("Look for issues with these labels that are great for new contributors:")

	labelsList := choseATask.CreateList(doyoucompute.BULLET)
	labelsList.Append("\"good first issue\" - Small, well-defined tasks perfect for beginners")
	labelsList.Append("\"help wanted\" - Tasks where maintainers would appreciate assistance")
	labelsList.Append("\"documentation\" - Opportunities to improve docs and examples")
	labelsList.Append("\"enhancement\" - New features or improvements to existing functionality")

	choseATask.WriteParagraph().
		Text("Don't see anything that interests you? Feel free to open a new issue to:")

	suggestionsList := choseATask.CreateList(doyoucompute.BULLET)
	suggestionsList.Append("Suggest new features or improvements")
	suggestionsList.Append("Report documentation gaps or unclear examples")
	suggestionsList.Append("Propose improvements")
	suggestionsList.Append("Ask questions about implementation details")

	return gettingStartedSection
}

func codeContributions() doyoucompute.Section {
	codeContributions := doyoucompute.NewSection("Code contributions")

	// Guidelines
	setupSection := codeContributions.CreateSection("Setting Up Your Development Environment")
	setupSection.WriteParagraph().
		Text("First, fork the repository on GitHub at").
		Link("https://github.com/MoonMoon1919/doyoucompute", "https://github.com/MoonMoon1919/doyoucompute").
		Text(" by clicking the \"Fork\" button.")

	setupSection.WriteParagraph().
		Text("Then clone your forked repository to your local machine:")

	setupSection.WriteCodeBlock("bash", []string{"git clone <your_fork_url> doyoucompute"}, doyoucompute.Static)
	setupSection.WriteCodeBlock("bash", []string{"cd doyoucompute"}, doyoucompute.Static)

	setupSection.WriteParagraph().
		Text("Install dependencies and verify you can run the tests:")

	setupSection.WriteCodeBlock("bash", []string{"go mod tidy"}, doyoucompute.Static)
	setupSection.WriteCodeBlock("bash", []string{"go test ./..."}, doyoucompute.Static)

	// Development Workflow
	workflowSection := codeContributions.CreateSection("Development Workflow")
	workflowSection.WriteParagraph().
		Text("Create a new branch for your feature or bug fix:")

	workflowSection.WriteCodeBlock("bash", []string{"git checkout -b feature/my-awesome-feature"}, doyoucompute.Static)

	workflowSection.WriteParagraph().
		Text("Make your changes and add tests for new functionality. Run tests to ensure changes work as expected:")

	workflowSection.WriteCodeBlock("bash", []string{"go test ./..."}, doyoucompute.Static)

	workflowSection.WriteParagraph().
		Text("If you're adding new features, consider adding example usage in the examples directory.")

	// Submitting
	submissionSection := codeContributions.CreateSection("Submitting Your Changes")

	submissionSection.WriteParagraph().
		Text("Once you're satisfied with your changes, commit them with a descriptive message:")

	submissionSection.WriteCodeBlock("bash", []string{"git add ."}, doyoucompute.Static)
	submissionSection.WriteCodeBlock("bash", []string{"git commit -m \"Add feature: descriptive commit message\""}, doyoucompute.Static)

	submissionSection.WriteParagraph().
		Text("Push your changes to your forked repository:")

	submissionSection.WriteCodeBlock("bash", []string{"git push origin feature/my-awesome-feature"}, doyoucompute.Static)

	submissionSection.WriteParagraph().
		Text("Finally, create a pull request:")

	submissionSteps := submissionSection.CreateList(doyoucompute.BULLET)
	submissionSteps.Append("Go to the original repository on GitHub")
	submissionSteps.Append("Click \"Compare & pull request\"")
	submissionSteps.Append("Provide a clear description of your changes")
	submissionSteps.Append("Reference any relevant issues using #issue-number")
	submissionSteps.Append("Wait for review and address any feedback")

	return codeContributions
}

func reportingBugs() doyoucompute.Section {
	bugsSections := doyoucompute.NewSection("Reporting bugs")

	checkingSection := bugsSections.CreateSection("Checking for Existing Reports")
	checkingSection.WriteParagraph().
		Text("Before reporting a new bug, search the").
		Link("issue tracker", "https://github.com/MoonMoon1919/doyoucompute/issues").
		Text("to see if someone else has already reported the same issue.").
		Text("Check both open and closed issues - the bug might have been fixed in a recent version.")

	creatingSection := bugsSections.CreateSection("Reporting new bugs")

	creatingSection.WriteParagraph().
		Text("If you can't find an existing report, create a new issue and fill out the bug report form.")

	return bugsSections
}

func writingDocs() doyoucompute.Section {
	docsSection := doyoucompute.NewSection("Writing documentation")

	// Review existing documentation
	docsSection.WriteParagraph().
		Text("Read the").
		Link("README", "./README.md").
		Text("to understand the project's structure and how it's used.")

	// Identify areas for improvement
	docsSection.WriteParagraph().
		Text("Look for documentation that is unclear, incomplete, or outdated.")

	// Make the changes
	docsSection.WriteParagraph().
		Text("Update the appropriate file in the").
		Link("docs folder", "./docs").
		Text("since we're using doyoucompute to generate it's own documents")

	return docsSection
}

func Contributing() doyoucompute.Document {
	doc := doyoucompute.NewDocument("Contributing")

	doc.AddSection(gettingStarted())

	guidelinesSection := doc.CreateSection("Contribution guidelines")
	guidelinesSection.AddSection(codeContributions())
	guidelinesSection.AddSection(reportingBugs())
	guidelinesSection.AddSection(writingDocs())

	// License section
	licenseSection := doc.CreateSection("License")
	licenseSection.WriteParagraph().
		Text("By contributing, you agree that your contributions will be licensed under the project's").
		Link("MIT License.", "./LICENSE")

	return doc
}

package documents

import "github.com/MoonMoon1919/doyoucompute"

func BugReport() doyoucompute.Document {
	document := doyoucompute.NewDocument("Bug Report")

	document.AddFrontmatter(*doyoucompute.NewFrontmatter(map[string]interface{}{
		"name":      "Bug report",
		"about":     "Report a bug",
		"title":     "",
		"labels":    "",
		"assignees": "",
	}))

	expectedBehavior := document.CreateSection("Expected behavior")
	expectedBehavior.WriteParagraph().Text("What should happen?")

	actualBehavior := document.CreateSection("Actual behavior")
	actualBehavior.WriteParagraph().Text("What actually happens?")

	stepsToRepro := document.CreateSection("Steps to reproduce")
	reproList := stepsToRepro.CreateList(doyoucompute.NUMBERED)
	reproList.Append("") // Intentionally empty
	reproList.Append("")
	reproList.Append("")

	environmentDetails := document.CreateSection("Environment details")
	environmentDetails.WriteParagraph().Text("Tell us what go version, os, package version, etc.")

	codeSamples := document.CreateSection("Code samples")
	codeSamples.WriteParagraph().Text("Share a snippet of code that demonstrates the bug.")

	errorMessages := document.CreateSection("Error Messages")
	errorMessages.WriteParagraph().Text("Add any relevant error messages/logs here.")

	return document
}

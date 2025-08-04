package documents

import "github.com/MoonMoon1919/doyoucompute"

func PullRequest() doyoucompute.Document {
	document := doyoucompute.NewDocument("Pull request template")

	description := document.CreateSection("Description")
	description.WriteParagraph().
		Text("What is this change and why are you making it?")

	issue := document.CreateSection("Related issue")
	issue.WriteParagraph().
		Text("Link to the relevant issue here.")

	testing := document.CreateSection("How I tested")
	testing.WriteParagraph().
		Text("How did you test these changes?")

	return document
}

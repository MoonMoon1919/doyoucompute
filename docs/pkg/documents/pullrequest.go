package documents

import "github.com/MoonMoon1919/doyoucompute"

func PullRequest() doyoucompute.Document {
	document := doyoucompute.NewDocument("Pull request template")

	description := document.CreateSection("Description")
	description.WriteComment("What is this change and why are you making it?")

	issue := document.CreateSection("Related issue")
	issue.WriteComment("Link to the relevant issue here.")

	testing := document.CreateSection("How I tested")
	testing.WriteComment("How did you test these changes?")

	return document
}

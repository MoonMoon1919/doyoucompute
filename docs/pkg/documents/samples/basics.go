package samples

import "github.com/MoonMoon1919/doyoucompute"

func basics() {
	doc, err := doyoucompute.NewDocument("My Project")
	if err != nil {
		panic(err)
	}

	// Add an introduction
	doc.WriteIntro().
		Text("Welcome to my project! ").
		Text("Follow these steps to get started.")

	// Add a setup section with executable commands
	setup := doc.CreateSection("Setup")
	setup.WriteParagraph().
		Text("First, install dependencies:")

	setup.WriteCodeBlock("bash", []string{"npm install"}, doyoucompute.Exec)

	setup.WriteParagraph().
		Text("Then start the development server:")

	setup.WriteCodeBlock("bash", []string{"npm run dev"}, doyoucompute.Exec)
}

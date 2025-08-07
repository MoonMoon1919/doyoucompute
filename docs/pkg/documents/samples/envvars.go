package samples

import "github.com/MoonMoon1919/doyoucompute"

func envvars() {
	setup := doyoucompute.NewSection("Setup")

	setup.WriteExecutable(
		"bash",
		[]string{"curl", "-H", "Authorization: Bearer $API_KEY", "api.example.com"},
		[]string{"API_KEY"})
}

package version

import "fmt"

func PrintVersion(buildVersion string, buildDate string, buildCommit string) {
	fmt.Print("Build version: ")
	if buildVersion != "" {
		fmt.Printf("%s\n", buildVersion)
	} else {
		fmt.Println("N/A")
	}

	fmt.Print("Build date: ")
	if buildDate != "" {
		fmt.Printf("%s\n", buildDate)
	} else {
		fmt.Println("N/A")
	}

	fmt.Print("Build commit: ")
	if buildCommit != "" {
		fmt.Printf("%s\n", buildCommit)
	} else {
		fmt.Println("N/A")
	}
}

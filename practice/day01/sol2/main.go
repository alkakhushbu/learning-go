package main

import "fmt"

func main() {

	var (
		projectName      string  = "Learning go"
		codeLinesCount   uint8   = 100
		bugsCount        int     = 2
		projectCompleted bool    = false
		averageLines     float64 = 15.4
		teamLeadName     string  = "Alka"
		deadline         int     = 30
	)

	fmt.Printf("Project name: %s\n", projectName)
	fmt.Printf("Code lines written: %d\n", codeLinesCount)
	fmt.Printf("Bugs found: %d\n", bugsCount)
	fmt.Printf("Project completed? %t\n", projectCompleted)
	fmt.Printf("Average lines of code written per hour: %f\n", averageLines)
	fmt.Printf("Team lead name: %s\n", teamLeadName)
	fmt.Printf("Project deadline in days: %d\n", deadline)

	var uintOverflow uint8 = 255
	fmt.Printf("uint8 value before overflow: %d\n", uintOverflow)
	uintOverflow += 1
	fmt.Printf("uint8 value after overflow: %d\n", uintOverflow)
}

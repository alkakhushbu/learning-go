package main

import (
	"fmt"
)

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

	fmt.Printf("Project name: %s, type: %T\n", projectName, projectName)
	fmt.Printf("Code lines written: %d, type: %T\n", codeLinesCount, codeLinesCount)
	fmt.Printf("Bugs found: %d, type: %T\n", bugsCount, bugsCount)
	fmt.Printf("Project completed? %t, type: %T\n", projectCompleted, projectCompleted)
	fmt.Printf("Average lines of code written per hour: %f, type: %T\n", averageLines, averageLines)
	fmt.Printf("Team lead name: %s, type: %T\n", teamLeadName, teamLeadName)
	fmt.Printf("Project deadline in days: %d, type: %T\n", deadline, deadline)
}

// q4. Print default values and Type names of variables from question 2 using printf
// // Quick Tip, Use %v if not sure about what verb should be used,
// // but don't use it in this question :)
// // but generally using %v should be fine

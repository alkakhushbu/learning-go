package sum

import "testing"

// to create a test file write _test.go, to signal it is a test file not regular file

// go test ./... -v // run all the test for the project

// So function names must start with word Test to signal it is a test
// helper functions could be present in this file
//that would not be part of test so the would not start with the word Test

func TestSumInt(t *testing.T) {
	// Figure out two things
	// What are inputs (parameters)
	// What are outputs (return values)

	input := []int{1, 2, 3, 4, 5}
	want := 15

	got := SumInt(input)

	// Checking if the function's output matches the expected output.
	// If they are not equal, the test will fail and print the error message with the expected and got values.
	if got != want {
		// test would continue on if test case fail
		t.Errorf("sum of 1 to 5 should be %v; got %v", want, got)

		// Uncomment next line to stop the test if it fails at this point.
		//t.Fatalf("sum of 1 to 5 should be %v; got %v", want, got)
	}

	want = 0
	got = SumInt(nil)
	if got != want {
		t.Errorf("sum of nil should not be %v; got %v", want, got)
	}

}

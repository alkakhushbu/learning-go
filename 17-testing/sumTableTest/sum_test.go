package sum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// to create a test file write _test.go, to signal it is a test file not regular file

// So function names must start with word Test to signal it is a test
// helper functions could be present in this file
//that would not be part of test so the would not start with the word Test

func TestSumInt(t *testing.T) {

	// table test
	tt := [...]struct {
		// always give a name to the test
		name  string // columns
		input []int
		want  int
	}{
		// rows in the table, each struct object denotes a row
		{
			name:  "one to five numbers",
			input: []int{1, 2, 3, 4, 5},
			want:  15,
		},
		{
			name:  "nil slice",
			input: nil,
			want:  0,
		},
		{
			name:  "1 -1",
			input: []int{1, -1},
			want:  0,
		},
	}

	for _, tc := range tt {

		// t.Run creates a subtest with a name defined in the struct
		t.Run(tc.name, func(t *testing.T) {
			got := SumInt(tc.input)

			// require would fail the current test
			require.Equal(t, tc.want, got)
		})

	}
}

//*******************************************************************//
// Benchmark
//*******************************************************************//

// go test -run none -bench=BenchmarkAppendAMillionElems
// go test -run none -bench=BenchmarkAppendAMillionElems -benchmem
func BenchmarkSumInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumInt([]int{1, 2, 3, 4, 5})
	}
}

func BenchmarkAppendAMillionElems(b *testing.B) {
	var a []int = make([]int, 10000000)
	for i := 0; i < b.N; i++ {
		a = AppendAMillionElems(a)
	}
}

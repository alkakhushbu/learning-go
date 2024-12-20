package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTotalSalary(t *testing.T) {
	tt := [...]struct {
		name  string
		input User
		want  int
	}{
		{
			name:  "First test",
			input: User{Name: "Jon", Email: "jon@email.com", Pay: 100, Bonus: 150},
			want:  250,
		},
		{
			name:  "Second test",
			input: User{Name: "Jenna", Email: "jenna@email.com"},
			want:  0,
		},
		{
			name:  "First test",
			input: User{Name: "Jane", Email: "jane@email.com", Pay: 20, Bonus: 50},
			want:  70,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.TotalSalary()
			require.Equal(t, tc.want, got)
		})
	}

}

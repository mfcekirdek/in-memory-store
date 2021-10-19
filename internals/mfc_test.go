//+build unit

package internals

import "testing"

func TestA(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"A", "YES"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := A(); got != tt.want {
				t.Errorf("A() = %v, want %v", got, tt.want)
			}
		})
	}
}

package diff

import "testing"

func TestChanged(t *testing.T) {
	tests := []struct {
		name      string
		current   string
		requested *string
		want      *string
	}{
		{name: "not requested", current: "a", requested: nil, want: nil},
		{name: "same value", current: "a", requested: new("a"), want: nil},
		{name: "new value", current: "a", requested: new("b"), want: new("b")},
		{name: "cleared value", current: "a", requested: new(""), want: new("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Changed(tt.current, tt.requested)
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Errorf("Changed() = %v, want %v", got, tt.want)
			}
		})
	}
}

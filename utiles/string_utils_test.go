package utiles

import (
	"testing"
)

func TestStrListContain(t *testing.T) {

	tests := []struct {
		name      string
		list      []string
		term      string
		isContain bool
	}{
		{
			name:      "test string list contain",
			list:      []string{"item1", "item2"},
			term:      "item1",
			isContain: true,
		},
		{
			name:      "test string list not contain",
			list:      []string{"item1", "item2"},
			term:      "item3",
			isContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrListContain(tt.list, tt.term)
			if result != tt.isContain {
				t.Errorf("StrListContain(%s, %s) returnd %v, we want %v", tt.list, tt.term, result, tt.isContain)
			}
		})
	}
}

package pocketsmith

import "testing"

func Test_isNil(t *testing.T) {
	tests := map[string]struct {
		i    interface{}
		want bool
	}{
		"empty interface": {
			want: true,
		},
		"i has a string value": {
			i:    "hello world",
			want: false,
		},
		"i has a map value": {
			i: map[string]int{
				"hello": 1,
				"world": 2,
			},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isNil(tt.i)
			if got != tt.want {
				t.Errorf("unexpected value returned; got: %v, want: %v\n", got, tt.want)
				return
			}
		})
	}
}

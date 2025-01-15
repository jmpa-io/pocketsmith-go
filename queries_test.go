package pocketsmith

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_setupQueries(t *testing.T) {
	tests := map[string]struct {
		queries map[string]string
		want    url.Values
	}{
		"setup queries": {
			queries: map[string]string{
				"page_size": "10",
				"hello":     "world",
				"this is":   "a test",
			},
			want: url.Values{
				"page_size": []string{"10"},
				"hello":     []string{"world"},
				"this is":   []string{"a test"},
			},
		},
		"check defaults page size": {
			want: url.Values{
				"page_size": []string{"100"},
			},
		},
	}
	for name, tt := range tests {

		// run tests.
		t.Run(name, func(t *testing.T) {
			got := setupQueries(tt.queries)

			// is there a mismatch from what we're expecting vs what we've got?
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"setupQueries() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
			}
		})
	}
}

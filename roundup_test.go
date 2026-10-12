package pocketsmith

import "testing"

func TestIsRoundUp(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"Woolworths, including -1.00 Round-up", true},
		{"Grill'd, including $0.42 Round-up", true},
		{"Fiftyfive, including -0.71 Round-up", true},
		{"Woolworths", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := IsRoundUp(tc.in); got != tc.want {
			t.Errorf("IsRoundUp(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestExtractRoundUp(t *testing.T) {
	cases := []struct {
		in        string
		wantPayee string
		wantAmt   float64
		wantOK    bool
	}{
		{"Woolworths, including -1.00 Round-up", "Woolworths", 1.00, true},
		{"Grill'd, including $0.42 Round-up", "Grill'd", 0.42, true},
		{"Fiftyfive, including -0.71 Round-up", "Fiftyfive", 0.71, true},
		{"Woolworths", "", 0, false},
		{"", "", 0, false},
	}
	for _, tc := range cases {
		payee, amt, ok := ExtractRoundUp(tc.in)
		if ok != tc.wantOK {
			t.Errorf("ExtractRoundUp(%q) ok=%v, want %v", tc.in, ok, tc.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if payee != tc.wantPayee {
			t.Errorf("ExtractRoundUp(%q) payee=%q, want %q", tc.in, payee, tc.wantPayee)
		}
		if amt != tc.wantAmt {
			t.Errorf("ExtractRoundUp(%q) amt=%v, want %v", tc.in, amt, tc.wantAmt)
		}
	}
}

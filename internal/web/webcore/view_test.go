package webcore

import "testing"

func TestTicketID(t *testing.T) {
	cases := map[int]string{0: "TKT-1000", 1: "TKT-1001", 42: "TKT-1042"}
	for index, want := range cases {
		if got := TicketID(index); got != want {
			t.Errorf("TicketID(%d) = %q, want %q", index, got, want)
		}
	}
}

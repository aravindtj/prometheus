package dns

import (
	"net"
	"testing"
	"time"
)

func getIP(s string) string {
	a, e := net.LookupAddr(s)
	if e != nil {
		return ""
	}
	return a[0]
}

// flaky, need to setup local server and test from
// that.
func testClientAXFR(t *testing.T) {
	if testing.Short() {
		return
	}
	m := new(Msg)
	m.SetAxfr("miek.nl.")

	server := getIP("linode.atoom.net")

	tr := new(Transfer)

	if a, err := tr.In(m, net.JoinHostPort(server, "53")); err != nil {
		t.Log("failed to setup axfr: " + err.Error())
		t.Fatal()
	} else {
		for ex := range a {
			if ex.Error != nil {
				t.Logf("error %s\n", ex.Error.Error())
				t.Fail()
				break
			}
			for _, rr := range ex.RR {
				t.Logf("%s\n", rr.String())
			}
		}
	}
}

// fails.
func testClientAXFRMultipleEnvelopes(t *testing.T) {
	if testing.Short() {
		return
	}
	m := new(Msg)
	m.SetAxfr("nlnetlabs.nl.")

	server := getIP("open.nlnetlabs.nl.")

	tr := new(Transfer)
	if a, err := tr.In(m, net.JoinHostPort(server, "53")); err != nil {
		t.Log("Failed to setup axfr" + err.Error() + "for server: " + server)
		t.Fail()
		return
	} else {
		for ex := range a {
			if ex.Error != nil {
				t.Logf("Error %s\n", ex.Error.Error())
				t.Fail()
				break
			}
		}
	}
}

func testClientTsigAXFR(t *testing.T) {
	if testing.Short() {
		return
	}
	m := new(Msg)
	m.SetAxfr("example.nl.")
	m.SetTsig("axfr.", HmacMD5, 300, time.Now().Unix())

	tr := new(Transfer)
	tr.TsigSecret = map[string]string{"axfr.": "so6ZGir4GPAqINNh9U5c3A=="}

	if a, err := tr.In(m, "176.58.119.54:53"); err != nil {
		t.Log("failed to setup axfr: " + err.Error())
		t.Fatal()
	} else {
		for ex := range a {
			if ex.Error != nil {
				t.Logf("error %s\n", ex.Error.Error())
				t.Fail()
				break
			}
			for _, rr := range ex.RR {
				t.Logf("%s\n", rr.String())
			}
		}
	}
}

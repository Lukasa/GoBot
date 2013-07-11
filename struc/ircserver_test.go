package struc

import "testing"

func TestIRCServerFromHostnamePort(t *testing.T) {
	hostnameports := []string{
		"chat.freenode.net:6667",
		"127.0.0.1:8080",
		"127.0.0.1",
	}

	names := []string{
		"chat.freenode.net",
		"127.0.0.1",
		"127.0.0.1",
	}

	ports := []uint64{
		uint64(6667),
		uint64(8080),
		uint64(6667),
	}

	for i, teststr := range hostnameports {
		server, err := NewIRCServerFromHostnamePort(teststr)

		if err != nil {
			t.Error(err)
		} else if server == nil {
			t.Error("No server returned.")
		} else {
			if server.Name != names[i] {
				t.Errorf(
					"Invalid name: expected %v, got %v",
					names[i],
					server.Name)
			} else if server.Port != ports[i] {
				t.Errorf(
					"Invalid port: expected %v, got %v",
					ports[i],
					server.Port)
			}
		}
	}
}

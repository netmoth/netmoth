package strategies

import "testing"

func TestStrategies_MapContainsAll(t *testing.T) {
	m := Strategies()
	for _, key := range []string{"pcap", "afpacket", "pfring", "ebpf"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("missing strategy %s", key)
		}
	}
}


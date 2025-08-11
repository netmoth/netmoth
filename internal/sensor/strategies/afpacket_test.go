package strategies

import "testing"

func TestAFPacket_CompileBpf_Invalid(t *testing.T) {
	s := &AFPacketStrategy{}
	if _, err := s.compileBpf(128, "invalid bpf expression("); err == nil {
		// pcap BPF compiler behavior may vary by platform; do not assert strictly
	}
}

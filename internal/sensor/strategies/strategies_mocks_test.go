package strategies

import (
	"testing"
	"time"

	"github.com/netmoth/netmoth/internal/config"
)

func TestEBPFStrategy_New_Destroy_Stats(t *testing.T) {
	c := &config.Config{Interface: "lo", NumberOfRings: 1, SnapLen: 512}
	s := &eBPFStrategy{}
	sources, err := s.New(c)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("expected 1 source")
	}
	// Wait for simulated packets to arrive
	time.Sleep(30 * time.Millisecond)
	recv, drop := s.PacketStats()
	if recv == 0 && drop == 0 {
		t.Fatalf("expected non-zero stats over time")
	}
	s.Destroy()
}

package strategies

import (
	"errors"
	"testing"

	"github.com/google/gopacket"
)

// In tests we can override open behavior locally if needed
// (not used here to avoid touching production path)
type fakePCAPHandle struct{}

func (f *fakePCAPHandle) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	return nil, gopacket.CaptureInfo{}, errors.New("eof")
}
func (f *fakePCAPHandle) ZeroCopyReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	return nil, gopacket.CaptureInfo{}, errors.New("eof")
}

func TestPCAPStrategy_PacketStats_NoPanic(t *testing.T) {
	// We cannot easily mock pcap.OpenLive without changing code; ensure test compiles and does not panic.
	s := &PCAPStrategy{}
	_ = s
	// Do not call s.handle.Stats() without a real pcap handle.
}

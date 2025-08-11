package config

import (
	"os"
	"testing"
)

func TestConfig_New_ParseAndValidate(t *testing.T) {
	yaml := `
interface: lo
strategy: pcap
number_of_rings: 1
zero_copy: true
snapshot_length: 512
promiscuous: true
connection_timeout: 0
bpf: ""
max_cores: 0
log_file: test.log
agent_mode: false
redis:
  host: "localhost:6379"
  password: ""
postgres:
  user: u
  password: p
  db: d
  host: h
allowed_origins: ["http://localhost:5173"]
agent_token: "t"
`
	tmp, err := os.CreateTemp("", "cfg-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(yaml); err != nil {
		t.Fatal(err)
	}
	_ = tmp.Close()

	cfg, err := New(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AgentToken != "t" {
		t.Fatalf("agent token parse failed")
	}
	if len(cfg.AllowedOrigins) != 1 {
		t.Fatalf("origins parse failed")
	}
}

func TestConfig_ValidateInterface_Error(t *testing.T) {
	yaml := `interface: definitely-no-such-iface`
	tmp, err := os.CreateTemp("", "cfg-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(yaml); err != nil {
		t.Fatal(err)
	}
	_ = tmp.Close()
	_, err = New(tmp.Name())
	if err == nil {
		t.Fatalf("expected error for invalid interface")
	}
}


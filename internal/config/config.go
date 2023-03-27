package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/google/gopacket/pcap"
	"gopkg.in/yaml.v3"
)

// Config is ...
type Config struct {
	Interface   string
	Promiscuous bool
	MaxCores    int `yaml:"max_cores"`
	ConnTimeout int `yaml:"connection_timeout"`
	SnapLen     int `yaml:"snapshot_length"`
	Bpf         string
	LogFile     string `yaml:"log_file"`
	Postgres    Postgres
	Redis       Redis
}

// Postgres is ...
type Postgres struct {
	User            string
	Password        string
	DB              string
	Host            string
	MaxConn         int `yaml:"max_conn"`
	MaxIDLecConn    int `yaml:"max_idlec_conn"`
	MaxLifeTimeConn int `yaml:"max_lifetime_conn"`
}

// Redis is ...
type Redis struct {
	Password string
	Host     string
}

// New is ...
func New(cf string) (cfg *Config, err error) {
	cfg = new(Config)
	contents, err := os.ReadFile(cf)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(contents, &cfg); err != nil {
		return nil, err
	}

	if cfg.MaxCores != 0 && cfg.MaxCores < runtime.NumCPU() {
		runtime.GOMAXPROCS(cfg.MaxCores)
	} else if cfg.MaxCores != 0 {
		fmt.Printf("[!] Warning: max_cores argument is invalid. Using %d cores instead", runtime.NumCPU())
	}
	if cfg.SnapLen == 0 {
		cfg.SnapLen = 262144
	}
	if cfg.LogFile == "" {
		cfg.LogFile = "analyzer.log"
	}

	if err := validateInterface(cfg.Interface); err != nil {
		return nil, err
	}

	if err := validateSnapshotLength(cfg.SnapLen); err != nil {
		return nil, err
	}

	return cfg, err
}

func validateInterface(iface string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		if device.Name == iface {
			return nil
		}
	}
	return errors.New("specified network interface does not exist")
}

func validateSnapshotLength(snapLen int) error {
	if snapLen < 64 {
		return errors.New("minimum snapshot length is 64")
	}
	if snapLen > 4294967295 {
		return errors.New("snapshot length must be an unsigned 32-bit integer")
	}
	return nil
}

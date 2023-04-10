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
	Redis         Redis
	Interface     string
	Strategy      string
	Bpf           string
	LogFile       string `yaml:"log_file"`
	Postgres      Postgres
	NumberOfRings int  `yaml:"number_of_rings"`
	MaxCores      int  `yaml:"max_cores"`
	SnapLen       int  `yaml:"snapshot_length"`
	ConnTimeout   int  `yaml:"connection_timeout"`
	ZeroCopy      bool `yaml:"zero_copy"`
	Promiscuous   bool
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
	if cfg.LogFile == "" {
		cfg.LogFile = "analyzer.log"
	}

	if err := validateInterface(cfg.Interface); err != nil {
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

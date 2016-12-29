package main

import (
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"github.com/martini-contrib/throttle"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Config struct {
	Hostname string `json:"hostname"`
	LogLevel string `json:"log_level,omitempty"`
	BaseUrl  string `json:"base_url,omitempty"`
	MountUi  string `json:"mount_ui,omitempty"`

	Listen map[string]protoConfig `json:"listen,omitempty"`

	DataDir     string                      `json:"datadir"`
	UsersConfig string                      `json:"users"`
	SyncSources []dsapid.SyncSourceResource `json:"sync,omitempty"`

	Throttle struct {
		Api throttleConfig `json:"api,omitempty"`
	} `json:"throttle,omitempty"`
}

type protoConfig struct {
	ListenAddress string     `json:"address,omitempty"`
	Key           string     `json:"key,omitempty"`
	Cert          string     `json:"cert,omitempty"`
	Acme          acmeConfig `json:"acme,omitempty"`
}

type acmeConfig struct {
	Email    string   `json:"email,omitempty"`
	Domains  []string `json:"domains,omitempty"`
	CacheDir string   `json:"cache_dir,omitempty"`
}

type throttleConfig struct {
	Limit  uint64   `json:"limit,omitempty"`
	Within Duration `json:"within,omitempty"`
}

func (me throttleConfig) ToQuota() *throttle.Quota {
	if me.Limit == 0 {
		return &throttle.Quota{
			Limit:  10,
			Within: time.Second,
		}
	}

	return &throttle.Quota{
		Limit:  me.Limit,
		Within: time.Duration(me.Within),
	}
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var s string
	var within time.Duration
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	within, err = time.ParseDuration(s)
	if err == nil {
		*d = Duration(within)
	}
	return err
}

func DefaultConfig() Config {
	return Config{
		Hostname: "localhost",
		LogLevel: flagLogLevel,
		BaseUrl:  "http://localhost:8000/",
		Listen: map[string]protoConfig{
			"http": protoConfig{
				ListenAddress: "0.0.0.0:8000",
			},
		},
	}
}

func (me *Config) Load(filename string) (err error) {
	if _, err := os.Stat(filename); err == nil {
		if data, err := ioutil.ReadFile(filename); err == nil {
			err = json.Unmarshal(data, me)
		}
	}

	return err
}

func (me *Config) Save(filename string) (err error) {
	os.MkdirAll(path.Dir(filename), 0770)

	if data, err := json.MarshalIndent(me, "", "  "); err == nil {
		err = ioutil.WriteFile(filename, data, 0666)
	}

	return err
}

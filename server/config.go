package main

import (
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Hostname string `json:"hostname"`
	BaseUrl  string `json:"base_url,omitempty"`
	MountUi  string `json:"mount_ui,omitempty"`

	Listen map[string]protoConfig `json:"listen,omitempty"`

	DataDir     string                      `json:"datadir"`
	UsersConfig string                      `json:"users"`
	SyncSources []dsapid.SyncSourceResource `json:"sync,omitempty"`
}

type protoConfig struct {
	ListenAddress string `json:"address,omitempty"`
	UseSSL        bool   `json:"ssl,omitempty"`
	Key           string `json:"key,omitempty"`
	Cert          string `json:"cert,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Hostname: "localhost",
		BaseUrl:  "http://localhost:8000/",
		Listen: map[string]protoConfig{
			"http": protoConfig{
				ListenAddress: "0.0.0.0:8000",
				UseSSL:        false,
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

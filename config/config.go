package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dswg-book/nautilus/peer"
)

type Config struct {
	Peers []*peer.Peer `json:"peers"`
}

func NewConfig() *Config {
	return &Config{
		Peers: []*peer.Peer{},
	}
}

func (c *Config) AddPeer(peer *peer.Peer) {
	c.Peers = append(c.Peers, peer)
}

func (c *Config) RemovePeer(peer *peer.Peer) {
	for i, p := range c.Peers {
		if p.Host == peer.Host && p.Port == peer.Port {
			c.Peers = append(c.Peers[:i], c.Peers[i+1:]...)
		}
	}
}

func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".nautilus")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}
	configPath := filepath.Join(configDir, "config.json")
	config, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, config, 0644)
}

func (c *Config) Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".nautilus")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}
	configPath := filepath.Join(configDir, "config.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		os.WriteFile(configPath, []byte("{\"peers\":[]}"), 0644)
		config, err = os.ReadFile(configPath)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(config, c)
}

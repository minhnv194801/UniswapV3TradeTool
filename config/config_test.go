package config

import (
	"testing"
)

func TestNewWatcherConfig(t *testing.T) {
	watcher, err := NewWatcherConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(watcher.GetDexConf())
	t.Log(watcher.GetCoinsMap())
}

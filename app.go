package main

import (
	"context"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// dataPath returns the path to a json file next to the binary
func dataPath(name string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(execPath), name), nil
}

// loadJSON reads a json file next to the binary; returns "{}" if missing.
func loadJSON(name string) string {
	path, err := dataPath(name)
	if err != nil {
		return "{}"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// saveJSON atomically writes data to a json file via a tmp file + os.Rename.
func saveJSON(name, data string) bool {
	path, err := dataPath(name)
	if err != nil {
		return false
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(data), 0644); err != nil {
		return false
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return false
	}
	return true
}

// LoadResults reads results.json next to the binary. Returns "{}" if missing.
func (a *App) LoadResults() string { return loadJSON("results.json") }

// SaveResults atomically writes results.json.
func (a *App) SaveResults(data string) bool { return saveJSON("results.json", data) }

// LoadSettings reads settings.json (timezone, language). Returns "{}" if missing.
func (a *App) LoadSettings() string { return loadJSON("settings.json") }

// SaveSettings atomically writes settings.json.
func (a *App) SaveSettings(data string) bool { return saveJSON("settings.json", data) }

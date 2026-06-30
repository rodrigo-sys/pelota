package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func cacheDir() string {
	var base string
	switch runtime.GOOS {
	case "darwin":
		base = filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "windows":
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = os.Getenv("TEMP")
		}
	default:
		base = os.Getenv("XDG_CACHE_HOME")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), ".cache")
		}
	}
	return filepath.Join(base, "la18hd")
}

func eventsCachePath(source string) string {
	return filepath.Join(cacheDir(), "events-"+source+".json")
}

func cacheFresh(path string, ttl time.Duration) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return time.Since(fi.ModTime()) < ttl
}

func loadEvents(source string) ([]Event, error) {
	path := eventsCachePath(source)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func saveEvents(events []Event, source string) error {
	path := eventsCachePath(source)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), "events-*.json")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(events); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, path)
}



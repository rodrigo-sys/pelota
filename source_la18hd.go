package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type la18hdSource struct{}

func (s *la18hdSource) Name() string { return "la18hd" }

func (s *la18hdSource) FetchEvents() ([]Event, error) {
	resp, err := http.Get("https://la18hd.com/eventos/json/agenda123.json")
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return events, nil
}

func (s *la18hdSource) Resolve(url string) (string, error) {
	if !strings.Contains(url, "canales.php") {
		return url, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	var buf [8192]byte
	n, _ := resp.Body.Read(buf[:])
	m := playbackRE.FindSubmatch(buf[:n])
	if len(m) > 1 {
		return string(m[1]), nil
	}
	return "", nil
}

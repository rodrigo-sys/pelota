package main

import (
	"fmt"
	"strings"
)

type Event struct {
	Title    string `json:"title"`
	Time     string `json:"time"`
	Date     string `json:"date"`
	League   string `json:"league,omitempty"`
	Language string `json:"language"`
	Link     string `json:"link"`
	Channel  string `json:"channel,omitempty"`
}

type Source interface {
	Name() string
	FetchEvents() ([]Event, error)
	Resolve(url string) (string, error)
}

var sources = map[string]Source{
	"la18hd":      &la18hdSource{},
	"pirlotv":     &pirlotvSource{},
	"pelotalibre": &pelotalibreSource{},
}

func lookupSource(name string) (Source, error) {
	s, ok := sources[strings.ToLower(name)]
	if !ok {
		avail := make([]string, 0, len(sources))
		for k := range sources {
			avail = append(avail, k)
		}
		return nil, fmt.Errorf("unknown source %q (available: %s)", name, strings.Join(avail, ", "))
	}
	return s, nil
}

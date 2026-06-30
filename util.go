package main

import (
	"io"
	"regexp"
)

var playbackRE = regexp.MustCompile(`var playbackURL = "([^"]+)`)

func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

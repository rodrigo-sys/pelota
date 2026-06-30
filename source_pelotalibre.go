package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type pelotalibreSource struct{}

func (s *pelotalibreSource) Name() string { return "pelotalibre" }

var (
	liRE     = regexp.MustCompile(`(?s)<li class="([A-Z0-9]+)">\s*<a[^>]*>(.*?)<span class="t">([^<]+)</span></a>\s*<ul>(.*?)</ul>`)
	chLinkRE = regexp.MustCompile(`(?s)<a[^>]*href="[^"]*eventos\.html\?r=([^"&]+)[^"]*"[^>]*>(.*?)<span`)
	chanSpan = regexp.MustCompile(`\s*<span.*`)
)

func (s *pelotalibreSource) FetchEvents() ([]Event, error) {
	req, _ := http.NewRequest("GET", "https://librepelota.su/es/agenda/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	html := string(body)

	var events []Event
	leagueMatch := regexp.MustCompile(`^(.+?):\s*(.*)`)

	for _, m := range liRE.FindAllStringSubmatch(html, -1) {
		fullTitle := strings.TrimSpace(m[2])
		time_ := strings.TrimSpace(m[3])
		ulContent := m[4]

		league := ""
		title := fullTitle
		if lm := leagueMatch.FindStringSubmatch(fullTitle); lm != nil {
			league = strings.TrimSpace(lm[1])
			title = strings.TrimSpace(lm[2])
		}

		for _, ch := range chLinkRE.FindAllStringSubmatch(ulContent, -1) {
			b64 := ch[1]
			chanFull := strings.TrimSpace(ch[2])

			decoded, _ := decodeB64(b64)

			chan_ := chanSpan.ReplaceAllString(chanFull, "")
			chan_ = strings.TrimSpace(chan_)
			if i := strings.Index(chan_, " | "); i >= 0 {
				chan_ = chan_[:i]
			}

			events = append(events, Event{
				Title:    title,
				Time:     time_,
				Date:     "Hoy",
				League:   league,
				Language: "",
				Link:     decoded,
				Channel:  strings.ToLower(strings.ReplaceAll(chan_, " ", "-")),
			})
		}
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}
	return events, nil
}

func decodeB64(s string) (string, error) {
	// Add padding if needed
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *pelotalibreSource) Resolve(url string) (string, error) {
	if !strings.Contains(url, "latamvidzfy.org") && !strings.Contains(url, "vidzenvivo.cc") {
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

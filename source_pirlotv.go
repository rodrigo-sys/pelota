package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type pirlotvSource struct{}

func (s *pirlotvSource) Name() string { return "pirlotv" }

var (
	trRE    = regexp.MustCompile(`(?s)<tr[^>]*>(.*?)</tr>`)
	timeRE  = regexp.MustCompile(`class='t'[^>]*>([^<]+)`)
	tdRE    = regexp.MustCompile(`(?s)<td[^>]*>(.*?)</td>`)
	hrefRE  = regexp.MustCompile(`href=["']([^"']+)`)
	boldRE  = regexp.MustCompile(`<b>(.*?)</b>`)
	vivoRE  = regexp.MustCompile(`(?i)\s+en\s+Vivo\s*$`)
	phpRE   = regexp.MustCompile(`\.php$`)
)

func (s *pirlotvSource) FetchEvents() ([]Event, error) {
	req, _ := http.NewRequest("GET", "https://pirlotvplay.dev/", nil)
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
	for _, tr := range trRE.FindAllStringSubmatch(html, -1) {
		trContent := tr[1]
		timeM := timeRE.FindStringSubmatch(trContent)
		if timeM == nil {
			continue
		}
		time_ := strings.TrimSpace(timeM[1])

		tds := tdRE.FindAllStringSubmatch(trContent, -1)
		if len(tds) < 3 {
			continue
		}
		td := tds[2][1]

		league := ""
		if i := strings.Index(td, "<a "); i >= 0 {
			league = strings.TrimSpace(td[:i])
		}

		hrefM := hrefRE.FindStringSubmatch(td)
		if hrefM == nil {
			continue
		}
		link := hrefM[1]

		boldM := boldRE.FindStringSubmatch(td)
		if boldM == nil {
			continue
		}
		title := vivoRE.ReplaceAllString(strings.TrimSpace(boldM[1]), "")

		if !strings.HasPrefix(link, "http") {
			link = "https://pirlotvplay.dev/" + strings.TrimLeft(link, "/")
		}

		chan_ := phpRE.ReplaceAllString(link[strings.LastIndex(link, "/")+1:], "")

		events = append(events, Event{
			Title:    title,
			Time:     time_,
			Date:     "Hoy",
			League:   league,
			Language: "",
			Link:     link,
			Channel:  chan_,
		})
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("no events found (site format may have changed)")
	}
	return events, nil
}

func (s *pirlotvSource) Resolve(url string) (string, error) {
	return "", nil
}

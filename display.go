package main

import (
	"fmt"
	"strings"
)

var prettyChannels = map[string]string{
	"tycsports":      "TyC Sports",
	"tyc-sports":     "TyC Sports",
	"dsports":        "DSports",
	"dsportsar":      "DSports Argentina",
	"dsportsplus":    "DSports+",
	"winsports":      "Win Sports+",
	"win-sports+":    "Win Sports+",
	"disney1":        "Disney+",
	"vtvplus":        "VTV+",
	"telefe":         "Telefe",
	"azteca7":        "Azteca 7",
	"azteca-7":       "Azteca 7",
	"caracol":        "Caracol TV",
	"caracol-tv":     "Caracol TV",
	"caracoltv":      "Caracol TV",
	"foxsports1_usa": "Fox Sports 1",
	"fox-one":        "Fox Sports 1",
	"fox1":           "Fox Sports 1",
	"sportv":         "SPORTV",
	"eventos15":      "Evento 15",
	"universo-us":    "Universo US",
	"dazn-liga":      "DAZN Liga",
	"directv":        "DirecTV Sports",
	"directv-sports": "DirecTV Sports",
	"directv-plus":   "DirecTV+",
	"telemundo":      "Telemundo",
	"telemundousa":   "Telemundo US",
	"univision-deportes": "Univision Deportes",
	"tnt-sports-premium": "TNT Sports Premium",
	"liga1-max":      "Liga1 Max",
	"vix1":           "ViX",
	"vix3":           "ViX",
	"vix+":           "ViX",
	"unicanal":       "Unicanal",
	"das-erste":      "Das Erste",
	"trece":          "Trece",
	"tsn-1":          "TSN 1",
	"tsn1":           "TSN 1",
	"cazetv":         "CazeTV",
	"tvela1":         "TV ElA",
}

func prettyChannel(name string) string {
	if p, ok := prettyChannels[name]; ok {
		return p
	}
	return name
}

func channelFromLink(link string) string {
	if i := strings.Index(link, "stream="); i >= 0 {
		return link[i+7:]
	}
	if i := strings.Index(link, "ver/"); i >= 0 {
		s := link[i+4:]
		if j := strings.LastIndexByte(s, '.'); j >= 0 {
			s = s[:j]
		}
		return s
	}
	// https://domain.tld/CHANNEL.php
	if strings.HasSuffix(link, ".php") {
		if i := strings.LastIndexByte(link, '/'); i >= 0 {
			s := link[i+1:]
			s = s[:len(s)-4]
			return s
		}
	}
	return "Directo"
}

func displayLang(s string) string {
	if s == "" {
		return "\u2014"
	}
	return s
}

type streamLine struct {
	lang  string
	chan_ string
	link  string
}

func streamDisplayLines(events []Event, title string) []streamLine {
	var out []streamLine
	seen := map[string]bool{}
	for _, e := range events {
		if e.Title != title {
			continue
		}
		if seen[e.Link] {
			continue
		}
		seen[e.Link] = true
		out = append(out, streamLine{
			lang:  displayLang(e.Language),
			chan_: prettyChannel(channelFromLink(e.Link)),
			link:  e.Link,
		})
	}
	return out
}

func (s streamLine) Display() string {
	return fmt.Sprintf("%s│%s", s.lang, s.chan_)
}

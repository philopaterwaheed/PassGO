package ui

import (
	"log"
	"strings"
)

const (
	GitHubURL    = "https://github.com/philopaterwaheed"
	LinkedInURL  = "https://www.linkedin.com/in/philopater-waheed-561292227/"
	PortfolioURL = "https://philopaterwaheed.github.io/portfolio/"
	CVURL        = "https://docs.google.com/document/d/14806tbKGYbfnIlEs6Ovi67zxWmiPSZgB/export?format=pdf&name=Philopater_Waheed_CV.pdf"
	EmailURL     = "mailto:philopaterwaheed9@gmail.com"
)

func NormalizeURL(raw string) string {
	url := strings.TrimSpace(raw)
	if url == "" {
		return ""
	}

	lower := strings.ToLower(url)
	if strings.Contains(lower, "://") ||
		strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "tel:") {
		return url
	}

	return "https://" + url
}

func OpenURLLogged(raw string) {
	if err := OpenURL(raw); err != nil {
		log.Printf("Failed to open URL %q: %v", NormalizeURL(raw), err)
	}
}

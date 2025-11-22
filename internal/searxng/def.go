//go:generate msgp -tests=false

package searxng

import (
	"cmp"
	"encoding/json/v2"
	"slices"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Different date formats
	formats := []string{
		"2006-01-02T15:04:05",       // Without timezone (what SearXNG returns)
		"2006-01-02T15:04:05Z07:00", // RFC3339 with timezone
		"2006-01-02T15:04:05Z",      // RFC3339 with Z timezone
		time.RFC3339,                // Standard RFC3339
	}

	var parseErr error
	for _, format := range formats {
		if parsedTime, err := time.Parse(format, s); err == nil {
			t.Time = parsedTime
			return nil
		} else {
			parseErr = err
		}
	}

	return parseErr
}

// MarshalJSON implements custom JSON marshaling
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format(time.RFC3339))
}

// SearchResponse represents the complete SearXNG JSON API response
type SearchResponse struct {
	Query               string        `json:"query"`
	NumberOfResults     int           `json:"number_of_results"`
	Results             []Result      `json:"results"`
	Answers             []Answer      `json:"answers"`
	Corrections         []string      `json:"corrections"`
	Infoboxes           []Infobox     `json:"infoboxes"`
	Suggestions         []string      `json:"suggestions"`
	UnresponsiveEngines []EngineError `json:"unresponsive_engines"`
}

func OrderResults(r *[]Result) {
	for k, _ := range *r {
		switch (*r)[k].Priority {
		case "high":
			(*r)[k].Score *= 2
		case "low":
			(*r)[k].Score /= 2
		}
	}
	slices.SortFunc(*r, func(i, j Result) int {
		return cmp.Compare(j.Score, i.Score)
	})
}
func OrderForContext(r *[]Result) {
	for k, _ := range *r {
		if strings.Contains((*r)[k].ParsedURL[1], "youtube.com") || strings.Contains((*r)[k].ParsedURL[1], "vimeo.com") {
			(*r)[k].Score = 0
		}
	}
	OrderResults(r)
}

// Result represents a single search result item
type Result struct {
	// Base Result fields
	URL       string    `json:"url"`
	Engine    string    `json:"engine"`
	ParsedURL [6]string `json:"parsed_url,omitempty"`

	// MainResult fields
	Template      string   `json:"template"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	ImgSrc        string   `json:"img_src,omitempty"`
	IframeSrc     string   `json:"iframe_src,omitempty"`
	AudioSrc      string   `json:"audio_src,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	PublishedDate *Time    `json:"publishedDate,omitempty"`
	PubDate       string   `json:"pubdate,omitempty"`
	Length        string   `json:"length,omitempty"`
	Views         string   `json:"views,omitempty"`
	Author        string   `json:"author,omitempty"`
	Metadata      string   `json:"metadata,omitempty"`
	Priority      string   `json:"priority,omitempty"`
	Engines       []string `json:"engines,omitempty"`
	OpenGroup     bool     `json:"open_group,omitempty"`
	CloseGroup    bool     `json:"close_group,omitempty"`
	Positions     []int    `json:"positions,omitempty"`
	Score         float64  `json:"score,omitempty"`
	Category      string   `json:"category,omitempty"`
}

// Answer represents an answer result type
type Answer struct {
	Answer   string `json:"answer"`
	Template string `json:"template"`
	Engine   string `json:"engine"`
	URL      string `json:"url,omitempty"`
	Position int    `json:"position,omitempty"`
}

// Infobox represents an infobox result type
type Infobox struct {
	Infobox    string             `json:"infobox"`
	Template   string             `json:"template"`
	Engine     string             `json:"engine"`
	Content    string             `json:"content,omitempty"`
	ImgSrc     string             `json:"img_src,omitempty"`
	ID         string             `json:"id,omitempty"`
	URLs       []InfoboxURL       `json:"urls,omitempty"`
	Attributes []InfoboxAttribute `json:"attributes,omitempty"`
}

// InfoboxURL represents a URL within an infobox
type InfoboxURL struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// InfoboxAttribute represents an attribute within an infobox
type InfoboxAttribute struct {
	Label string            `json:"label"`
	Value string            `json:"value"`
	Image map[string]string `json:"image,omitempty"`
}

// EngineError represents an unresponsive engine error
type EngineError []string

/*type EngineError struct {
	Engine string `json:"engine"`
	Error  string `json:"error"`
	Type   string `json:"type,omitempty"`
}*/

// LegacyResult represents a legacy dictionary-based result for backward compatibility
type LegacyResult struct {
	// Base fields
	URL       string `json:"url"`
	Template  string `json:"template"`
	Engine    string `json:"engine"`
	ParsedURL string `json:"parsed_url,omitempty"`

	// MainResult fields
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	ImgSrc        string   `json:"img_src,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	Priority      string   `json:"priority,omitempty"`
	Engines       []string `json:"engines,omitempty"`
	Positions     []int    `json:"positions,omitempty"`
	Score         float64  `json:"score,omitempty"`
	Category      string   `json:"category,omitempty"`
	PublishedDate *Time    `json:"publishedDate,omitempty"`
	PubDate       string   `json:"pubdate,omitempty"`

	// Infobox fields
	URLs       []InfoboxURL       `json:"urls,omitempty"`
	Attributes []InfoboxAttribute `json:"attributes,omitempty"`

	// Legacy type indicators
	Answer          string `json:"answer,omitempty"`
	Suggestion      string `json:"suggestion,omitempty"`
	Correction      string `json:"correction,omitempty"`
	NumberOfResults int    `json:"number_of_results,omitempty"`
	EngineData      string `json:"engine_data,omitempty"`
}

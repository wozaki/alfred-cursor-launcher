package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

// Item represents an Alfred workflow item
type Item struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
	Icon     *Icon  `json:"icon,omitempty"`
	Valid    *bool  `json:"valid,omitempty"`
}

// Icon represents an Alfred workflow icon
type Icon struct {
	Path string `json:"path"`
}

// ScriptFilter represents the Alfred Script Filter JSON format
type ScriptFilter struct {
	Items []Item `json:"items"`
}

// NewScriptFilter creates a new ScriptFilter instance
func NewScriptFilter() *ScriptFilter {
	return &ScriptFilter{
		Items: []Item{},
	}
}

// AddItem adds an item to the script filter
func (sf *ScriptFilter) AddItem(item Item) {
	sf.Items = append(sf.Items, item)
}

// AddErrorItem adds an error item to the script filter
func (sf *ScriptFilter) AddErrorItem(message string) {
	valid := false
	sf.AddItem(Item{
		Title:    "Error: " + message,
		Subtitle: "Please check the logs for details",
		Valid:    &valid,
	})
}

// Print outputs the script filter as JSON to stdout, exits on error
func (sf *ScriptFilter) Print() {
	data, err := json.Marshal(sf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(data))
}

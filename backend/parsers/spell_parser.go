package parsers

import (
	"regexp"
	"strings"
)

// ParseSpell extracts spell name and description from HTML content
func ParseSpell(content string) (string, string) {
	// Parse Name
	// <title>Throw Corrosive Vial - Spells - Turtle WoW Database</title>
	var name string
	titleRegex := regexp.MustCompile(`<title>(.*?) - Spells`)
	if matches := titleRegex.FindStringSubmatch(content); len(matches) > 1 {
		name = matches[1]
	}

	// Parse Description
	var description string

	// Strategy 1: Look for the specific "Description" table cell
	// <th class="q1" style="border-top: none;">Description</th></tr><tr><td style="padding-top: 5px;">Throw a corrosive vial...</td></tr>
	descRegex := regexp.MustCompile(`<th>Description</th></tr><tr><td[^>]*>(.*?)</td></tr>`)

	// Strategy 2: Look for ANY text following "Use:" or "Equip:" in the main content area
	// This captures "Use: Throw a corrosive vial..." or "Equip: Poisons +5."
	// We want to capture until a significant HTML tag like <br>, </div>, or </td>
	// Removed </span> from stop list as descriptions often contain colored spans
	useRegex := regexp.MustCompile(`((Use|Equip): .*?)(<br>|</div>|</td>|<table|</tr>)`)

	if matches := descRegex.FindStringSubmatch(content); len(matches) > 1 {
		description = matches[1]
	} else if matches := useRegex.FindStringSubmatch(content); len(matches) > 1 {
		description = matches[1]
	} else {
		// Strategy 3: Look for <span class="q"> that is NOT the name or "Rank X"
		// Often custom spells have the description in a simple yellow span
		spanRegex := regexp.MustCompile(`<span class="q">([^<]+)</span>`)
		matches := spanRegex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			// Filter out common non-description strings
			// We moved the "contains name" check to allow "Poisons +5" for "Poisons" spell
			if len(m) > 1 && len(m[1]) > 5 && !strings.Contains(m[1], "Rank") {
				// Pick the longest one as a heuristic
				if len(m[1]) > len(description) {
					description = m[1]
				}
			}
		}
	}

	// Clean up HTML
	if description != "" {
		description = strings.ReplaceAll(description, "<br>", "\n")
		stripTags := regexp.MustCompile(`<[^>]*>`)
		description = stripTags.ReplaceAllString(description, "")
		description = strings.TrimSpace(description)
		description = cleanSpellDescription(description)
	}

	return name, description
}

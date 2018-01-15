package utils

import (
	"errors"
	"strings"
	"unicode"

	"github.com/dhickie/go-lgtv/control"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

// ErrNoMatchFound is returned if no match could be found for the target name or label
var ErrNoMatchFound = errors.New("No match was found for the target name or label")

// MatchInput will find the closest matched input based on the provided input name
func MatchInput(inputName string, inputs []control.Input) (control.Input, error) {
	// Remove all the whitespace from the input name, fuzzymatch doesn't like it
	stripped := stripWhitespace(inputName)

	// Build a list of input IDs and labels
	ids := make([]string, 0)
	labels := make([]string, 0)
	for _, v := range inputs {
		ids = append(ids, v.ID)
		labels = append(labels, v.Label)
	}

	// Rank all the matches to either the ID list of the label list
	idRanks := fuzzy.RankFind(stripped, ids)
	labelRanks := fuzzy.RankFind(stripped, labels)

	// Iterate through both and pick the one with the lowest distance
	var closestMatch = fuzzy.Rank{Distance: 99999}
	var matchFound = false
	var isLabel = false
	for _, v := range idRanks {
		matchFound = true
		if v.Distance < closestMatch.Distance {
			closestMatch = v
		}
	}
	for _, v := range labelRanks {
		matchFound = true
		if v.Distance < closestMatch.Distance {
			isLabel = true
			closestMatch = v
		}
	}

	// If no match was found at all, then return an error
	if !matchFound {
		return control.Input{}, ErrNoMatchFound
	}

	// Get the input that had this as its closest match
	for _, v := range inputs {
		if isLabel {
			if v.Label == closestMatch.Target {
				return v, nil
			}
		} else {
			if v.ID == closestMatch.Target {
				return v, nil
			}
		}
	}

	return control.Input{}, ErrNoMatchFound
}

func stripWhitespace(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)
}

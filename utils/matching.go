package utils

import (
	"errors"
	"strings"
	"unicode"

	"github.com/dhickie/hickhub/models"

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

	// If no match was found at all, then return an error
	if len(idRanks) == 0 && len(labelRanks) == 0 {
		return control.Input{}, ErrNoMatchFound
	}

	// Iterate through both and pick the one with the lowest distance
	var closestIDMatch = findLowestDistance(idRanks)
	var closestLabelMatch = findLowestDistance(labelRanks)

	isLabel := false
	closestMatch := closestIDMatch
	if closestLabelMatch.Distance < closestIDMatch.Distance {
		isLabel = true
		closestMatch = closestLabelMatch
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

// MatchChannel will find the closest matched channel based on the info in target
func MatchChannel(target models.SetChannelDetail, channels []control.Channel) (control.Channel, error) {
	// If the target has a channel number, then just use that
	if target.ChannelNumber > 0 {
		for _, v := range channels {
			if v.ChannelNumber == target.ChannelNumber {
				return v, nil
			}

		}
		return control.Channel{}, ErrNoMatchFound
	}

	// Strip whitespace out of the target channel name
	stripped := stripWhitespace(target.ChannelName)

	// Make a list of channel names
	names := make([]string, 0)
	for _, v := range channels {
		names = append(names, v.ChannelName)
	}

	// Rank the names against the target & find the closest match
	ranks := fuzzy.RankFind(stripped, names)
	closestMatch := findLowestDistance(ranks)

	// Go through the channels until we find this one
	for _, v := range channels {
		if v.ChannelName == closestMatch.Target {
			return v, nil
		}
	}

	return control.Channel{}, ErrNoMatchFound
}

// MatchApp will find the closest matched app based on the target app name
func MatchApp(appName string, apps []control.App) (control.App, error) {
	// Make a list of app names
	appNames := make([]string, 0)
	for _, v := range apps {
		appNames = append(appNames, v.Name)
	}

	// Strip whitespace out of the target app name
	stripped := stripWhitespace(appName)

	// Rank the app names & find the closest match
	ranks := fuzzy.RankFind(stripped, appNames)
	closestMatch := findLowestDistance(ranks)

	// Go through the apps and find the one with this name
	for _, v := range apps {
		if v.Name == closestMatch.Target {
			return v, nil
		}
	}

	return control.App{}, ErrNoMatchFound
}

func stripWhitespace(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)
}

func findLowestDistance(ranks []fuzzy.Rank) fuzzy.Rank {
	var closestMatch = fuzzy.Rank{Distance: 999999}
	for _, v := range ranks {
		if v.Distance < closestMatch.Distance {
			closestMatch = v
		}
	}

	return closestMatch
}

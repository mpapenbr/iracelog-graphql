package storage

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mpapenbr/iracelog-graphql/internal/events"
)

var ErrNoKeys = fmt.Errorf("no keys in searchArg")

func ExtractEventSearchKeys(arg string) (*events.EventSearchKeys, error) {
	ret := &events.EventSearchKeys{}
	foundKeys := false
	for _, v := range []string{"name", "car", "team", "driver", "track"} {
		regex := regexp.MustCompile(
			fmt.Sprintf("(?i)%s:\\s*(?P<arg>.+?)(?P<rest>(\\w+:)|$)", v))

		matches := regex.FindStringSubmatch(arg)
		if matches != nil {
			foundKeys = true
			m := strings.TrimSpace(matches[regex.SubexpIndex("arg")])
			switch v {
			case "name":
				ret.Name = m
			case "car":
				ret.Car = handleRegexSpecials(m)
			case "driver":
				ret.Driver = handleRegexSpecials(m)
			case "team":
				ret.Team = handleRegexSpecials(m)
			case "track":
				ret.Track = m
			}
		}
	}
	if foundKeys {
		return ret, nil
	}
	return nil, ErrNoKeys
}

// need to "extra" escape the \ since this value is used in a string in SQL
// see events.AdvancedEventSearch
func handleRegexSpecials(s string) string {
	return strings.ReplaceAll(regexp.QuoteMeta(s), "\\", "\\\\")
}

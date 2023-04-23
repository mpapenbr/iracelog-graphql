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
				ret.Car = m
			case "driver":
				ret.Driver = m
			case "team":
				ret.Team = m
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

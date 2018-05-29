package galaxy

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

// ParseFile returns the events of the given JSON file.
func ParseFile(jsonPath string) ([]Event, error) {
	buf, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var data [][]interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, errors.WithStack(err)
	}
	var events []Event
	for _, d := range data {
		event, err := parseEvent(d)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		events = append(events, event)
	}
	return events, nil
}

// parseEvent parses the given raw event from Python JSON format.
func parseEvent(data []interface{}) (Event, error) {
	coords := data[2].([]interface{})
	date, err := time.Parse("2006-01-02 15:04:05", data[1].(string))
	if err != nil {
		return Event{}, errors.WithStack(err)
	}
	return Event{
		T: date.UnixNano(),
		S: Point{X: coords[0].(float64), Y: coords[1].(float64)},
	}, nil
}

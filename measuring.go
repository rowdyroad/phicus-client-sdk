package phicus

import (
	"encoding/json"
	"errors"
	"time"
)

// Measuring struct
type Measuring struct {
	MeasuringID  string    `json:"-"`
	Key          string    `json:"key"`
	Value        string    `json:"value"`
	Lat          *float64  `json:"lat"`
	Lng          *float64  `json:"lng"`
	Display      *string   `json:"display"`
	FixationTime time.Time `json:"time"`
	Attachments  []string  `json:"attachments"`
}

// MarshalBinary marhaling Measuring struct
func (m Measuring) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary unmarhaling Measuring struct
func (m *Measuring) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}

// NewMeasuring creates new Measuring
func NewMeasuring(
	key string,
	value string,
	lat *float64,
	lng *float64,
	display *string,
	attachments []string) (*Measuring, error) {
	if key == "" {
		return nil, errors.New("Not valid measuring: key is undefined")
	}
	return &Measuring{
		Key:          key,
		Value:        value,
		Lat:          lat,
		Lng:          lng,
		Display:      display,
		FixationTime: time.Now(),
		Attachments:  attachments,
	}, nil
}

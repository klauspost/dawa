package dawa

import (
	"encoding/gob"
	"strings"
	"time"
)

// Since date/time is not a standard encoded field, we must create out own type.
type AwsTime time.Time

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		panic(err)
	}
	// Register it as Gob
	gob.Register(AwsTime{})
}

// ParseTime will return the time encoding for a single field
// It the input must be AWS formatted encoding
func ParseTime(s string) (*AwsTime, error) {
	result, err := time.ParseInLocation("2006-01-02T15:04:05.000", string(s), location)
	if err != nil {
		return nil, err
	}
	t := AwsTime(result)
	return &t, nil
}

// MustParseTime will return the time encoding for a single field
// It the input must be AWS formatted encoding
func MustParseTime(s string) AwsTime {
	result, err := time.ParseInLocation("2006-01-02T15:04:05.000", string(s), location)
	if err != nil {
		panic(err)
	}
	return AwsTime(result)
}

func (t AwsTime) MarshalText() (text []byte, err error) {
	return t.MarshalJSON()
}

func (t *AwsTime) UnmarshalText(text []byte) error {
	return t.UnmarshalJSON(text)
}

// UnmarshalJSON a single time field
// It will attempt AWS encoding, and if that fails standard UnmarshalJSON for time.Time
func (t *AwsTime) UnmarshalJSON(b []byte) error {
	unquoted := strings.Trim(string(b), "\" ")
	result, err := time.ParseInLocation("2006-01-02T15:04:05.000", unquoted, location)

	// Could not parse, attempt standard unmarshall
	if err != nil {
		var t2 time.Time
		err = t2.UnmarshalJSON([]byte(unquoted))
		if err != nil {
			return err
		}
		*t = AwsTime(t2)
		return nil
	}

	*t = AwsTime(result)
	return nil
}

// Time will return the underlying time.Time object
func (t AwsTime) Time() time.Time {
	return time.Time(t)
}

// MarshalJSON will send it as ordinary Javascipt date
func (t AwsTime) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// GobEncode (as time.Time)
func (t AwsTime) GobEncode() ([]byte, error) {
	return time.Time(t).GobEncode()
}

// GobDecode (as time.Time)
func (t *AwsTime) GobDecode(data []byte) error {
	return (*time.Time)(t).GobDecode(data)
}

/*
// GetBSON provides BSON encoding of the Kid
func (t AwsTime) GetBSON() (interface{}, error) {
	return time.Time(t), nil
}

// SetBSON provides BSON decoding
func (t *AwsTime) SetBSON(raw bson.Raw) error {
	var t2 time.Time
	err := raw.Unmarshal(&t2)
	*t = AwsTime(t2)
	return errgo.Mask(err)
}
*/

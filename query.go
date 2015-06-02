package dawa

import (
	"encoding/json"
	"fmt"
	"github.com/kpawlik/geojson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DefaultHost is the default host used for queries.
var DefaultHost = "http://dawa.aws.dk"

type parameter interface {
	Param() string
	Key() string
	IsMulti() bool
	AllValues() []string
	Merge(parameter) error
}

// A generic query structure
type query struct {
	host     string
	path     string
	params   map[string]parameter
	keys     []string // Keys in the order they were added
	warnings []error
}

type queryGeoJSON struct {
	query
}

// Add a key/value pair as additional parameter. It will be added as key=value on the URL.
// The values should not be delivered URL-encoded, that will be handled by the library.
func (q *query) Add(key, value string) {
	q.add(&textQuery{Name: key, Values: []string{value}, Multi: false, Null: true})
}

// This will  return any warnings that may have been generated while building the query.
func (q query) Warnings() []error {
	return q.warnings
}

// Returns true if any warnings have been generated.
func (q query) HasWarnings() bool {
	return len(q.warnings) > 0
}

func (q *query) add(p parameter) {
	if q.params == nil {
		q.params = make(map[string]parameter)
	}
	key := p.Key()
	pOld, ok := q.params[key]
	if ok {
		if !p.IsMulti() {
			q.warnings = append(q.warnings, fmt.Errorf("Ignoring second value of key %s", key))
			return
		}
		err := pOld.Merge(p)
		if err != nil {
			q.warnings = append(q.warnings, fmt.Errorf("Error while adding second value of key %s:%s", key, err.Error()))
		}
		q.params[key] = pOld
		return
	}
	q.keys = append(q.keys, key)
	q.params[key] = p
}

// WithHost allows overriding the host for this query.
//
// The default value is http://dawa.aws.dk
func (q *query) WithHost(s string) {
	q.host = s
}

// Replace the path of the query with something else.
func (q *query) OnPath(s string) {
	q.path = s
}

// Returns the URL for the generated query.
func (q query) URL() string {
	out := q.host + q.path
	if len(q.params) == 0 {
		return out
	}
	out += "?"

	for i, key := range q.keys {
		out += q.params[key].Param()
		if i < len(q.keys)-1 {
			out += "&"
		}
	}
	return out
}

type textQuery struct {
	Name   string
	Values []string
	Multi  bool
	Null   bool
}

// Merge other parameters into this one.
func (t *textQuery) Merge(other parameter) error {
	if t.Key() != other.Key() {
		return fmt.Errorf("merge: key value mismatch '%s' != '%s'", t.Key(), other.Key())
	}
	if !t.Multi {
		return fmt.Errorf("merge: cannot merge multiple values of key %s", t.Key())
	}
	t.Values = append(t.Values, other.AllValues()...)
	return nil
}

// Returns the entire parameter
func (t textQuery) Param() string {
	out := url.QueryEscape(t.Name) + "="
	if t.Null && len(t.Values) == 0 {
		return out
	}
	for i, val := range t.Values {
		out += url.QueryEscape(val)
		if !t.Multi {
			break
		}
		if i < len(t.Values)-1 {
			out += "|"
		}
	}
	return out
}

// Return values as string array, unencoded
func (t textQuery) AllValues() []string {
	return t.Values
}

func (t textQuery) Key() string {
	return t.Name
}

func (t textQuery) IsMulti() bool {
	return t.Multi
}

type RequestError struct {
	Type    string        `json:"type"`
	Title   string        `json:"title"`
	Details []interface{} `json:"details"`
	URL     string
}

func (r RequestError) Error() string {
	if r.Type == "" {
		return fmt.Sprintf("Error with request %s", r.URL)
	}
	return fmt.Sprintf("%s:%s. Details:%v. Request URL:%s", r.Type, r.Title, r.Details, r.URL)
}

// Perform the Request, and return the request result.
// If an error occurs during the request, or an error is reported
// this is returned.
// In some cases the error will be a RequestError type.
func (q query) Request() (io.ReadCloser, error) {
	url := q.URL()
	resp, err := http.Get(q.URL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 400 {
		return resp.Body, nil
	}
	u, e2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e2 != nil || len(u) == 0 {
		return nil, fmt.Errorf("Error with request %s", url)
	}
	var rerr RequestError
	_ = json.Unmarshal(u, &rerr)
	rerr.URL = url
	return nil, rerr
}

// Perform the Request, and return the request result as a geojson featurecollection.
// If an error occurs during the request, or an error is reported
// this is returned.
// In some cases the error will be a RequestError type.
func (q queryGeoJSON) GeoJSON() (*geojson.FeatureCollection, error) {
	q.Add("format", "geojson")
	url := q.URL()
	resp, err := http.Get(q.URL())
	if err != nil {
		return nil, err
	}

	u, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil || len(u) == 0 {
		return nil, fmt.Errorf("Error with request %s", url)
	}

	if resp.StatusCode >= 400 {
		rerr := RequestError{URL: url}
		_ = json.Unmarshal(u, &rerr)
		return nil, rerr
	}
	var fc geojson.FeatureCollection
	err = json.Unmarshal(u, &fc)
	if err != nil {
		return nil, err
	}
	return &fc, nil
}

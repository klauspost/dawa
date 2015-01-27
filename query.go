package dawa

import (
	"net/url"
)

// DefaultHost is the default host used for queries.
var DefaultHost = "http://dawa.aws.dk"

type parameter interface {
	Param() string
}

// A generic query structure
type query struct {
	host   string
	path   string
	params []parameter
}

func (q *query) Add(p parameter) {
	q.params = append(q.params, p)
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
	for i, value := range q.params {
		out += value.Param()
		if i != len(q.params)-1 {
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

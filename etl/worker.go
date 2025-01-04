package etl

import (
	"encoding/csv"
	"encoding/json"
)

const (
	InitSize  = iota
	ChunkSize // We can adjust the queue size accordingly, say 100 chunks
)

type subOrganizationOf struct {
	Name *string `json:"name"` // default nil reference
}

type publisher struct {
	Name              *string           `json:"name"`              // default nil reference
	SubOrganizationOf subOrganizationOf `json:"subOrganizationOf"` // default nil reference
}

type contactPoint struct {
	Fn *string `json:"fn"` // default nil instead of undefined
}

type document struct {
	Modified *string `json:"modified"` // default nil reference

	Publisher publisher    `json:"publisher"`    // default nil reference
	Contact   contactPoint `json:"contactPoint"` // default nil reference
	Keyword   []string     `json:"keyword"`      // array is a reference type
}

type csvData struct {
	Modified       string
	PubName        string
	PubSubOrgName  string
	ContactPointFn string
	Keyword        string
}

type accumulator []csvData

// for nil references replace with empty value
func validate(field *string) string {
	if field == nil {
		return ""
	}
	return *field
}

// method unmarshalling document data, single ownership of pipeline(channel)
// 'job' is ReadOnly Pipeline ensuring immutability
var extract = func(decoder *json.Decoder) <-chan document {
	jobs := make(chan document, ChunkSize) // buffered channel with initial cap of 1
	go func() {
		defer close(jobs)
		var doc document
		if _, err := decoder.Token(); err != nil {
			panic(err)
		}
		for decoder.More() {
			if err := decoder.Decode(&doc); err == nil {
				jobs <- doc
			}
		}
	}()
	return jobs
}

// single ownership of pipeline(channel)
// receive from 'job' and transform to csv data
var transform = func(results <-chan document) <-chan accumulator {
	jobs := make(chan accumulator, ChunkSize)
	go func() {
		defer close(jobs)
		for result := range results {
			// extract valid fields from document for each result
			// marshal into csvData
			// send csvData to a channel
			jobs <- mapper(result)
		}
	}()
	// return the channel as read only
	return jobs
}

/*
load each slice predicated by the keyword array
i.e. keywords = ["key1","key2","key3"] slice is:

	[

		{"modified": "meh",
		"publisher.name": "foo",
	    	"publisher.sub.name": "foobar",
	    	"contactPoint.fn": "bar",
	    	"keyword": "key1"},

		{"modified": "meh",
		"publisher.name": "foo",
	    	"publisher.sub.name": "foobar",
	    	"contactPoint.fn": "bar",
	    	"keyword": "key2"},
		...
	]
*/
var load = func(results <-chan accumulator, writer *csv.Writer) {
	for result := range results {
		if err := writeCsv(writer, result); err != nil {
			panic(err)
		}
	}
}

var mapper = func(docData document) []csvData {
	_accumulator := make(accumulator, 0)
	for _, keyword := range docData.Keyword {

		_accumulator = append(_accumulator, csvData{
			Modified:       validate(docData.Modified),
			PubName:        validate(docData.Publisher.Name),
			PubSubOrgName:  validate(docData.Publisher.SubOrganizationOf.Name),
			ContactPointFn: validate(docData.Contact.Fn),
			Keyword:        keyword,
		})
	}
	return _accumulator
}

package etl

import (
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
	Name              *string            `json:"name"`              // default nil reference
	SubOrganizationOf *subOrganizationOf `json:"subOrganizationOf"` // default nil reference
}

type contactPoint struct {
	Fn *string `json:"fn"` // default nil instead of undefined
}

type document struct {
	Modified *string `json:"modified"` // default nil reference

	Publisher *publisher    `json:"publisher"`    // default nil reference
	Contact   *contactPoint `json:"contactPoint"` // default nil reference
	Keyword   []string      `json:"keyword"`
}

type csvData struct {
	Modified       string `csv:"modified"`
	PubName        string `csv:"publisher.name"`
	PubSubOrgName  string `csv:"publisher.subOrganizationOf.name"`
	ContactPointFn string `csv:"contactPoint.fn"`
	Keyword        string `csv:"keyword"`
}

// method unmarshalling document data, single ownership of pipeline(channel)
// 'job' is ReadOnly Pipeline ensuring immutability
var extract = func(decoder *json.Decoder) <-chan document {
	jobs := make(chan document, ChunkSize) // buffered channel with initial cap of 1
	go func() {
		defer close(jobs)
		var doc document
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
var transform = func(results <-chan document) <-chan []csvData {
	jobs := make(chan []csvData, ChunkSize)
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
var load = func(results <-chan []csvData, et *ETL) {
	for result := range results {
		if err := writeCsv(et.out, result); err != nil {
			panic(err)
		}
	}
}

var mapper = func(docData document) []csvData {
	bucketSize := len(docData.Keyword)
	accumulator := make([]csvData, bucketSize)
	for _, keyword := range docData.Keyword {
		accumulator = append(accumulator, csvData{
			Modified:       *docData.Modified,
			PubName:        *docData.Publisher.Name,
			PubSubOrgName:  *docData.Publisher.Name,
			ContactPointFn: *docData.Contact.Fn,
			Keyword:        keyword,
		})
	}
	return accumulator
}

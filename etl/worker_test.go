package etl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	emptySubOrgName = ""
	modified        = "2017-05-15"
	publisherName   = "pub-name"
	subOrgName      = "sub-org-name"
	contactName     = "contact-name"

	validSingleJson = `[
{
    "modified": "2017-05-15",
	"publisher": {
	    "name": "pub-name-1",
	    "subOrganizationOf": {
		"name": "sub-org-name-1"
	    }
	},
	"contactPoint": {
	    "fn": "contact-name-1"
	},
	"keyword": [
	    "keyword1",
	"keyword2",
	"keyword3"
	]
}
]`
	validNestedJson = `[
{
    "modified": "2017-05-15",
	"publisher": {
	    "name": "pub-name-2",
	    "subOrganizationOf": {
		"name": "sub-org-name-2"
	    }
	},
	"contactPoint": {
	    "fn": "contact-name-2"
	},
	"keyword": [
	    "keyword1",
	"keyword2",
	"keyword3"
	]
},
{
    "modified": "2017-05-15",
    "publisher": {
	"name": "pub-name-3",
	"subOrganizationOf": {
	    "name": "sub-org-name-3"
	}
    },
    "contactPoint": {
	"fn": "contact-name-3"
    },
    "keyword": [
	"keyword4",
    "keyword5",
    "keyword6"
    ]
}
]`
	validSingleJsonWithoutSubOrganzationOf = `[{"modified": "2015-05-15", "contactPoint": { "fn": "contact-name-6"}, "keyword": ["keyword7", "keyword8", "keyword9"]}]`
	invalidJson                            = `{
    "modified": "2017-05-15",
	"publisher": {
	    "name": "pub-name-4",
	    "subOrganizationOf": {
		"name": "sub-org-name-4"
	    }
	},
	"contactPoint": {
	    "fn": "contact-name-4"
	},
	"keyword": [
	    "keyword10",
	"keyword11",
	"keyword12"
	]
},
{
    "modified": "2017-05-20",
    "publisher": {
	"name": "pub-name-5",
	"subOrganizationOf": {
	    "name": "sub-org-name-5"
	}
    },
    "contactPoint": {
	"fn": "contact-name-5"
    },
    "keyword": [
	"keyword13",
    "keyword14",
    "keyword15"
    ]
}`
)

const (
	Zero = iota
	One
	Two
	Three
	Four
	Five
	Six
)

func TestExtract_ValidJson_ReturnDoc(t *testing.T) {

	var tests = []struct {
		json      string
		output    document
		outputErr bool
	}{
		{validSingleJson, document{
			Modified: &modified,
			Publisher: publisher{
				Name: formatString(publisherName, One),
				SubOrganizationOf: subOrganizationOf{
					Name: formatString(subOrgName, One),
				},
			},
			Contact: contactPoint{Fn: formatString(contactName, One)},
			Keyword: []string{"keyword1", "keyword2", "keyword3"},
		}, false},
		{validSingleJsonWithoutSubOrganzationOf, document{
			Modified:  &modified,
			Publisher: publisher{},
			Contact:   contactPoint{Fn: formatString(contactName, Two)},
			Keyword:   []string{"keyword7", "keyword8", "keyword9"},
		}, false},
	}

	for _, test := range tests {
		et := newMockETL(test.json)
		doc := extract(json.NewDecoder(et.in))
		got := <-doc
		if validate(got.Publisher.Name) != validate(test.output.Publisher.Name) {
			t.Errorf("[TestExtract_ValidJson_ReturnDoc]: %v want %v", validate(got.Publisher.Name), validate(test.output.Publisher.Name))
		}
	}
}

func TestExtract_InvalidJson_ExceptionThrown(t *testing.T) {
	var tests = []struct {
		json      string
		output    document
		outputErr bool
	}{
		{invalidJson, document{}, true},
	}

	for _, test := range tests {
		et := newMockETL(test.json)
		doc := extract(json.NewDecoder(et.in))
		// illustrate that error closes channel without any data
		for {
			select {
			case <-doc:
				return
			default:
				t.Logf("[TestExtract_InvalidJson_ExceptionThrown]: Found Error: %t", test.outputErr)
			}
			break
		}
	}
}

func TestExtract_NestedValidJson_ReturnDocList(t *testing.T) {
	var tests = []struct {
		json      string
		output    []document
		outputErr bool
	}{
		{validNestedJson, []document{
			{
				Modified: &modified,
				Publisher: publisher{
					Name: formatString(publisherName, Two),
					SubOrganizationOf: subOrganizationOf{
						Name: formatString(subOrgName, Two),
					},
				},
				Contact: contactPoint{Fn: formatString(contactName, Two)},
				Keyword: []string{"keyword1", "keyword2", "keyword3"},
			},
			{
				Modified: &modified,
				Publisher: publisher{
					Name: formatString(publisherName, Three),
					SubOrganizationOf: subOrganizationOf{
						Name: formatString(subOrgName, Three),
					},
				},
				Contact: contactPoint{Fn: formatString(contactName, Three)},
				Keyword: []string{"keyword4", "keyword5", "keyword6"},
			},
		}, false},
	}

	for _, test := range tests {
		et := newMockETL(test.json)
		docs := extract(json.NewDecoder(et.in))
		for got := range docs {
			out := test.output[1]
			if *got.Publisher.Name != *out.Publisher.Name {
				t.Errorf("[TestExtract_NestedValidJson_ReturnDocList]: want\n %v found\n %v", *out.Publisher.Name, *got.Publisher.Name)
			}
		}
	}
}

func TestTransform_DocWithoutSomeFields_ReturnCsvWithEmptyValue(t *testing.T) {
	var tests = []struct {
		doc    document
		output []csvData
	}{
		{
			document{
				Modified: &modified,
				Publisher: publisher{
					Name: &publisherName,
					SubOrganizationOf: subOrganizationOf{
						Name: &subOrgName,
					},
				},
				Contact: contactPoint{Fn: &contactName},
				Keyword: []string{"random-keyword1", "random-keyword2", "random-keyword3"},
			}, []csvData{
				{
					Modified:       modified,
					PubName:        publisherName,
					PubSubOrgName:  publisherName,
					ContactPointFn: contactName,
					Keyword:        "random-keyword1",
				},
				{
					Modified:       modified,
					PubName:        publisherName,
					PubSubOrgName:  publisherName,
					ContactPointFn: contactName,
					Keyword:        "random-keyword2",
				},
				{
					Modified:       modified,
					PubName:        publisherName,
					PubSubOrgName:  publisherName,
					ContactPointFn: contactName,
					Keyword:        "random-keyword3",
				},
			},
		},
	}

	job := make(chan document)
	for _, test := range tests {
		go func() {
			defer close(job)
			job <- test.doc
		}()
		result := transform(job)
		for got := range result {
			if len(got) != len(test.output) {
				t.Errorf("[TestTransform_DocWithoutSomeFields_ReturnCsvWithEmptyValue]: got %d want %d", len(got), len(test.output))
			}
		}
	}
}

func newMockETL(jsonData string) *ETL {
	reader := strings.NewReader(jsonData)
	writer := bufio.NewWriter(os.Stdout)
	return NewEtl(reader, writer)
}

func formatString(field string, idx int) *string {
	str := fmt.Sprintf("%s-%d", field, idx)
	return &str
}

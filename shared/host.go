package shared

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/pkg/errors"
)

// SchedularType is an enum for the type of schedulars used to configure the
// host schedular
type SchedularType string

// Enums of SchedularType
const (
	RSS     SchedularType = "RSS"
	Sitemap SchedularType = "Sitemap"
)

// Host is the root level configuration object for the crawl3 library
type Host struct {
	Host      string              `json:"host"`
	Alias     []string            `json:"alias"`
	Schedular []HostSchedularOpts `json:"schedular"`
	Extractor []ExtractorOpts     `json:"extraction"`
}

// HostSchedularOpts provides the configuration of a schedular
type HostSchedularOpts struct {
	Type      SchedularType          `json:"type"`
	Frequency string                 `json:"frequency"`
	Data      map[string]interface{} `json:"data"`
}

// ExtractorOpts provide the extraction options to pull data from a website
type ExtractorOpts struct {
	Type      string               `json:"@type"`
	PageMatch []string             `json:"@pageMatcher"`
	Fields    map[string]FieldRule `json:"fields"`
}

// FieldRule defines the data required to pull a single piece of a data from a website
type FieldRule struct {
	Kind           string   `json:"type"`
	Matcher        string   `json:"matcher"`
	Content        string   `json:"content"`
	ContentMatcher []string `json:"content_match"`
	ExcludeMatch   []string `json:"excludeMatch"`
}

// LoadHostsFromDir will look in a directory for a list of JSON files and if possible
// load them into a slice of Host objects
func LoadHostsFromDir(dirname string) ([]Host, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to read hosts from dir {%s}", dirname)
	}

	var hosts = []Host{}
	for _, file := range files {
		filename := path.Join(dirname, file.Name())
		host, err := parseHost(filename)
		if err != nil {
			log.Println(err)
			continue
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}

func parseHost(path string) (Host, error) {
	hostJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return Host{}, errors.Wrap(err, "unable to read data from file")
	}

	return parseJSONHost(hostJSON)
}

func parseJSONHost(data []byte) (Host, error) {
	var host Host
	err := json.Unmarshal(data, &host)
	if err != nil {
		return Host{}, errors.Wrap(err, "unable to parse json data")
	}

	return host, nil
}

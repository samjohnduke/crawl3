package crawler

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/samjohnduke/crawl3/shared"
)

// JSONExtractor takes an array
type JSONExtractor struct {
	URL   string
	Rules []shared.ExtractorOpts
}

// Register the extractor in the list of extractors
func (je *JSONExtractor) Register(ex Extractors) {
	ex.Add(je)
}

// Match the extractor to a url host
func (je *JSONExtractor) Match() (url string) {
	return je.URL
}

// Extract a data object from the goquery.Document
func (je *JSONExtractor) Extract(doc *goquery.Document) (interface{}, error) {
	var harvestedData = []interface{}{}

	for _, rule := range je.Rules {
		var sel *goquery.Selection
		for _, pm := range rule.PageMatch {
			sel = doc.Find(pm)
		}

		if sel.Length() <= 0 {
			continue
		}

		data := map[string]interface{}{
			"type": rule.Type,
		}

		for field, fieldRule := range rule.Fields {

			var fieldSel *goquery.Selection
			fieldSel = sel.Find(fieldRule.Matcher)

			if fieldSel.Length() <= 0 {
				log.Printf("skipping { %s } as field not found in document", field)
				continue
			}

			if len(fieldRule.ExcludeMatch) > 0 {
				fieldSel = fieldSel.Children()

				for _, exclude := range fieldRule.ExcludeMatch {
					fieldSel = fieldSel.Not(exclude)
				}
			}

			var result interface{}

			switch fieldRule.Kind {
			case "String":
				if fieldRule.Content == "innerHTML" {
					result = fieldSel.Text()
				} else {
					rows := fieldSel.Map(func(_ int, sel *goquery.Selection) string {
						val, _ := sel.Attr(fieldRule.Content)
						return strings.TrimSpace(val)
					})
					result = strings.Join(rows, "\n\n")
				}
				break

			case "[]String":
				if fieldRule.Content == "innerHTML" {
					result = fieldSel.Map(func(_ int, sel *goquery.Selection) string {
						return strings.TrimSpace(sel.Text())
					})
				} else {
					result = fieldSel.Map(func(_ int, sel *goquery.Selection) string {
						val, _ := sel.Attr(fieldRule.Content)
						return strings.TrimSpace(val)
					})
				}
				break

			case "Time":
				if fieldRule.Content == "innerHTML" {
					var err error
					tString := strings.TrimSpace(fieldSel.First().Text())
					result, err = dateparse.ParseAny(tString)
					if err != nil {
						log.Println(err)
					}
				} else {
					t := fieldSel.First()
					tString, _ := t.Attr(fieldRule.Content)
					var err error
					result, err = dateparse.ParseAny(tString)
					if err != nil {
						log.Println(err)
					}
				}
				break
			default:
				log.Println("error, incorrect field rule type value, {" + fieldRule.Kind + "}")
			}

			data[field] = result
		}

		harvestedData = append(harvestedData, data)
	}
	return harvestedData, nil
}

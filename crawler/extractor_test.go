package crawler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDomainExtractor(t *testing.T) {
	resp, err := http.Get("https://www.domain.com.au/bourke-street-melbourne-vic-3000-10828929")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}

	extractor := funcExs[0]
	_, err = extractor.Extract(doc)
	if err != nil {
		t.Error(err)
	}
}

func TestABCExtractor(t *testing.T) {
	resp, err := http.Get("http://www.abc.net.au/news/2018-02-23/barnaby-joyce-resigns/9477942")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}

	extractor := funcExs[1]
	_, err = extractor.Extract(doc)
	if err != nil {
		t.Error(err)
	}
}

var funcExs = []*funcExtractor{
	&funcExtractor{"www.domain.com.au", func(doc *goquery.Document) (interface{}, error) {

		result := doc.Find(".listing-details__root").Length()
		if result == 0 {
			fmt.Println("no data found")
			return nil, nil
		}

		data := map[string]interface{}{}
		data["title"] = doc.Find("h1").Text()
		data["summary-title"] = doc.Find(".listing-details__summary-title").Text()
		data["features"] = doc.Find(".listing-details__summary-right-column .property-feature__feature-text-container").Map(func(_ int, sel *goquery.Selection) string { return sel.Text() })
		data["callout-details"] = doc.Find(".listing-details__summary-strip-items").Map(func(_ int, sel *goquery.Selection) string { return sel.Text() })
		data["description"] = doc.Find(".listing-details__description").Map(func(_ int, sel *goquery.Selection) string { return sel.Text() })

		return data, nil
	}},
	&funcExtractor{"www.abc.net.au", func(doc *goquery.Document) (interface{}, error) {
		result := doc.Find(".news.story_page").Length()
		if result == 0 {
			return nil, nil
		}

		data := map[string]interface{}{}
		data["title"] = doc.Find("h1").Text()
		data["byline"] = doc.Find(".byline a").Text()
		a := doc.Find(".article").
			Children().
			Not(".tools").
			Not(".attached-content").
			Not(".byline").
			Not(".btn-group").
			Not("h1").
			Not(".published").
			Not(".inline-content").
			Not(".topics").
			Not(".authorpromo").
			Not(".published").
			Map(func(_ int, sel *goquery.Selection) string {
				cont := strings.TrimSpace(sel.Text())
				return cont
			})
		b := make([]string, 0, len(a))
		for _, e := range a {
			if e != "" {
				b = append(b, e)
			}
		}

		data["article"] = strings.Join(b, "\n")

		return data, nil
	}},
}

package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/samjohnduke/crawl3/crawler"
)

var funcExs = []crawler.FuncExtractor{
	abcExtractor,
	sbsExtractor,
	ageExtractor,
	brisbaneTimesExtractor,
	hearldsunExtractor,
	alternetExtractor,
	canberraTimesExtractor,
}

type entity struct {
	Name string
	Link string
	Type string
}

var fairfaxEx = func(doc *goquery.Document) (interface{}, error) {
	article := map[string]interface{}{}

	article["title"] = doc.Find("h1").Text()

	article["images"] = doc.Find("article section div img").Map(func(_ int, sel *goquery.Selection) string {
		link, _ := sel.Attr("src")
		return link
	})

	article["content"] = strings.Join(doc.Find("article section div").Map(func(_ int, sel *goquery.Selection) string {
		return sel.Text() + "\n"
	}), " ")

	articles := []map[string]interface{}{article}
	data := map[string]interface{}{
		"articles": articles,
	}

	return data, nil
}

var hearldsunExtractor = crawler.FuncExtractor{
	URL: "www.heraldsun.com.au",
	Method: func(doc *goquery.Document) (interface{}, error) {
		article := map[string]interface{}{}

		article["title"] = doc.Find("h1").Text()

		article["images"] = doc.Find(".story .tg-tlc-storybody img").Map(func(_ int, sel *goquery.Selection) string {
			link, _ := sel.Attr("src")
			return link
		})

		article["content"] = strings.Join(doc.Find(".story .tg-tlc-storybody").Map(func(_ int, sel *goquery.Selection) string {
			return sel.Text() + "\n"
		}), " ")

		articles := []map[string]interface{}{article}
		data := map[string]interface{}{
			"articles": articles,
		}

		return data, nil
	},
}
var abcExtractor = crawler.FuncExtractor{
	URL: "www.abc.net.au",
	Method: func(doc *goquery.Document) (interface{}, error) {
		result := doc.Find(".news.story_page").Length()
		if result == 0 {
			return nil, nil
		}

		entities := []entity{}
		articles := []map[string]interface{}{}

		data := map[string]interface{}{}
		data["title"] = doc.Find("h1").Text()

		doc.Find(".byline a").Each(func(_ int, sel *goquery.Selection) {
			name := sel.Text()
			link, _ := sel.Attr("href")

			link = "http://www.abc.net.au" + link

			entities = append(entities, entity{
				Name: name,
				Link: link,
				Type: "Person",
			})
		})

		data["authors"] = entities

		data["images"] = doc.Find(".article img").Map(func(_ int, sel *goquery.Selection) string {
			link, _ := sel.Attr("src")
			return link
		})

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

		data["content"] = strings.Join(b, "\n")

		articles = append(articles, data)

		finalData := map[string]interface{}{
			"entities": entities,
			"articles": articles,
		}

		return finalData, nil
	},
}

var sbsExtractor = crawler.FuncExtractor{
	URL: "www.sbs.com.au",
	Method: func(doc *goquery.Document) (interface{}, error) {
		result := doc.Find(".layout--article").Length()
		if result == 0 {
			return nil, nil
		}

		article := map[string]interface{}{}

		article["title"] = doc.Find("h1").Text()

		article["images"] = doc.Find(".article img").Map(func(_ int, sel *goquery.Selection) string {
			link, _ := sel.Attr("src")
			return link
		})

		article["summary"] = doc.Find(".text-abstract__content").Text()

		article["content"] = doc.Find(".text-body").Map(func(_ int, sel *goquery.Selection) string {
			return sel.Text() + "\n"
		})

		articles := []map[string]interface{}{article}
		data := map[string]interface{}{
			"articles": articles,
		}

		return data, nil
	},
}

var ageExtractor = crawler.FuncExtractor{
	URL:    "www.theage.com.au",
	Method: fairfaxEx,
}

var brisbaneTimesExtractor = crawler.FuncExtractor{
	URL:    "www.brisbanetimes.com.au",
	Method: fairfaxEx,
}

var alternetExtractor = crawler.FuncExtractor{
	URL: "www.alternet.org",
	Method: func(doc *goquery.Document) (interface{}, error) {
		entities := []entity{}
		articles := []map[string]interface{}{}

		doc.Find(".byline a").Each(func(_ int, sel *goquery.Selection) {
			name := sel.Text()
			link, _ := sel.Attr("href")

			link = "http://alternet.org" + link

			entities = append(entities, entity{
				Name: name,
				Link: link,
				Type: "Person",
			})
		})

		article := map[string]interface{}{}

		article["title"] = doc.Find("h1").Text()
		article["summary"] = doc.Find("teaser").Text()
		article["authors"] = entities

		article["images"] = doc.Find(".the_body img").Map(func(_ int, sel *goquery.Selection) string {
			link, _ := sel.Attr("src")
			return link
		})

		article["content"] = strings.Join(doc.Find(".the_body").Map(func(_ int, sel *goquery.Selection) string {
			return sel.Text() + "\n"
		}), " ")

		articles = append(articles, article)

		data := map[string]interface{}{
			"entities": entities,
			"articles": articles,
		}

		return data, nil
	},
}

var canberraTimesExtractor = crawler.FuncExtractor{
	URL: "www.canberratimes.com.au",
	Method: func(doc *goquery.Document) (interface{}, error) {
		result := doc.Find(".template--article").Length()
		if result == 0 {
			return nil, nil
		}

		articles := []map[string]interface{}{}

		article := map[string]interface{}{}

		article["title"] = doc.Find("h1").Text()

		article["images"] = doc.Find(".article__body img").Map(func(_ int, sel *goquery.Selection) string {
			link, _ := sel.Attr("src")
			return link
		})

		article["content"] = strings.Join(doc.Find(".article__body").
			Children().
			Not("#subscribe-newsletter").
			Not("figure").
			Not(".adWrapper").
			Map(func(_ int, sel *goquery.Selection) string {
				return sel.Text() + "\n"
			}),
			" ")

		articles = append(articles, article)

		data := map[string]interface{}{
			"articles": articles,
		}

		return data, nil
	},
}

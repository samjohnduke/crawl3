package schedular

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type sitemapIndex struct {
	Sitemap []sitemap `xml:"sitemap"`
}

type sitemap struct {
	Loc     string      `xml:"loc"`
	Lastmod sitemaptime `xml:"lastmod"`
}

type url2 struct {
	Loc             string      `xml:"loc"`
	Lastmod         sitemaptime `xml:"lastmod"`
	ChangeFrequency string      `xml:"changefreq"`
	Priority        float32     `xml:"priority"`
	Altertive       []link      `xml:"link"`
}

type urlSet struct {
	URLS []url2 `xml:"url"`
}

type link struct {
	Lang string `xml:"hreflang,attr"`
	URL  string `xml:"href,attr"`
}

// A SiteMap contains either the index or the urlset
type SiteMap struct {
	SitemapIndex sitemapIndex
	URLset       urlSet
}

type sitemaptime struct {
	time.Time
}

func (c *sitemaptime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02" // yyyymmdd date format
	const longForm = "2006-01-02T15:04:05Z07:00"
	const longForm2 = "2006-01-02T15:04:0507:00"
	const longForm3 = "2006-01-02T15:04Z07:00"
	const longForm4 = "2006-01-02T15:04:05Z0700"

	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(longForm, v)
	if err == nil {
		*c = sitemaptime{parse}
		return nil
	}

	parse, err = time.Parse(longForm2, v)
	if err == nil {
		*c = sitemaptime{parse}
		return nil
	}

	parse, err = time.Parse(longForm3, v)
	if err == nil {
		*c = sitemaptime{parse}
		return nil
	}

	parse, err = time.Parse(longForm4, v)
	if err == nil {
		*c = sitemaptime{parse}
		return nil
	}

	parse, err = time.Parse(shortForm, v)
	if err != nil {
		return err
	}

	*c = sitemaptime{parse}

	return nil
}

// Fill recursively gets all the data from a sitemap as URLS, not matter how many there are,
// This could be dangerous, but its fine for now
// TODO: Limit this in some way so that it caps at a number
func (sm *SiteMap) Fill() error {
	for _, in := range sm.SitemapIndex.Sitemap {
		resp, err := http.Get(in.Loc)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		mm, err := ParseSitemap(out)
		if err != nil {
			return err
		}

		sm.URLset.URLS = append(sm.URLset.URLS, mm.URLset.URLS...)
	}
	return nil
}

// ParseSitemapFromURL first gets the sitemaps data from the web and then parses it.
func ParseSitemapFromURL(url string) (*SiteMap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out []byte
	out, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(out)

	if contentType == "application/x-gzip" {
		var gr io.ReadCloser
		gr, err = gzip.NewReader(bytes.NewBuffer(out))
		defer gr.Close()
		out, err = ioutil.ReadAll(gr)
	}

	return ParseSitemap(out)
}

// ParseSitemap turns an xml byte list into a SiteMap object
func ParseSitemap(sitemap []byte) (*SiteMap, error) {
	var us urlSet
	err := xml.Unmarshal(sitemap, &us)
	if err != nil {
		return nil, err
	}

	if len(us.URLS) > 0 {
		sm := SiteMap{URLset: us}
		return &sm, nil
	}

	var smi sitemapIndex
	err = xml.Unmarshal(sitemap, &smi)
	if err != nil {
		return nil, err
	}

	if len(smi.Sitemap) > 0 {
		sm := SiteMap{SitemapIndex: smi}
		return &sm, nil
	}

	return nil, errors.New("unable to parse sitemap")
}

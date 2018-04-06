package schedular

import (
	"testing"
)

func TestParseSiteMap(t *testing.T) {
	s, err := ParseSitemap([]byte(testUrlset))
	if err != nil {
		t.Error(err)
	}

	if len(s.URLset.URLS) != 2 {
		t.Errorf("Incorrectly Parsed sitemap, %d != 2", len(s.URLset.URLS))
	}

	si, err := ParseSitemap([]byte(testSitemapIndex))
	if err != nil {
		t.Error(err)
	}

	if len(si.SitemapIndex.Sitemap) != 4 {
		t.Errorf("Incorrectly Parsed sitemap, %d != 2", len(si.SitemapIndex.Sitemap))
	}
}

func TestParseSiteMapFill(t *testing.T) {
	si, err := ParseSitemap([]byte(testSitemapIndex))
	if err != nil {
		t.Error(err)
	}

	if len(si.SitemapIndex.Sitemap) != 4 {
		t.Errorf("Incorrectly Parsed sitemap, %d != 2", len(si.SitemapIndex.Sitemap))
	}

	err = si.Fill()
	if err != nil {
		t.Error(err)
	}
}

func TestParseFromURL(t *testing.T) {
	si, err := ParseSitemapFromURL("https://theconversation.com/sitemap.xml")
	if err != nil {
		t.Error(err)
	}

	if len(si.SitemapIndex.Sitemap) <= 0 {
		t.Errorf("Unable to fetch index")
	}
}

var testUrlset = `
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns:xhtml="http://www.w3.org/1999/xhtml" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<url>
    <loc>http://theconversation.com/id</loc>
    <lastmod>2018-03-12T10:21:59Z</lastmod>
    <changefreq>hourly</changefreq>
    <priority>1.0</priority>
    <xhtml:link rel="alternate" hreflang="id" href="http://theconversation.com/id"/>
    <xhtml:link rel="alternate" hreflang="en-ca" href="http://theconversation.com/ca"/>
    <xhtml:link rel="alternate" hreflang="fr" href="http://theconversation.com/fr"/>
    <xhtml:link rel="alternate" hreflang="en" href="http://theconversation.com/global"/>
    <xhtml:link rel="alternate" hreflang="en-us" href="http://theconversation.com/us"/>
    <xhtml:link rel="alternate" hreflang="en-au" href="http://theconversation.com/au"/>
    <xhtml:link rel="alternate" hreflang="en-za" href="http://theconversation.com/africa"/>
    <xhtml:link rel="alternate" hreflang="en-gb" href="http://theconversation.com/uk"/>
    <xhtml:link rel="alternate" hreflang="x-default" href="http://theconversation.com/"/>
</url>
<url>
<loc>http://theconversation.com/ca</loc>
<lastmod>2018-03-13T01:07:06Z</lastmod>
<changefreq>hourly</changefreq>
<priority>1.0</priority>
<xhtml:link rel="alternate" hreflang="id" href="http://theconversation.com/id"/>
<xhtml:link rel="alternate" hreflang="en-ca" href="http://theconversation.com/ca"/>
<xhtml:link rel="alternate" hreflang="fr" href="http://theconversation.com/fr"/>
<xhtml:link rel="alternate" hreflang="en" href="http://theconversation.com/global"/>
<xhtml:link rel="alternate" hreflang="en-us" href="http://theconversation.com/us"/>
<xhtml:link rel="alternate" hreflang="en-au" href="http://theconversation.com/au"/>
<xhtml:link rel="alternate" hreflang="en-za" href="http://theconversation.com/africa"/>
<xhtml:link rel="alternate" hreflang="en-gb" href="http://theconversation.com/uk"/>
<xhtml:link rel="alternate" hreflang="x-default" href="http://theconversation.com/"/>
</url>
</urlset>
`

var testSitemapIndex = `
<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>http://theconversation.com/sitemap_general.xml</loc>
    <lastmod>2018-03-13T00:00:00Z</lastmod>
  </sitemap>
  <sitemap>
    <loc>http://theconversation.com/sitemap_topics_0_9.xml</loc>
    <lastmod>2018-03-13T00:00:00Z</lastmod>
  </sitemap>
  <sitemap>
    <loc>http://theconversation.com/sitemap_topics_a.xml</loc>
    <lastmod>2018-03-13T00:00:00Z</lastmod>
  </sitemap>
  <sitemap>
    <loc>http://theconversation.com/sitemap_topics_b.xml</loc>
    <lastmod>2018-03-13T00:00:00Z</lastmod>
  </sitemap>
</sitemapindex>
`

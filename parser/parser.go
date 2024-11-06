package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func Parse(resp *http.Response) (storage.PageData, []string, error) {

	var data storage.PageData
	var links []string

	if resp.StatusCode != 200 {
		return data, links, fmt.Errorf("Non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return data, links, fmt.Errorf("Error loading body: %s", err)
	}

	data.Title = doc.Find("title").Text()
	data.Desciption, _ = doc.Find("meta[name=description]").Attr("content")

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			link := normalizeURL(resp.Request.URL, href)
			if link != "" {
				links = append(links, link)
			}
		}
	})

	return data, links, nil
}

func normalizeURL(base *url.URL, href string) string {
	u, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return ""
	}
	return base.ResolveReference(u).String()
}

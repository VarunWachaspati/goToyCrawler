package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Parse - Parses the HTML response for URLs to crawl belonging to same domain
func Parse(response, currentURL string) ([]string, []string, error) {
	doc, er1 := goquery.NewDocumentFromReader(strings.NewReader(response))
	if er1 != nil {
		log.Println(er1)
		return nil, nil, er1
	}
	urls := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		url, exists := s.Attr("href")
		if !exists || strings.HasPrefix(url, "#") {
			return false
		}
		url = ResolveToBaseDomain(url, currentURL)
		if IsSameDomain(url, baseDomain) {
			return true
		}
		return false
	}).Map(func(i int, s *goquery.Selection) string {
		url, _ := s.Attr("href")
		return SanitizeURL(ResolveToBaseDomain(url, currentURL))
	})
	assetUrls := doc.Find("img,script,link").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, hrefExists := s.Attr("href")
		src, srcExists := s.Attr("src")
		if srcExists && IsSameDomain(ResolveToBaseDomain(src, currentURL), baseDomain) {
			return true
		} else if hrefExists && IsSameDomain(ResolveToBaseDomain(href, currentURL), baseDomain) {
			return true
		}
		return false
	}).Map(func(i int, s *goquery.Selection) string {
		var result string
		url, exists := s.Attr("href")
		if result, _ = s.Attr("src"); exists {
			result = url
		}
		return SanitizeURL(ResolveToBaseDomain(result, currentURL))
	})
	return urls, assetUrls, nil
}

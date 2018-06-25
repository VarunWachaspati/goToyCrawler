package main

import (
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

// GetBaseDomain - Utility function to get the base domain
func GetBaseDomain(url string) string {
	regex, err := regexp.Compile("(http[s]*:\\/\\/[w]*[\\.]?[a-z0-9]+[\\.][a-z]+)")
	if err != nil {
		log.Println("Failed to compile regex")
	}
	return regex.FindString(url)
}

// ByteSliceToString - Utility function to get String from Byte Slice using reflection
// reference - https://golang.org/pkg/reflect/#SliceHeader
func ByteSliceToString(bArray []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bArray))
	stringHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	stringPointer := (*string)(unsafe.Pointer(&stringHeader))
	return *stringPointer
}

// GetHost - Returns hostname of the url
func GetHost(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return u.Host, nil
}

// IsSameDomain - To check whether the crawled URL belongs to base Domain
func IsSameDomain(url, baseDomain string) bool {
	domain := GetBaseDomain(url)
	return domain == baseDomain
}

// ResolveToBaseDomain - To resolve relative URLs to absolute URL to crawl them properly
func ResolveToBaseDomain(domain, baseDomain string) string {
	link, err := url.Parse(domain)
	if err != nil {
		log.Println(err)
		return domain
	}
	baseLink, err := url.Parse(baseDomain)
	if err != nil {
		log.Println(err)
		return baseDomain
	}
	return SanitizeURL(baseLink.ResolveReference(link).String())
}

// FormatURL - Function to append proper Protocol to URL
func FormatURL(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return link
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		return u.String()
	}
	// Assumption http will work, should try for https also
	u.Scheme = "http"
	return u.String()

}

// SanitizeURL - Cleans up trailing punctutaions and query params
func SanitizeURL(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Println("Error while parsing URL -" + err.Error())
	}
	return strings.TrimRight(u.Scheme+"://"+u.Host+u.Path, "/")
}

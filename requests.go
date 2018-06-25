package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
)

var userAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/66.0.3359.181 Chrome/66.0.3359.181 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
}

//GetPage - Function to make HTTP request and return the response in String
func GetPage(url string) (respBody, finalURL string, err error) {
	timeout := time.Duration(20 * time.Second)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // TODO: set to false in production
		},
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
	req, er1 := http.NewRequest("GET", url, nil)
	if er1 != nil {
		log.Println("Error in Requests [" + url + "] : " + er1.Error())
		return "", "", er1
	}
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])
	res, er2 := client.Do(req)
	if er2 != nil {
		log.Println("Error in Requests [" + url + "] : " + er2.Error())
		return "", "", er2
	}
	defer res.Body.Close()
	// This will be a bottleneck for large responses
	// Explore buffered IO alternatives
	respBArray, er3 := ioutil.ReadAll(res.Body)
	if er3 != nil {
		log.Println("Error in Requests : " + er3.Error())
		return "", res.Request.URL.String(), er3
	}
	return ByteSliceToString(respBArray), SanitizeURL(res.Request.URL.String()), nil
}

func checkIfValidURL(url string) bool {
	if !valid.IsURL(url) {
		return false
	}
	// Instead of GetPage(uses GET) explore HEAD request for performance optimization
	_, _, err := GetPage(FormatURL(url))
	if err != nil {
		return false
	}
	return true
}

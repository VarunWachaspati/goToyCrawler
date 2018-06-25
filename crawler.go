package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/disiqueira/gotree"
)

// Link - Represents the Link/URL to be crawled
type Link struct {
	url          string
	timestamp    time.Time
	depth, tries int
	tree         gotree.Tree
	err          error
}

// DataStore - Interface to represent various datastore like Map, Redis, etc.
type DataStore interface {
	add(url string) bool
	contains(url string) bool
}

func worker(work chan *Link, wg *sync.WaitGroup, cache DataStore) {
	defer wg.Done()
	for {
		select {
		case link := <-work:
			crawl(link, work, cache)
		case <-time.After(1 * time.Second):
			if len(work) == 0 {
				return
			}
		}
	}
}

func newLink(url string, depth int) *Link {
	return &Link{
		url:       url,
		timestamp: time.Now(),
		depth:     depth,
		tries:     0,
		err:       nil,
		tree:      gotree.New(url),
	}
}

// Crawler - Intializes the GoRoutine Workers to crawl the given Domain
func Crawler(url string, cache DataStore) gotree.Tree {
	// Channel on which to be crawled Links will be sent
	in := make(chan *Link, channelSize)
	var wg sync.WaitGroup
	// Spawn Worker Go Routines
	for i := 0; i < runtime.NumCPU()*3; i++ {
		wg.Add(1)
		go worker(in, &wg, cache)
	}
	// Add baseDomain to channel, with base depth 1
	base := newLink(url, 1)
	in <- base
	cache.add(url)
	wg.Wait()
	close(in)
	return base.tree
}

func crawl(url *Link, work chan *Link, cache DataStore) {
	if url.tries > maxRetry {
		return
	}
	if url.depth > *maxDepth {
		return
	}
	// Politeness Delay
	if *delay != 0 {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(*delay)+*delay/2))
	}
	resp, finalURL, err := GetPage(url.url)
	// In case of some error, we should retry the URL
	if err != nil {
		url.err = err
		url.tries = url.tries + 1
		url.timestamp = time.Now()
		work <- url
		log.Println("Error while fetching URL[" + url.url + "] - " + err.Error())
		return
	}
	log.Println("Crawling [ " + url.url + " ] ")
	cache.add(finalURL)
	links, assetLinks, err := Parse(resp, url.url)
	for _, link := range links {
		if !cache.contains(link) {
			childLink := newLink(link, url.depth+1)
			work <- childLink
			cache.add(link)
			url.tree.AddTree(childLink.tree)
		}
	}
	if *staticAssetFlag {
		for _, link := range assetLinks {
			if !cache.contains(link) {
				url.tree.AddTree(gotree.New(link))
				cache.add(link)
			}
		}
	}

}

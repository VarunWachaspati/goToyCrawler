package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/disiqueira/gotree"
)

// Globals which specify the init configuration of the Crawler
var baseDomain string
var maxDepth *int
var delay *int
var staticAssetFlag *bool
var maxRetry = 2
var channelSize = 128

func main() {
	maxDepth = flag.Int("depth", 3, "Specify the depth till which the crawler should crawl till")
	delay = flag.Int("delay", 0, "Specify the delay between each request to a domain")
	staticAssetFlag = flag.Bool("staticAsset", true, "Specify whether static assets should be listed in the sitemap")
	output := flag.String("output", "", "Specify the file to which the sitemap needs to be written. Default is StdOut")
	help := flag.Bool("h", false, "Prints the usage of the command line utility")

	flag.Parse()
	if *help {
		helperMessage()
		return
	}
	// Validations on the flag and Arguments
	if len(flag.Args()) < 1 {
		log.Fatalln("[Error] Missing URL. Enter the URL to be crawled")
	}
	if *maxDepth < 1 {
		log.Fatalln("[Error] Depth cannot be less than 1")
	}
	if !checkIfValidURL(flag.Arg(0)) {
		log.Fatalln("[Error] Kindly recheck the URL if it's working or not")
	}

	// Set Logger
	f, err := os.OpenFile("crawler.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("error opening file: " + err.Error())
		log.Println("Printing Logs to stdout")
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	baseDomain = GetBaseDomain(FormatURL(flag.Arg(0)))
	localCache := &LocalStore{store: make(map[string]bool)}
	log.Println("**********************************************")
	fmt.Println("Started Crawling : " + flag.Arg(0))
	log.Println("Started Crawling : " + flag.Arg(0))
	sitemap := Crawler(baseDomain, localCache)
	printSiteMap(sitemap, flag.Arg(0), *output)
	log.Println("Crawling Finished for Domain : " + flag.Arg(0))
	log.Println("**********************************************")
}

func helperMessage() {
	fmt.Println("GoCrawler Utility Usage - ./crawler [-opt val] <URL>")
	fmt.Println("Command Line flags - ")
	fmt.Println("1. -depth <level> : Specifies the level of depth till which the crawler should crawl from the given URL\n" +
		"                    Deafult Level : 3. Level cannot be less than 1.")
	fmt.Println("2. -delay <milliseconds> : Specify the politeness delay between each request to the domain\n" +
		"                           Deafult Delay : 0")
	fmt.Println("3. -staticAsset <boolean> : Specify specify whether static assets should be listed on the sitemap or not.\n" +
		"                            Deafult : true")
	fmt.Println("4. -output <file> : Specify the file to which the sitemap needs to be written\n" +
		"                    Deafult : StdOut")
	fmt.Println("5. -h : Prints the commandline utility usages.")
}

func printSiteMap(sitemap gotree.Tree, url, output string) {
	var w io.Writer
	if w = os.Stdout; output != "" {
		f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Unable to Print to output file, Writing to StdOut")
		} else {
			w = io.MultiWriter(f, os.Stdout)
		}
	}
	fmt.Fprintln(w, "SiteMap for the URL : "+url)
	fmt.Fprintln(w, "==============================================")
	fmt.Fprintln(w, sitemap.Print())
	fmt.Fprintln(w, "==============================================")
}

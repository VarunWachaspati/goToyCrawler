# GoLang Crawler
A Toy Web Crawler which crawls over the given domain and lists down the linkages between different webpages and the static assets each page uses via a sitemap (printed in a tree format). 

Crawler doesn't crawl external links, only the links belonging to the same domain.

#### Assumptions
- Assumes Standard installation of GoLang 1.8.3
- Tested on amd64 Linux

## Usage

#### Generic
```
> make install
> make test
> make build
> ./crawler -depth <Level> -delay <Milliseconds> -output <File> <URL>
> cat <File>
```
##### Commandline Flags
| Flag           | Default  | Description  |
| -------------- |:--------:| :------------|
| -depth \<Level>| 3 | Specifies the level of depth till which the crawler should crawl from the given URL |
| -delay \<milliseconds>| 0 | Specify the politeness delay between each request to the domain |
| -staticAsset \<boolean>| true | Specify specify whether static assets should be listed on the sitemap or not |
| -output \<file>| StdOut | Specify the file to which the sitemap needs to be written |
| -h| - | Prints the commandline utility usages |
#### Sample
```
> ./crawler -output sitemap.txt tomblomfield.com
> cat sitemap.txt
```
#### Help
```
> ./crawler -h
```

## Design Decisions
1. Each URL is fetched by a separate GoRoutine as fetching is an IO intensive operation, hence can leverage GoLang concurrency model to increase performance.
2. Spawns only ```runtime.NumCPU() * 3``` Worker GoRoutines, to have more control over parallel requests and avoid memory allocation issues when crawling a huge website which might arise if each URL spawned it's own GoRoutine.
3. Set timeout of 20 seconds in HTTP Client. The default in ```net/http``` package is no timeout which can lead to unresponsive behavior.
4. Used generic interface ```DataStore``` which can be implemented later by a persistent/ in-memory data store like SQLite/Redis instead of the ```LocalStore```(Map with Mutex) used right now.
5. ByteArray to string conversion done using ```reflect``` and ```unsafe``` packages by changing ```sliceHeader``` to ```stringHeader```. In comparison to ```string(byte[])```, this methodology would require no copy of data, hence more performant. 
```
func ByteSliceToString(bArray []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bArray))
	stringHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	stringPointer := (*string)(unsafe.Pointer(&stringHeader))
	return *stringPointer
}
```  
6. As all requests are to the same domain, added politeness delay by sleeping the goroutine for some random time in a range [```time.Sleep(time.Millisecond * time.Duration(rand.Intn(*delay)+*delay/2))```] set using ```-delay <Milliseconds>``` flag.
If nothing specified, then by default no delay is induced.

## Issues
1. Set ```InsecureSkipVerify: true``` in ```TLSClientConfig``` which skips TLS verification as of now.
2. Cannot parse dynamically generated webpages as of now. Can be done by injecting a javascript engine/ headless chrome object and render dynamic content.
3. Cannot resolve few URL which eventually resolve to another URL like
	1. [tomblomfield.com/page/1](tomblomfield.com/page/1) being the same as [tomblomfield.com](tomblomfield.com)
	2. [https://www.sitemaps.org/index.php](https://www.sitemaps.org/index.php) resolves to [https://www.sitemaps.org/index.html](https://www.sitemaps.org/index.html) and [https://www.sitemaps.org](https://www.sitemaps.org)
4. Lack of Integration Tests.
5. The User Experience is not great as there is no progress bar showing the status of crawler.
6. Design not extensible to build a distributed crawler.
7. Makefile not up to standards, and should add dependency management package like [dep](https://github.com/golang/dep). Currently listed all external dependencies in the Makefile.

## Dependencies
1. [gotree](github.com/disiqueira/gotree) - Used for printing Sitemap in tree format
2. [goquery](github.com/PuerkitoBio/goquery) - Used for parsing HTML pages
3. [govalidator](github.com/asaskevich/govalidator) - Used for validating URL
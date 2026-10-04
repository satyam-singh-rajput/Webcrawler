package Helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	gocache "github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	"github.com/gocolly/colly"
	goCache "github.com/patrickmn/go-cache"
)

func getLinks(url string) []string {
	c := colly.NewCollector()
	var links []string
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
		fmt.Println(links)
	})
	c.Visit(url)
	return links
}

func worker(jobs chan string, results chan<- map[string][]string, cxt context.Context,
	cacheManager *gocache.Cache[string]) {
	for url := range jobs {
		fmt.Println("processing this Url: ", url)
		_, err := cacheManager.Get(cxt, url)
		if err == nil {
			fmt.Println("this Url is already in the cache")
		} else {
			links := getLinks(url)
			for _, link := range links {
				jobs <- link
			}
			cacheManager.Set(cxt, url, "Visited")
			results <- map[string][]string{
				url: links,
			}
		}
	}
}
func serialize(mp map[string][]string) {

	var buildTree func(url string, path map[string]bool) map[string]any

	buildTree = func(url string, path map[string]bool) map[string]any {
		if path[url] {
			return map[string]any{}
		}
		path[url] = true
		node := make(map[string]any)
		for _, child := range mp[url] {
			node[child] = buildTree(child, path)
		}
		delete(path, url)
		return node
	}
	root := make(map[string]any)
	for url := range mp {
		root[url] = buildTree(url, make(map[string]bool))
	}
	res, err := os.Create("CrawlerResult.json")
	if err != nil {
		fmt.Println("Issue while creating CrawlerResult.json", err)
	}
	defer res.Close()
	encoder := json.NewEncoder(res)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(root)
	if err != nil {
		fmt.Println("Error serializing CrawlerResult.json", err)
	}
}
func scraper(url string) {
	cxt := context.Background()
	//q := createQueue()
	cacheClient := goCache.New(5*time.Minute, 10*time.Minute)
	cacheStore := gocache_store.NewGoCache(cacheClient)
	cacheManger := gocache.New[string](cacheStore)
	jobs := make(chan string)
	results := make(chan map[string][]string)

	for i := 1; i <= 10; i++ {
		go worker(jobs, results, cxt, cacheManger)
	}
	jobs <- url
	close(results)
	mp := make(map[string][]string)
	for result := range results {
		for url, links := range result {
			mp[url] = links
		}
	}
	close(jobs)
	serialize(mp)
}

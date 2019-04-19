package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"sync/atomic"
)

const startID = 1000000
const path = "https://goodster.ru"

var ID = int64(startID)

type parseUrl struct {
	url   string
	level int
}

func main() {
	startCategories := getStartCategories()
	pararellParser(startCategories)
}

func getID() int64 {
	return atomic.AddInt64(&ID, int64(1))
}

func parseCatalogPage(url parseUrl) []parseUrl {
	var answer []parseUrl

	doc, err := goquery.NewDocument(path + url.url)
	if err != nil {
		panic(err)
	}

	doc.Find("li.categoriesListBoxTitle").Each(func(i int, s *goquery.Selection) {
		id := getID()
		node := s.Find("a").Eq(0)
		href, ok := node.Attr("href")
		if !ok {
			panic("qwe")
		}
		title, err := node.Find("b").Eq(0).Html()
		if err != nil {
			panic(err)
		}
		answer = append(answer, parseUrl{
			url:   href,
			level: url.level + 1,
		})
		fmt.Printf(`{id: %d, title: "%s", href: "%s", lvl: %d}`, id, title, path+href, url.level+1)
		fmt.Println()
	})

	return answer
}

func pararellParser(urls []parseUrl) {
	wg := sync.WaitGroup{}
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url parseUrl) {
			pararellParser(parseCatalogPage(url))
			wg.Done()
		}(url)
	}
	wg.Wait()
}

func getStartCategories() []parseUrl {
	var answer []parseUrl

	//+ "/tools/ajax/mainNav.php"
	doc, err := goquery.NewDocument(path)
	if err != nil {
		panic(err)
	}

	doc.Find(".mainNav .mainNavLevel1").Each(func(i int, s *goquery.Selection) {
		id := getID()
		node := s.Find("a").Eq(0)
		href, ok := node.Attr("href")
		if !ok {
			panic("qwe")
		}
		if href == "#" {
			//	s.Find(".mainNavTitle>a").Each(func(i int, s *goquery.Selection) {
			//		href, ok := s.Attr("href")
			//		if !ok {
			//			panic("qwe")
			//		}
			//		title, err := s.Html()
			//		if err != nil {
			//			panic(err)
			//		}
			//
			//		answer = append(answer, parseUrl{
			//			url: href,
			//			level: 1,
			//		})
			//
			//		fmt.Printf(`{id: %d, title: "%s", href: "%s", lvl: %d}`, id, title, path + href, 1)
			//		fmt.Println()
			//	})
			return
		}

		title, err := node.Html()
		if err != nil {
			panic(err)
		}

		answer = append(answer, parseUrl{
			url:   href,
			level: 1,
		})

		fmt.Printf(`{id: %d, title: "%s", href: "%s", lvl: %d}`, id, title, path+href, 1)
		fmt.Println()
	})

	return answer
}

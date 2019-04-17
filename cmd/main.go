package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
)

func main() {
	url := "https://goodster.ru/cat/61803/"
	fmt.Println("request: " + url)

	//fmt.Println(parseCatalogPage(url))
	pararellParser([]string{url})
}

func parseCatalogPage(url string) []string {
	var answer []string

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("li.categoriesListBoxTitle").Each(func(i int, s *goquery.Selection) {
		node := s.Find("a").Eq(0)
		href, ok := node.Attr("href")
		if !ok {
			panic("qwe")
		}
		title, err := node.Find("b").Eq(0).Html()
		if err != nil {
			panic(err)
		}
		answer = append(answer, href)
		fmt.Println(title, href)
	})

	return answer
}

func pararellParser(urls []string) {
	wg := sync.WaitGroup{}
	wg.Add(len(urls))
	for _, url := range urls {
		go func() {
			pararellParser(parseCatalogPage(url))
			wg.Done()
		}()
	}
	wg.Wait()
}

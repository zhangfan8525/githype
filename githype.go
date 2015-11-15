package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	GitHubUrlStub            = "http://github.com"
	TrendingUrlStub          = "https://github.com/trending"
	RepoNameHtmlClass        = "repo-list-name"
	RepoDescriptionHtmlClass = "repo-list-description"
)

const Usage = `
  USAGE
      githype [options]
  OPTIONS
      -l, -lang 
          Set the programming language to search for
      -s, -since
          Set the trending time span, daily, weekly, monthly
      -t, -top
          Set the number of items ot display
`

func main() {
	lang := flag.String("l", "", "programming language")
	flag.StringVar(lang, "lang", "", "programming language")
	since := flag.String("s", "", "trending since")
	flag.StringVar(since, "since", "", "trending since")
	top := flag.Int("t", 25, "items to display")
	flag.IntVar(top, "top", 25, "items to display")
	flag.Usage = func() { fmt.Println(Usage) }
	flag.Parse()

	url := fmt.Sprintf("%s?l=%s&since=%s", TrendingUrlStub, *lang, *since)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	nameNextLine := false
	descNextLine := false
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	cnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		if nameNextLine {
			cnt++
			s := strings.Split(line, "\"")
			fmt.Printf(" %d. %s\n", cnt, GitHubUrlStub+s[1])
			nameNextLine = false
		}
		if descNextLine {
			fmt.Printf("%s\n\n", line)
			descNextLine = false
			if cnt >= *top {
				break
			}
		}
		if strings.Contains(line, RepoNameHtmlClass) {
			nameNextLine = true
		}
		if strings.Contains(line, RepoDescriptionHtmlClass) {
			descNextLine = true
		}
	}
}

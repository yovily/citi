package main

import (
	"fmt"
	"net/http"
)

type Endpoint struct {
	Urls []string
}

func main() {
	ep := Endpoint{
		Urls: []string{
			"https://www.google.com",
			"https://www.facebook.com",
			"https://www.twitter.com",
			"https://www.instagram.com",
			"https://www.linkedin.com",
			"https://www.youtube.com",
			"https://www.github.com",
		},
	}
	c := make(chan string)

	for _, url := range ep.Urls {
		go checkUrl(url, c)
	}
	for range ep.Urls {
		fmt.Println(<- c)
	}


}

func checkUrl(url string, channel chan string) {
	_, err := http.Get(url)
	if err != nil {
		channel <- fmt.Sprintf("%s is down", url)
	}
	channel <- fmt.Sprintf("%s yep, its up", url)
}

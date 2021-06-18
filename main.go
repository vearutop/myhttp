package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vearutop/myhttp/internal"
)

func main() {
	c := internal.Fetcher{
		OnError: func(err error, link string) {
			fmt.Println(link, err.Error())
		},
		OnSuccess: func(hash, link string) {
			fmt.Println(link, hash)
		},
	}

	flag.IntVar(&c.Concurrency, "parallel", 10, "maximum number of concurrent requests")
	flag.Parse()

	c.Links = flag.Args()
	if len(c.Links) == 0 {
		flag.Usage()

		return
	}

	c.Fetch(context.Background())
}

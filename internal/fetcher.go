// Package internal provides implementation details that are not part of public API.
package internal

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Fetcher is a concurrent HTTP resource fetcher.
type Fetcher struct {
	Concurrency int
	Links       []string
	OnError     func(err error, link string)
	OnSuccess   func(hash, link string)
}

// Fetch downloads and hashes a batch of HTTP resources.
func (c Fetcher) Fetch(ctx context.Context) {
	if c.Concurrency == 0 {
		c.Concurrency = 10
	}

	semaphore := make(chan struct{}, c.Concurrency)

	for _, l := range c.Links {
		l := l

		semaphore <- struct{}{}

		go func() {
			defer func() {
				<-semaphore
			}()

			if err := c.do(ctx, l); err != nil {
				if c.OnError != nil {
					c.OnError(err, l)
				}
			}
		}()
	}

	// Wait for jobs to finish by filling semaphore to full capacity.
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}

	close(semaphore)
}

func (c Fetcher) do(ctx context.Context, l string) (err error) {
	if !strings.HasPrefix(l, "http://") && !strings.HasPrefix(l, "https://") {
		l = "http://" + l
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// If there are any redirects, they will not be followed.
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	defer func() {
		clErr := resp.Body.Close()
		if err == nil && clErr != nil {
			err = clErr
		}
	}()

	h := md5.New()

	_, err = io.Copy(h, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.OnSuccess != nil {
		c.OnSuccess(fmt.Sprintf("%x", h.Sum(nil)), l)
	}

	return nil
}

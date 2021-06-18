package internal_test

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/vearutop/myhttp/internal"
)

func TestFetcher_Fetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(r.URL.String()))
		if err != nil {
			t.Fatal(err)
		}
	}))

	defer srv.Close()

	f := internal.Fetcher{
		OnError: func(err error, link string) {
			t.Fatal(err)
		},
		OnSuccess: func(hash, link string) {
			h := fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimPrefix(link, srv.URL))))
			if h != hash {
				t.Fatalf("hash mismatch, expected %s, received %s", h, hash)
			}
		},
	}
	for i := 0; i < 100; i++ {
		f.Links = append(f.Links, strings.TrimPrefix(srv.URL, "http://")+"/"+strconv.Itoa(i))
	}

	f.Fetch(context.Background())
}

func BenchmarkFetcher_Fetch(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(r.URL.String()))
		if err != nil {
			b.Fatal(err)
		}
	}))

	defer srv.Close()

	f := internal.Fetcher{
		OnError: func(err error, link string) {
			b.Fatal(err)
		},
		OnSuccess: func(hash, link string) {
			h := fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimPrefix(link, srv.URL))))
			if h != hash {
				b.Fatalf("hash mismatch, expected %s, received %s", h, hash)
			}
		},
	}

	b.ReportAllocs()

	batch := 0

	for i := 0; i < b.N; i++ {
		f.Links = append(f.Links, strings.TrimPrefix(srv.URL, "http://")+"/"+strconv.Itoa(i))

		batch++

		if batch >= 100 {
			f.Fetch(context.Background())

			batch = 0
			f.Links = f.Links[0:0]
		}
	}

	f.Fetch(context.Background())
}

package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/heretic1321/pokedex/internal/store"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)

// Minimal shape of the API response (reuse your real type if you prefer).
type areaResp struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []singleArea  `json:"results"`
}
type singleArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func TestServiceCache_MissThenHit_SameURL(t *testing.T) {
	var hits int64

	// Fake PokeAPI server with pagination.
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)

		if !strings.HasPrefix(r.URL.Path, "/api/v2/location-area/") {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		offset := atoiDefault(q.Get("offset"), 0)
		limit := atoiDefault(q.Get("limit"), 20)

		next := fmt.Sprintf("%s/api/v2/location-area/?offset=%d&limit=%d", srv.URL, offset+limit, limit)
		prev := fmt.Sprintf("%s/api/v2/location-area/?offset=%d&limit=%d", srv.URL, max0(offset-limit), limit)

		// normalize prev to offset=0 for first page (so it equals the first URL)
		if offset == 0 {
			prev = fmt.Sprintf("%s/api/v2/location-area/?offset=0&limit=%d", srv.URL, limit)
		}

		var results []singleArea
		for i := 0; i < limit; i++ {
			results = append(results, singleArea{
				Name: fmt.Sprintf("area-%d", offset+i),
				URL:  fmt.Sprintf("%s/area/%d", srv.URL, offset+i),
			})
		}
		resp := areaResp{
			Count:    9999,
			Next:     strPtr(next),
			Previous: strPtr(prev),
			Results:  results,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	// Real client against our fake server.
	httpClient := srv.Client()
	api := pokeapi.New(httpClient, srv.URL+"/api/v2/")
	cache := store.New(10 * time.Second) // no eviction during test
	svc := New(api, cache)

	// Use explicit offset/limit so the key stays stable across map/mapb.
	firstURL := api.Base + "location-area/?offset=0&limit=20"

	// 1) MISS -> fetch + cache
	var page1 areaResp
	if err := svc.FetchResults(firstURL, &page1); err != nil {
		t.Fatalf("FetchResults(1) error: %v", err)
	}
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Fatalf("expected 1 server hit, got %d", got)
	}

	// 2) HIT -> served from cache (no additional server hit)
	var page1again areaResp
	if err := svc.FetchResults(firstURL, &page1again); err != nil {
		t.Fatalf("FetchResults(1 again) error: %v", err)
	}
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Fatalf("expected cache hit (still 1 server hit), got %d", got)
	}
	if page1.Results[0].Name != page1again.Results[0].Name {
		t.Fatalf("expected same cached content")
	}

	// 3) mapb scenario: call with Previous (which we normalized to the same string)
	if page1.Previous == nil {
		t.Fatalf("previous must not be nil on first page")
	}
	var prevPage areaResp
	if err := svc.FetchResults(*page1.Previous, &prevPage); err != nil {
		t.Fatalf("FetchResults(previous) error: %v", err)
	}
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Fatalf("expected previous to be cache hit as well, got server hits=%d", got)
	}
}

func TestServiceCache_DifferentButEquivalentURLs_MissThenMissWithoutNormalization(t *testing.T) {
	var hits int64
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		// Return an empty but valid page
		resp := areaResp{
			Count:    1,
			Next:     nil,
			Previous: nil,
			Results:  []singleArea{{Name: "area-0", URL: srv.URL + "/area/0"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	httpClient := srv.Client()
	api := pokeapi.New(httpClient, srv.URL+"/api/v2/")
	cache := store.New(10 * time.Second)
	svc := New(api, cache)

	// Two equivalent URLs that differ as strings:
	rawA := api.Base + "location-area/"                 // no query
	rawB := api.Base + "location-area/?offset=0&limit=20"

	var a areaResp
	if err := svc.FetchResults(rawA, &a); err != nil {
		t.Fatalf("FetchResults(rawA) error: %v", err)
	}
	if want, got := int64(1), atomic.LoadInt64(&hits); got != want {
		t.Fatalf("want hits=%d got=%d", want, got)
	}

	var b areaResp
	if err := svc.FetchResults(rawB, &b); err != nil {
		t.Fatalf("FetchResults(rawB) error: %v", err)
	}
	// Because keys differ, this will be a MISS again (hits=2) unless you normalize keys.
	if want, got := int64(2), atomic.LoadInt64(&hits); got != want {
		t.Fatalf("expected a miss due to different cache keys, hits=%d", got)
	}
}

func TestCache_EvictsAfterInterval(t *testing.T) {
	c := store.New(50 * time.Millisecond)
	key := "k"
	c.Add(key, []byte("v"))

	// Immediately present
	if b, ok := c.Get(key); !ok || string(b) != "v" {
		t.Fatalf("expected cached value")
	}

	// Wait past interval (readLoop evicts > interval old)
	time.Sleep(120 * time.Millisecond)

	if _, ok := c.Get(key); ok {
		t.Fatalf("expected value to be evicted")
	}
}

// --- helpers ---

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}
func max0(x int) int {
	if x < 0 {
		return 0
	}
	return x
}
func strPtr(s string) *string { return &s }


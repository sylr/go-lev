package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/sylr/go-lev/rand"
	qdsync "sylr.dev/libqd/sync"
)

type randomResponse struct {
	Tests int64                     `json:"tests"`
	Total int64                     `json:"total"`
	Time  string                    `json:"time"`
	Hists map[string]map[string]int `json:"hits"`
}

func stopped(w http.ResponseWriter, r *http.Request) bool {
	if Stopped {
		http.NotFound(w, r)
		return true
	}

	return false
}

func httpGetRandom(w http.ResponseWriter, r *http.Request) {
	if stopped(w, r) {
		return
	}

	count := 500
	max := 22
	httpParams := r.URL.Query()

	if c, ok := httpParams["count"]; ok {
		count, _ = strconv.Atoi(c[0])
	}

	if m, ok := httpParams["max"]; ok {
		max, _ = strconv.Atoi(m[0])
	}

	resp := randomResponse{}
	hashes := rand.GetRandomHashSlice(count)

	start := time.Now()
	wg := qdsync.NewBoundedWaitGroup(runtime.NumCPU())
	mu := sync.Mutex{}
	resp.Hists = make(map[string]map[string]int)

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int) {
			for j := i + 1; j < count; j++ {
				lev := levenshtein.ComputeDistance(hashes[i], hashes[j])

				if lev < max {
					mu.Lock()
					resp.Total++
					if _, ok := resp.Hists[hashes[i]]; !ok {
						resp.Hists[hashes[i]] = make(map[string]int)
					}
					resp.Hists[hashes[i]][hashes[j]] = lev
					mu.Unlock()
				}
			}

			atomic.AddInt64(&resp.Tests, int64(count-(i+1)))

			wg.Done()
		}(i)
	}

	wg.Wait()

	resp.Time = fmt.Sprintf("%f secs", time.Since(start).Seconds())

	goLevDistanceResultsHistogram.WithLabelValues(strconv.Itoa(count), strconv.Itoa(max)).Observe(float64(resp.Total))

	if _, ok := httpParams["json"]; ok {
		w.Header().Add("Content-Type", "application/json")
		b, _ := json.Marshal(resp)
		w.Write(b)
	} else {
		w.Header().Add("Content-Type", "text")
		for k, v := range resp.Hists {
			for i, j := range v {
				fmt.Fprintf(w, "%s to %s = %d\n", k, i, j)
			}
		}

		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Tests: %12d\n", resp.Tests)
		fmt.Fprintf(w, "Found: %12d\n", resp.Total)
		fmt.Fprintf(w, "Time: %13s\n", resp.Time)
	}
}

func httpGetDistance(w http.ResponseWriter, r *http.Request) {
	if stopped(w, r) {
		return
	}

	var strs []string
	var ok bool
	httpParams := r.URL.Query()

	w.Header().Add("Content-Type", "text")

	if strs, ok = httpParams["string"]; !ok || len(strs) != 2 {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Bad parameters\n")
		return
	}

	distance := levenshtein.ComputeDistance(strs[0], strs[1])

	fmt.Fprintf(w, "%d\n", distance)
}

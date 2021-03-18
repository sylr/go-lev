package v1

import (
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

	tests := int64(0)
	count := 500
	max := 22
	httpParams := r.URL.Query()

	w.Header().Add("Content-Type", "text")

	if c, ok := httpParams["count"]; ok {
		count, _ = strconv.Atoi(c[0])
	}

	if c, ok := httpParams["max"]; ok {
		max, _ = strconv.Atoi(c[0])
	}

	hashes := rand.GetRandomHashSlice(count)

	start := time.Now()
	wg := qdsync.NewBoundedWaitGroup(runtime.NumCPU())
	mu := sync.Mutex{}
	total := 0

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int) {
			for j := i + 1; j < count; j++ {
				lev := levenshtein.ComputeDistance(hashes[i], hashes[j])

				if lev < max {
					mu.Lock()
					total++
					fmt.Fprintf(w, "%s to %s = %d\n", hashes[i], hashes[j], lev)
					mu.Unlock()
				}
			}

			atomic.AddInt64(&tests, int64(count-(i+1)))

			wg.Done()
		}(i)
	}

	wg.Wait()

	end := time.Since(start)

	goLevDistanceResultsHistogram.WithLabelValues(strconv.Itoa(count), strconv.Itoa(max)).Observe(float64(total))

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Tests: %12d\n", tests)
	fmt.Fprintf(w, "Found: %12d\n", total)
	fmt.Fprintf(w, "Time: %13f secs\n", end.Seconds())
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

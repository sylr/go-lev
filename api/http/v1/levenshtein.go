package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/sylr/go-lev/rand"
)

func httpGetRandom(w http.ResponseWriter, r *http.Request) {
	if stopped(w, r) {
		return
	}

	tests := 0
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

	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			lev := levenshtein.ComputeDistance(hashes[i], hashes[j])

			if lev < max {
				fmt.Fprintf(w, "%s to %s = %d\n", hashes[i], hashes[j], lev)
			}

			tests++
		}
	}

	end := time.Since(start)

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Tests: %12d\n", tests)
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

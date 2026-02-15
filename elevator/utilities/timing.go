package utilities

import "time"

// TimeN runs fn n times and returns total duration and average duration
func TimeN(n int, fn func()) (total time.Duration, avg time.Duration) {
	start := time.Now()

	for i := 0; i < n; i++ {
		fn()
	}

	total = time.Since(start)
	avg = total / time.Duration(n)
	return
}

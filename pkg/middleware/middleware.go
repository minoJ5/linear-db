package middleware

import (
	"fmt"
	"net/http"
	"runtime"
)

func Memory(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var start, end runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&start)
		next.ServeHTTP(w, r)
		runtime.ReadMemStats(&end)
		total := float64(end.TotalAlloc-start.TotalAlloc) / (1024 * 1024)
		malloc := float64(end.Mallocs-start.Mallocs) / (1024 * 1024)
		current := float64(end.Alloc) / (1024 * 1024)
		fmt.Printf("Memory: Used [ %.7f MB ] with mallocs [ %.7f MB ] now [ %.7f MB ]\n", total, malloc, current)
	})
}

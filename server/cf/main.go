package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	host := GetEnvOrDefault("HOST", "0.0.0.0")
	port := GetEnvOrDefault("PORT", "8080")
	server := fmt.Sprintf("%s:%s", host, port)
	print := fmt.Sprintf("Running on %s", server)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(print))
	})
	fmt.Println(print)
	http.ListenAndServe(server, nil)
}

func GetEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

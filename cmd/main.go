package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func main() {
	// Get the port and destination URL from environment variables or default values
	port := getEnv("PROXY_PORT", "4199")
	baseURL := getEnv("BASE_URL", "https://my.homeassistant.io")

	// Set up the HTTP server and handle incoming requests
	http.HandleFunc("/", handleRequest(baseURL))

	log.Printf("Starting proxy server on port %s, forwarding to HA at %s", port, baseURL)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRequest(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Construct the full destination URL
		destURL, err := url.Parse(baseURL)
		if err != nil {
			log.Printf("Error parsing base URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logRequest(r, http.StatusInternalServerError, start, "")
			return
		}
		destURL.Path = path.Join(destURL.Path, r.URL.Path)

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			logRequest(r, http.StatusMethodNotAllowed, start, "")
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form data: %v", err)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			logRequest(r, http.StatusBadRequest, start, "")
			return
		}

		// Extract solar radiation
		solarRadiation := r.FormValue("solarradiation")

		// Forward the request to the destination URL
		resp, err := http.PostForm(destURL.String(), r.Form)
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			http.Error(w, "Failed to forward request", http.StatusBadGateway)
			logRequest(r, http.StatusBadGateway, start, solarRadiation)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		// Copy response from the destination to the client
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error copying response: %v", err)
		}

		logRequest(r, resp.StatusCode, start, solarRadiation)
	}
}

func logRequest(r *http.Request, status int, start time.Time, solarRadiation string) {
	duration := time.Since(start)
	log.Printf(
		"Time: %s | IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v | Solar Radiation: %s",
		start.Format(time.RFC3339),
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
		status,
		duration,
		solarRadiation,
	)
}

// getEnv fetches the environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

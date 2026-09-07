package pkg

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func VerifyRepo(repo Repo) bool {
	// 1. Basic non-empty field validations
	if repo.Name == "" || repo.Type == "" || repo.OS == "" || repo.Arch == "" || repo.LocalDir == "" {
		return false
	}

	if !repo.Enabled {
		return false
	}

	// 2. Universal OS validation
	osTarget := strings.ToLower(repo.OS)
	if !supportedOS[osTarget] {
		return false
	}

	// 3. Universal Architecture validation
	archTarget := strings.ToLower(repo.Arch)
	if !supportedArch[archTarget] {
		return false
	}

	// 4. Validate local path presence and integrity
	if !checkLocalResource(repo.LocalDir) {
		return false
	}

	// 5. Validate remote URL if provided
	if repo.URL != "" {
		if !checkRemoteResource(repo.URL) {
			return false
		}
	}

	return true
}

func checkLocalResource(path string) bool {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return false
	}
	if !info.IsDir() && info.Size() == 0 {
		return false
	}
	return true
}

func checkRemoteResource(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(rawURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return true
	}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Range", "bytes=0-0") // Fetch only 1 byte to minimize bandwidth
	resp, err = client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Install Repo downloads, unpacks, or places binaries/libraries from repo.URL into repo.LocalDir.
func Install(repo Repo) error {
	// 1. Validate repository metadata
	if !VerifyRepo(repo) {
		return fmt.Errorf("invalid repository configuration: %+v", repo)
	}

	// 2. Ensure target local directory exists
	if err := os.MkdirAll(repo.LocalDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", repo.LocalDir, err)
	}

	// If no remote URL is given, assume local setup and return early
	if repo.URL == "" {
		return nil
	}

	// 3. Download remote artifact to temporary file
	tmpFile, err := os.CreateTemp("", "repo-install-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if err := downloadFile(repo.URL, tmpFile); err != nil {
		return fmt.Errorf("failed to download resource from %s: %w", repo.URL, err)
	}

	// Rewind file pointer for extraction/copying
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// 4. Extract or copy file based on extension
	lowerURL := strings.ToLower(repo.URL)
	switch {
	case strings.HasSuffix(lowerURL, ".tar.gz") || strings.HasSuffix(lowerURL, ".tgz"):
		if err := extractTarGz(tmpFile, repo.LocalDir); err != nil {
			return fmt.Errorf("failed to extract tar.gz archive: %w", err)
		}

	case strings.HasSuffix(lowerURL, ".zip"):
		// Zip reader requires file size info
		fi, err := tmpFile.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat temp file for zip extraction: %w", err)
		}
		if err := extractZip(tmpFile, fi.Size(), repo.LocalDir); err != nil {
			return fmt.Errorf("failed to extract zip archive: %w", err)
		}

	default:
		// Raw binary, DLL, or .so file — copy directly into LocalDir
		fileName := filepath.Base(repo.URL)
		destPath := filepath.Join(repo.LocalDir, fileName)

		destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, tmpFile); err != nil {
			return fmt.Errorf("failed to write raw binary to %s: %w", destPath, err)
		}
	}

	return nil
}

// downloadFile performs an HTTP GET download with timeout limits.
func downloadFile(rawURL string, dst *os.File) error {
	client := http.Client{
		Timeout: 10 * time.Minute, // Timeout for large binary downloads
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	_, err = io.Copy(dst, resp.Body)
	return err
}

// extractTarGz extracts a .tar.gz stream into targetDir with path traversal protection (Zip Slip).
func extractTarGz(r io.Reader, targetDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath, err := safeJoin(targetDir, header.Name)
		if err != nil {
			continue // Skip malicious or invalid archive paths
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

// extractZip extracts a .zip archive into targetDir with path traversal protection.
func extractZip(ra io.ReaderAt, size int64, targetDir string) error {
	zr, err := zip.NewReader(ra, size)
	if err != nil {
		return err
	}

	for _, file := range zr.File {
		targetPath, err := safeJoin(targetDir, file.Name)
		if err != nil {
			continue // Skip path traversal vulnerabilities
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(f, rc)
		f.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// safeJoin protects against path traversal vulnerabilities (Zip Slip / Tar Slip).
func safeJoin(baseDir, relativePath string) (string, error) {
	cleanPath := filepath.Clean(relativePath)
	joined := filepath.Join(baseDir, cleanPath)

	// Ensure destination remains within base directory boundary
	if !strings.HasPrefix(joined, filepath.Clean(baseDir)+string(filepath.Separator)) && joined != filepath.Clean(baseDir) {
		return "", fmt.Errorf("illegal file path: %s", relativePath)
	}
	return joined, nil
}

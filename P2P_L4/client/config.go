package main

import (
	_ "embed"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed icon.png
var TrayIcon []byte

var ServerURL = "ssltun.fehmicorp.in"

type MenuOptions struct {
	Title   string
	Type    string
	Tag     string
	Dynamic bool
	Func    func(*application.Context)
}

type AppConfig struct {
	AppName         string
	Description     string
	Icon            string
	Version         string
	Domain          string
	InstallationDir string `json:"installDir"`
}

// GetExpandedInstallDir safely resolves environment variables in system paths
func (a AppConfig) GetExpandedInstallDir() string {
	expanded := os.ExpandEnv(a.InstallationDir)
	return filepath.Clean(expanded)
}

var HttpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// Ensure idle connections drop quickly when network routes shift
		IdleConnTimeout: 15 * time.Second,
	},
}

var Conf = AppConfig{
	AppName:         "Shalimar SAP Agent",
	Description:     "Shalimar Private SAP access Agent",
	Icon:            "assets/icon.png",
	Version:         "v1.0.1",
	Domain:          "shalimarcorp.org",
	InstallationDir: `%ProgramFiles%\shalimar\agent`, // Windows standard environment format
}

var Target = struct {
	OS   string
	Arch string
}{
	OS:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

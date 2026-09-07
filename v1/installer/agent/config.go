package main

import (
	"os"
	"runtime"

	"github.com/fehmicorp/pkg/v1/utils/os/win"
)

var Target = struct {
	OS       string
	Arch     string
	Hostname string
}{
	OS:       runtime.GOOS,
	Arch:     runtime.GOARCH,
	Hostname: "",
}

type AppConfig struct {
	AppName         string
	Description     string
	Icon            string
	Version         string
	Domain          string
	InstallationDir string `json:"installDir"`
}

var Conf = AppConfig{
	AppName:         "Fehmi Agent Installer",
	Description:     "Fehmi Agent Installer is a web-based application that allows users to easily install and manage Fehmi agents on their systems. It provides a user-friendly interface for configuring and deploying agents, making it simple for users to monitor and control their cloud infrastructure.",
	Icon:            "assets/logo.png",
	Version:         "v1.0.1",
	Domain:          "fehmicorp.in",
	InstallationDir: `%ProgramFiles%\fehmi\agent`,
}

type SystemPayload struct {
	TargetOS   string `json:"targetOs"`
	TargetArch string `json:"targetArch"`
	IsAdmin    bool   `json:"isAdmin"`
	Hostname   string `json:"hostname"`
	AgentVer   string `json:"agentVersion"`
}

type SystemService struct{}

func (s *SystemService) GetSystemInfo() SystemPayload {
	Target.Hostname, _ = os.Hostname()
	isAdmin := win.IsAdmin()
	return SystemPayload{
		TargetOS:   Target.OS,
		TargetArch: Target.Arch,
		IsAdmin:    isAdmin,
		Hostname:   Target.Hostname,
		AgentVer:   Conf.Version,
	}
}

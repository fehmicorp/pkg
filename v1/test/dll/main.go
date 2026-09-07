package main

import "C" // REQUIRED for c-shared exports

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/fehmicorp/pkg/v1/utils/os/win"
	"github.com/fehmicorp/pkg/v1/utils/os/win/tray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var InstanceId string = "FCAgent_systray"
var app *application.App

var Conf = struct {
	AppName     string
	Description string
}{
	AppName:     "Fehmi Agent",
	Description: "Version 1.0.1",
}

//go:embed icon.png
var Icon []byte

func main() {}

var Menu = []tray.MenuOptions{
	{
		Title:  "SAP Console",
		Type:   "window",
		Target: "https://google.com",
	},
}

//export Run
func Run() {
	_ = win.EnsureSingleInstance(InstanceId, nil, func() { os.Exit(0) })
	defer win.ReleaseSingleInstance()

	app = application.New(application.Options{
		Name:        Conf.AppName,
		Description: Conf.Description,
	})

	tooltip := fmt.Sprintf("%s\n%s", Conf.AppName, Conf.Description)
	onLeftClick := func() {
		tray.OpenBrowser(Menu[0].Target)
	}
	onRightClick := func() {}

	err := tray.RunMenu(app, Icon, tooltip, Menu, true, onLeftClick, onRightClick)
	if err != nil {
		fmt.Printf("Error starting tray: %v\n", err)
		os.Exit(1)
	}
}

//export Quit
func Quit() {
	if app != nil {
		app.Quit()
	}
}

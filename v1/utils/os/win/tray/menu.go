package tray

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type MenuOptions struct {
	Title  string
	Type   string // "web" or "browser"
	Target string // URL or route path
}

func PrepareMenuItems(app *application.App, menu []MenuOptions) []MenuItemConfig {
	var items []MenuItemConfig
	for _, opt := range menu {
		items = append(items, MenuItemConfig{
			Title:   opt.Title,
			OnClick: getFuncForMenu(app, opt),
		})
	}
	return items
}

func getFuncForMenu(app *application.App, option MenuOptions) func(ctx *application.Context) {
	switch option.Type {
	case "window":
		return func(ctx *application.Context) {
			OpenWindow(app, option.Title, option.Target)
		}
	case "browser":
		return func(ctx *application.Context) {
			OpenBrowser(option.Target)
		}
	case "openapp":
		return func(ctx *application.Context) {
			OpenApp(option.Target)
		}
	default:
		return nil
	}
}

func OpenApp(targetExe string, args ...string) error {
	cmd := exec.Command(targetExe, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch %s: %w", targetExe, err)
	}
	return nil
}

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func OpenWindow(app *application.App, title string, url string) {
	if app == nil {
		log.Println("Cannot open window: app instance is nil")
		return
	}

	// 1. Check if window already exists using app.Window.GetByName
	if win, ok := app.Window.GetByName(title); ok && win != nil {
		win.Show()
		win.Focus()
		return
	}

	// 2. Create window using app.Window.NewWithOptions or app.NewWebviewWindow
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   title,
		Title:  title,
		Width:  1024,
		Height: 768,
		URL:    url,
	})

	window.Show()
}

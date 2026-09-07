package tray

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type TrayManager struct {
	App  *application.App
	Tray *application.SystemTray
}

type MenuItemConfig struct {
	Title   string
	OnClick func(ctx *application.Context)
}

type TrayConfig struct {
	IconData     []byte
	Tooltip      string
	CustomMenus  []MenuOptions
	ShowQuit     bool
	OnLeftClick  func()
	OnRightClick func()
}

// NewTrayManager initializes the system tray with icon data, tooltips, and custom menus
func NewTrayManager(app *application.App, cfg *TrayConfig) *TrayManager {
	menuItems := PrepareMenuItems(app, cfg.CustomMenus)
	systray := app.SystemTray.New()
	if len(cfg.IconData) > 0 {
		systray.SetIcon(cfg.IconData)
	}
	if cfg.Tooltip != "" {
		systray.SetTooltip(cfg.Tooltip)
	}
	tm := &TrayManager{
		App:  app,
		Tray: systray,
	}
	menu := tm.BuildMenu(menuItems, cfg.ShowQuit)
	systray.SetMenu(menu)
	if cfg.OnLeftClick != nil {
		systray.OnDoubleClick(func() {
			cfg.OnLeftClick()
		})
	}
	if cfg.OnRightClick != nil {
		systray.OnRightClick(func() {
			cfg.OnRightClick()
			systray.OpenMenu()
		})
	}
	return tm
}
func (tm *TrayManager) BuildMenu(customMenus []MenuItemConfig, showQuit bool) *application.Menu {
	menu := application.NewMenu()

	for _, item := range customMenus {
		title := item.Title
		handler := item.OnClick
		menu.Add(title).OnClick(func(ctx *application.Context) {
			if handler != nil {
				handler(ctx)
			}
		})
	}

	if showQuit {
		menu.AddSeparator()
		menu.Add("Quit").OnClick(func(ctx *application.Context) {
			tm.App.Quit()
		})
	}

	return menu
}

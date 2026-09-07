package tray

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func RunMenu(
	app *application.App,
	TrayIcon []byte,
	Tooltip string,
	menu []MenuOptions,
	Quit bool,
	onLeftClick func(),
	onRightClick func(),
) (*application.SystemTray, error) {
	var conf = TrayConfig{
		IconData:     TrayIcon,
		Tooltip:      Tooltip,
		CustomMenus:  menu,
		ShowQuit:     Quit,
		OnLeftClick:  onLeftClick,
		OnRightClick: onRightClick,
	}
	systray := NewTrayManager(app, &conf)
	if err := app.Run(); err != nil {
		return nil, fmt.Errorf("❌ System tray service failed: %v", err)
	}
	return systray.Tray, nil
}

func SetupTray(
	app *application.App,
	TrayIcon []byte,
	Tooltip string,
	menu []MenuOptions,
	Quit bool,
	onLeftClick func(),
	onRightClick func(),
) *application.SystemTray {
	conf := TrayConfig{
		IconData:     TrayIcon,
		Tooltip:      Tooltip,
		CustomMenus:  menu,
		ShowQuit:     Quit,
		OnLeftClick:  onLeftClick,
		OnRightClick: onRightClick,
	}

	tm := NewTrayManager(app, &conf)
	return tm.Tray
}

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
) error {
	var conf = TrayConfig{
		IconData:     TrayIcon,
		Tooltip:      Tooltip,
		CustomMenus:  menu,
		ShowQuit:     Quit,
		OnLeftClick:  onLeftClick,
		OnRightClick: onRightClick,
	}
	NewTrayManager(app, &conf)
	if err := app.Run(); err != nil {
		return fmt.Errorf("❌ System tray service failed: %v", err)
	}
	return nil
}

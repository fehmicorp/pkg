package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	isConnected int32              // 0 = Disconnected, 1 = Connected
	cancelFunc  context.CancelFunc // Used to terminate the running client session
)

func PrepareMenuItems() []systray.MenuItemConfig {
	var customMenus []systray.MenuItemConfig

	// 1. Static Menu Items (SAP Console, etc.)
	menuList := []MenuOptions{
		{
			Title: "SAP Console",
			Type:  "web",
			Func: func(ctx *application.Context) {
				OpenBrowser("https://vhrpedclcc01.sap.shalimarcorp.org:8443")
			},
		},
	}

	for _, menu := range menuList {
		currentMenu := menu
		itemConfig := systray.MenuItemConfig{
			Title: currentMenu.Title,
			OnClick: func(ctx *application.Context) {
				if currentMenu.Func != nil {
					currentMenu.Func(ctx)
				}
			},
		}
		customMenus = append(customMenus, itemConfig)
	}

	// 2. Dynamic VPN Connect / Disconnect Action Item
	var vpnMenuItem *systray.MenuItemConfig

	vpnMenuItem = &systray.MenuItemConfig{
		Title: "Connect",
		OnClick: func(appCtx *application.Context) {
			// DISCONNECT ACTION: If running, cancel the tunnel session
			if atomic.LoadInt32(&isConnected) == 1 {
				if cancelFunc != nil {
					cancelFunc()
				}
				atomic.StoreInt32(&isConnected, 0)
				vpnMenuItem.Title = "Connect"
				return
			}

			// CONNECT ACTION: Thread-safe single execution check
			if atomic.CompareAndSwapInt32(&isConnected, 0, 1) {
				vpnMenuItem.Title = "Disconnect"

				// Spawn tunnel execution off the UI thread
				go func() {
					defer func() {
						atomic.StoreInt32(&isConnected, 0)
						vpnMenuItem.Title = "Connect"
						if HttpClient != nil {
							HttpClient.CloseIdleConnections()
						}
					}()

					sigCtx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
					defer stopSignal()

					var clientCtx context.Context
					clientCtx, cancelFunc = context.WithCancel(sigCtx)

					runTunnel(clientCtx)
				}()
			}
		},
	}

	customMenus = append(customMenus, *vpnMenuItem)
	return customMenus
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

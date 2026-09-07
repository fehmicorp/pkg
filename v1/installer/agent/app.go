package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) GetConfig() AppConfig {
	return Conf
}
func (a *App) WindowMinimize() {
	runtime.WindowMinimise(a.ctx)
}
func (a *App) WindowToggleMaximize() {
	runtime.WindowToggleMaximise(a.ctx)
}
func (a *App) WindowClose() {
	runtime.Quit(a.ctx)
}

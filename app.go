package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"simUI/config"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	config.Ctx = ctx

	//做一些配置
	runtime.WindowSetDarkTheme(ctx) //设为黑色主题
}

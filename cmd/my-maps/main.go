package main

import (
	"github.com/oneworkspacestudio/my-maps/internal/app"
	"github.com/oneworkspacestudio/my-maps/internal/webui"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"os"
)

var version = "0.1.0-preview.5"

func main() {
	if e := apphost.Main(apphost.Config{ID: "maps", Name: "My Maps", Version: version, Listen: "127.0.0.1:9015", Geolocation: true, MediaCapture: true}, webui.Files(), func(h *apphost.Host) (apphost.Application, error) { return app.New(h) }); e != nil {
		apphost.Report(e)
		os.Exit(1)
	}
}

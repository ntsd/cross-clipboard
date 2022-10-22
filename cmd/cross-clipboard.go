package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/ntsd/cross-clipboard/ui"
)

func main() {
	isTerminalMode := flag.Bool("t", false, "run in terminal mode")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	crossClipboard, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if isTerminalMode != nil && *isTerminalMode {
		for {
			select {
			case l := <-crossClipboard.LogChan:
				log.Println("log: ", l)
			case err := <-crossClipboard.ErrorChan:
				var fatalErr *xerror.FatalError
				if errors.As(err, &fatalErr) {
					log.Fatal(fmt.Errorf("fatal error: %w", fatalErr))
				}
				log.Println(fmt.Errorf("runtime error: %w", err))
			case clipboards := <-crossClipboard.ClipboardManager.ClipboardsChannel:
				_ = clipboards
			case devices := <-crossClipboard.DeviceManager.DevicesChannel:
				for _, dv := range devices {
					if dv.Status == device.StatusPending {
						fmt.Printf("device %s wanted to connect (Y/n)", dv.Name)
						var input string
						fmt.Scanln(&input)
						if input == "n" {
							dv.Block()
						} else {
							err = dv.Trust()
							if err != nil {
								log.Println(fmt.Errorf("can not trust device: %w", err))
							}
						}
						crossClipboard.DeviceManager.UpdateDevice(dv)
					}
				}
			}
		}
	} else {
		view := ui.NewView(crossClipboard)
		view.Start()
	}
}

package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/mobile"
	"github.com/ntsd/cross-clipboard/mobile/gui"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// cfg := &config.Config{
	// 	GroupName:  "default",
	// 	ListenHost: "0.0.0.0",
	// 	ListenPort: 4001,
	// }
	crossclipboard.NewCrossClipboard(cfg)

	g := gui.NewGUI()
	mobile.SetGame(g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}

package clipboard

import (
	"context"

	"golang.design/x/clipboard"
)

func main() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		println(string(data))
		// calling event to send the clipboard data to the server
	}
}

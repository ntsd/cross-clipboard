package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	WINDOW_NAME = "Cross Clipbaord"
)

func main() {
	a := app.New()
	w := a.NewWindow(WINDOW_NAME)

	hello := widget.NewLabel("Hello Fyne!")
	vBox := container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	)
	vScroll := container.NewVScroll(vBox)

	w.SetContent(vScroll)

	w.ShowAndRun()
}

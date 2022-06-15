package ui

import "github.com/rivo/tview"

type Page func(prevPage func(), nextPage func()) (title string, content tview.Primitive)

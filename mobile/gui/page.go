package gui

import (
	"github.com/ebitenui/ebitenui/widget"
)

type page struct {
	title   string
	content widget.PreferredSizeLocateableWidget
}

func clipboardsPage(res *uiResources) *page {
	c := newPageContentContainer()
	return &page{
		title:   "Clipboards",
		content: c,
	}
}

func devicesPage(res *uiResources) *page {
	c := newPageContentContainer()
	return &page{
		title:   "Devices",
		content: c,
	}
}

func settingPage(res *uiResources) *page {
	c := newPageContentContainer()
	return &page{
		title:   "Setting",
		content: c,
	}
}

func logPage(res *uiResources) *page {
	c := newPageContentContainer()
	return &page{
		title:   "Log",
		content: c,
	}
}

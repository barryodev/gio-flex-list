package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"log"
	"math/rand"
	"os"
	"time"
)

// UI holds all of the application state.
type UI struct {
	// Theme is used to hold the fonts used throughout the application.
	theme *material.Theme

	firstList *layout.List
	feeds []*feed

	secondList *layout.List
	entries []*entry

	textBox *widget.Editor
}

type feed struct {
	name string
	url string
}

type entry struct {
	title string
	url string
	contents string
}

var defaultMargin = unit.Dp(10)

func main() {
	ui := &UI{}
	ui.theme = material.NewTheme(gofont.Collection())

	ui.feeds = createDummyFeeds()
	ui.entries = createDummyEntries()

	ui.firstList = &layout.List{
		Axis: layout.Vertical,
	}

	ui.secondList = &layout.List{
		Axis: layout.Vertical,
	}

	ui.textBox = new(widget.Editor)
	ui.textBox.SetText(ui.entries[0].contents)

	go func() {
		w := app.NewWindow(
			app.Title("Counter"),
			app.Size(unit.Dp(540), unit.Dp(350)),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}


func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			}

		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	// inset is used to add padding around the window border.
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.flexed(gtx)
	})
}

func (ui *UI) flexed(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			return ui.layoutFeeds(gtx)
		}),
		layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
					return ui.layoutEntries(gtx)
				}),
				layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
					return material.Editor(ui.theme, ui.textBox, "Hint").Layout(gtx)
				}),
			)

		}),

	)

}


func (ui *UI) layoutFeeds(gtx layout.Context) layout.Dimensions {
	return ui.firstList.Layout(gtx, len(ui.feeds), func(gtx layout.Context, i int) layout.Dimensions {
		in := layout.UniformInset(unit.Dp(3))
		return in.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			feedNameWidget := material.Body1(ui.theme, ui.feeds[i].name)
			feedNameWidget.MaxLines = 1
			return feedNameWidget.Layout(gtx)
		})
	})
}

func (ui *UI) layoutEntries(gtx layout.Context) layout.Dimensions {
	return ui.secondList.Layout(gtx, len(ui.entries), func(gtx layout.Context, i int) layout.Dimensions {
		in := layout.UniformInset(unit.Dp(3))
		return in.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			entryNameWidget := material.Body1(ui.theme, ui.entries[i].title)
			entryNameWidget.MaxLines = 1
			return entryNameWidget.Layout(gtx)
		})
	})
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString() string {
	max := 10
	min := 5
	return StringWithCharset(rand.Intn(max - min) + min, charset)
}

func createDummyEntries() []*entry {
	dummyEntries := make([]*entry, 10)

	for i := 0; i < 10; i++ {
		fakeEntry := entry{title: RandomString(),
			url: "http://www." + RandomString() + ".com/" + RandomString(),
			contents: createDummyContents()}

		dummyEntries[i] = &fakeEntry
	}

	return dummyEntries
}

func createDummyContents() string {
	min := 30
	max := 60

	totalWords := rand.Intn(max - min) + min

	dummyContents := ""
	for i := 0; i < totalWords; i++ {
		dummyContents += RandomString() + " "
	}

	return dummyContents
}

func createDummyFeeds() []*feed {
	dummyFeeds := make([]*feed, 10)

	for i := 0; i < 10; i++ {
		fakeFeed := feed{name: RandomString(), url: "http://www." + RandomString() + ".com/" + RandomString() }
		dummyFeeds[i] = &fakeFeed
	}

	return dummyFeeds
}
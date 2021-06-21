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

const dummyContents string = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Nulla malesuada pellentesque elit eget gravida cum. Ultricies mi quis hendrerit dolor. Lectus quam id leo in vitae turpis massa sed elementum. Luctus venenatis lectus magna fringilla urna. Sollicitudin ac orci phasellus egestas tellus rutrum tellus. Eget est lorem ipsum dolor sit amet consectetur. Adipiscing tristique risus nec feugiat in fermentum posuere urna nec. Sed elementum tempus egestas sed sed risus pretium quam. Eget velit aliquet sagittis id consectetur purus ut faucibus. Risus sed vulputate odio ut enim blandit volutpat maecenas volutpat. Lacus vel facilisis volutpat est velit egestas. Massa sed elementum tempus egestas sed sed. Blandit volutpat maecenas volutpat blandit.

Quam nulla porttitor massa id neque aliquam vestibulum morbi blandit. Est velit egestas dui id ornare arcu odio. Eget nulla facilisi etiam dignissim diam. Condimentum mattis pellentesque id nibh. Suspendisse sed nisi lacus sed viverra tellus. Ut tellus elementum sagittis vitae et leo duis ut diam. Neque gravida in fermentum et sollicitudin ac orci phasellus. Ultricies mi quis hendrerit dolor magna eget est lorem. Sit amet mauris commodo quis imperdiet. Ut pharetra sit amet aliquam id diam. Volutpat est velit egestas dui id ornare arcu odio ut. Mauris vitae ultricies leo integer malesuada nunc. Ut sem viverra aliquet eget sit amet tellus cras. Odio tempor orci dapibus ultrices in iaculis nunc sed augue. Vestibulum rhoncus est pellentesque elit.
`

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

func RandomString(numberOfWords, lengthOfWord int) string {
	var randomWords string

	for i := 0; i < numberOfWords; i++ {
		randomWords += StringWithCharset(lengthOfWord, charset)

		if i != numberOfWords - 1 {
			randomWords += " "
		}
	}

	return randomWords
}

func createDummyEntries() []*entry {
	dummyEntries := make([]*entry, 20)

	for i := 0; i < 20; i++ {
		fakeEntry := entry{title: RandomString(5, 7),
			url: "http://www." + RandomString(1, 20) + ".com/" + RandomString(1, 10),
			contents: dummyContents}

		dummyEntries[i] = &fakeEntry
	}

	return dummyEntries
}


func createDummyFeeds() []*feed {
	dummyFeeds := make([]*feed, 50)

	for i := 0; i < 50; i++ {
		fakeFeed := feed{name: RandomString(3, 6), url: "http://www." + RandomString(1, 20) + ".com/" + RandomString(1, 10) }
		dummyFeeds[i] = &fakeFeed
	}

	return dummyFeeds
}
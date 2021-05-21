package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"os"
)

// UI holds all of the application state.
type UI struct {
	// Theme is used to hold the fonts used throughout the application.
	Theme *material.Theme
}

var defaultMargin = unit.Dp(10)

func main() {
	ui := &UI{}
	ui.Theme = material.NewTheme(gofont.Collection())

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
		return flexed(gtx)
	})
}
/*
// Layout lays out the counter and handles input.
func (counter *Counter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// Flex layout lays out widgets from left to right by default.
	return layout.Flex{}.Layout(gtx,
		// We use weight 1 for both text and count to make them the same size.
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// We center align the text to the area available.
			return layout.Center.Layout(gtx,
				// Body1 is the default text size for reading.
				material.Body1(th, strconv.Itoa(counter.Count)).Layout)
		}),
		// We use an empty widget to add spacing between the text
		// and the button.
		layout.Rigid(layout.Spacer{Height: defaultMargin}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// For every click on the button increment the count.
			for range counter.increase.Clicks() {
				counter.Count++
			}
			// Finally display the button.
			return material.Button(th, &counter.increase, "Count").Layout(gtx)
		}),
	)
}*/


func flexed(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			return ColorBox(gtx, gtx.Constraints.Min, red)
		}),
		layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
					return ColorBox(gtx, gtx.Constraints.Min, green)
				}),
				layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
					return ColorBox(gtx, gtx.Constraints.Min, blue)
				}),
			)
		}),
	)
}



// Test colors.
var (
	//background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

// ColorBox creates a widget with the specified dimensions and color.
func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
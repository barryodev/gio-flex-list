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
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"os"
)

// UI holds all of the application state.
type UI struct {
	// Theme is used to hold the fonts used throughout the application.
	theme *material.Theme

	firstList *layout.List
	secondList *layout.List
	textBox *widget.Editor
}

var defaultMargin = unit.Dp(10)

func main() {
	ui := &UI{}
	ui.theme = material.NewTheme(gofont.Collection())

	ui.firstList = &layout.List{
		Axis: layout.Vertical,
	}

	ui.secondList = &layout.List{
		Axis: layout.Vertical,
	}

	ui.textBox = new(widget.Editor)
	ui.textBox.SetText(longText)

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
			return ColorBox(gtx, gtx.Constraints.Min, red)
		}),
		layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
					return ColorBox(gtx, gtx.Constraints.Min, green)
				}),
				layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
					return material.Editor(ui.theme, ui.textBox, "Hint").Layout(gtx)
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

const longText = `1. I learned from my grandfather, Verus, to use good manners, and to put restraint on anger. 

2. In the famous memory of my father I had a pattern of modesty and manliness. 

3. Of my mother I learned to be pious and generous; to keep myself not only from evil deeds, but even from evil thoughts; and to live with a simplicity which is far from customary among the rich. 

4. I owe it to my great-grandfather that I did not attend public lectures and discussions, but had good and able teachers at home; and I owe him also the knowledge that for things of this nature a man should count no expense too great.

5. My tutor taught me not to favour either green or blue at the chariot races, nor, in the contests of gladiators, to be a supporter either of light or heavy armed. He taught me also to endure labour; not to need many things; to serve myself without troubling others; not to intermeddle in the affairs of others, and not easily to listen to slanders against them.

6. Of Diognetus I had the lesson not to busy myself about vain things; not to credit the great professions of such as pretend to work wonders, or of sorcerers about their charms, and their expelling of Demons and the like; not to keep quails (for fighting or divination), nor to run after such things; to suffer freedom of speech in others, and to apply myself heartily to philosophy. Him also I must thank for my hearing first Bacchius, then Tandasis and Marcianus; that I wrote dialogues in my youth, and took a liking to the philosopher's pallet and skins, and to the other things which, by the Grecian discipline, belong to that profession.

7. To Rusticus I owe my first apprehensions that my nature needed reform and cure; and that I did not fall into the ambition of the common Sophists, either by composing speculative writings or by declaiming harangues of exhortation in public; further, that I never strove to be admired by ostentation of great patience in an ascetic life, or by display of activity and application; that I gave over the study of rhetoric, poetry, and the graces of language; and that I did not pace my house in my senatorial robes, or practise any similar affectation. I observed also the simplicity of style in his letters, particularly in that which he wrote to my mother from Sinuessa. I learned from him to be easily appeased, and to be readily reconciled with those who had displeased me or given cause of offence, so soon as they inclined to make their peace; to read with care; not to rest satisfied with a slight and superficial knowledge; nor quickly to assent to great talkers. I have him to thank that I met with the discourses of Epictetus, which he furnished me from his own library.

8. From Apollonius I learned true liberty, and tenacity of purpose; to regard nothing else, even in the smallest degree, but reason always; and always to remain unaltered in the agonies of pain, in the losses of children, or in long diseases. He afforded me a living example of how the same man can, upon occasion, be most yielding and most inflexible. He was patient in exposition; and, as might well be seen, esteemed his fine skill and ability in teaching others the principles of philosophy as the least of his endowments. It was from him that I learned how to receive from friends what are thought favours without seeming humbled by the giver or insensible to the gift.`

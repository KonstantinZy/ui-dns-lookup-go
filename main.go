package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const WINDOWHEIGHT = 400
const WINDOWWIDTH = 800
const WINDOWBORDER = 3

func main() {
	myApp := app.New()
	w := myApp.NewWindow("DNS lookup app")

	makeUI(&w)

	w.Resize(fyne.NewSize(WINDOWWIDTH, WINDOWHEIGHT))

	w.ShowAndRun()
}

func makeUI(w *fyne.Window) {

	/* main container */
	borderRightCanvas := canvas.NewRectangle(color.Transparent)
	borderBottomCanvas := canvas.NewRectangle(color.Transparent)
	borderLeftCanvas := canvas.NewRectangle(color.Transparent)
	borderBottomCanvas.SetMinSize(fyne.NewSize(WINDOWWIDTH, WINDOWBORDER))
	borderRightCanvas.SetMinSize(fyne.NewSize(WINDOWBORDER, WINDOWHEIGHT))
	borderLeftCanvas.SetMinSize(fyne.NewSize(WINDOWBORDER, WINDOWHEIGHT))
	mainCont := container.New(layout.NewBorderLayout(nil, borderBottomCanvas, borderLeftCanvas, borderRightCanvas), borderBottomCanvas, borderLeftCanvas, borderRightCanvas)

	/* addres container */
	addrCont := container.NewHBox()
	addrCont.Add(canvas.NewText("Label", color.Opaque))
	addrCont.Add(canvas.NewText("Addres", color.Opaque))
	addrCont.Add(canvas.NewText("Button", color.Opaque))

	/* tabs container sith results */
	resContWrap := container.NewMax()
	resCont := container.NewAppTabs(
		container.NewTabItem("MX", widget.NewButton("MX", func() { fmt.Println("Tapped") })),
		container.NewTabItem("A", widget.NewLabel("A")),
		container.NewTabItem("NS", widget.NewLabel("NS")),
	)
	resCont.SetTabLocation(container.TabLocationLeading)

	resContWrap.Add(resCont)
	gridCont := container.NewBorder(addrCont, nil, nil, nil, resContWrap)

	/* set UI */
	mainCont.Add(gridCont)

	(*w).SetContent(mainCont)
}

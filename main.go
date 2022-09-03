package main

import (
	"errors"
	"image/color"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const WINDOWHEIGHT = 400
const WINDOWWIDTH = 800
const WINDOWBORDER = 3

var (
	addres    *widget.Entry
	searchBtn *widget.Button
	progres   *widget.ProgressBarInfinite
)

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

	/* addres + search container */
	label := canvas.NewText("Addres for lookup:", color.Opaque)

	addres = widget.NewEntry()
	addres.Validator = func(val string) error {
		re := regexp.MustCompile(`^(https*:\/\/)*[^\.]+?\.\w{2,}$`)
		if !re.MatchString(val) {
			return errors.New("wrong adres format")
		}
		return nil
	}
	addres.SetOnValidationChanged(func(err error) {
		if err == nil {
			searchBtn.Enable()
		} else {
			searchBtn.Disable()
		}
	})
	addres.SetPlaceHolder("addres")

	controlCont := container.NewHBox()
	searchBtn = widget.NewButtonWithIcon("Lookup", theme.SearchIcon(), func() {
		progres.Show()
		progres.Start()
		lookupAddr()
	})
	searchBtn.Disable()
	controlCont.Add(searchBtn)
	progres = widget.NewProgressBarInfinite()
	progres.Hide()
	controlCont.Add(progres)

	addrCont := container.NewBorder(nil, nil, label, controlCont, addres)

	/* tabs container sith results */
	resContWrap := container.NewMax()
	resCont := container.NewAppTabs(
		container.NewTabItem("MX", widget.NewLabel("MX")),
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

func lookupAddr() {
	time.Sleep(time.Second * 5)
	progres.Stop()
	progres.Hide()
}

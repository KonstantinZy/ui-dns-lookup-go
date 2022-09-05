package main

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/KonstantinZy/ui-dns-lookup-go/tabs"
	entry "github.com/KonstantinZy/ui-dns-lookup-go/ui"
)

const WINDOWHEIGHT = 400
const WINDOWWIDTH = 800
const WINDOWBORDER = 3

var (
	addres           *entry.DateEntry
	searchBtn        *widget.Button
	progres          *widget.ProgressBarInfinite
	state            = binding.NewString()
	stopLookupSignal = make(chan struct{})
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("DNS lookup app")

	makeUI(&w)

	w.Resize(fyne.NewSize(WINDOWWIDTH, WINDOWHEIGHT))

	w.ShowAndRun()
}

func makeUI(w *fyne.Window) {

	state.Set("Input addres for lookup")

	/* main container */
	borderRightCanvas := canvas.NewRectangle(color.Transparent)
	borderBottomCanvas := canvas.NewRectangle(color.Transparent)
	borderLeftCanvas := canvas.NewRectangle(color.Transparent)
	borderBottomCanvas.SetMinSize(fyne.NewSize(WINDOWWIDTH, WINDOWBORDER))
	borderRightCanvas.SetMinSize(fyne.NewSize(WINDOWBORDER, WINDOWHEIGHT))
	borderLeftCanvas.SetMinSize(fyne.NewSize(WINDOWBORDER, WINDOWHEIGHT))
	mainCont := container.New(layout.NewBorderLayout(nil, borderBottomCanvas, borderLeftCanvas, borderRightCanvas), borderBottomCanvas, borderLeftCanvas, borderRightCanvas)

	/* search line */
	label := canvas.NewText("Addres for lookup:", color.Opaque)

	/* ----------- addres entry field ----------------- */
	addres = entry.NewDateEntry()
	addres.Validator = func(val string) error {
		re := regexp.MustCompile(`^[\w\.]+\.[a-z]{2,}$`)
		if !re.MatchString(val) {
			return errors.New("wrong addres format")
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

	addres.SetOnFocusLost(func() {
		ValidateAndWriteStatus(addres)
	})
	addres.SetOnEnter(func() {
		if err := ValidateAndWriteStatus(addres); err == nil {
			startLookup()
		}
	})

	addres.SetPlaceHolder("addres")

	controlCont := container.NewHBox()
	searchBtn = widget.NewButtonWithIcon("Lookup", theme.SearchIcon(), func() {
		startLookup()
	})
	searchBtn.Disable()
	controlCont.Add(searchBtn)
	progres = widget.NewProgressBarInfinite()
	progres.Hide()
	controlCont.Add(progres)

	addrCont := container.NewBorder(nil, nil, label, controlCont, addres)

	/* tabs container sith results */
	resCont := container.NewAppTabs()
	resCont.SetTabLocation(container.TabLocationLeading)

	tabsForApp := tabs.Get()
	for _, tab := range tabsForApp {
		resCont.Append(
			container.NewTabItem(
				tab.Name,
				widget.NewListWithData(
					tab.Records,
					tab.CreateItem,
					tab.UpdateItem,
				),
			),
		)
	}

	stateLine := widget.NewLabelWithData(state)

	resContWrap := container.NewBorder(nil, stateLine, nil, nil, resCont)

	gridCont := container.NewBorder(addrCont, nil, nil, nil, resContWrap)

	/* set UI */
	mainCont.Add(gridCont)

	(*w).SetContent(mainCont)
}

// ------------ function with main work logic ------------------
func startLookup() {
	lockInputAndShowProgres()

	// ------ progres in state line --------
	go func(line *binding.String, signal *chan struct{}) {
		t := 0
	outer:
		for {
			select {
			case <-*signal:
				break outer
			default:
				(*line).Set(fmt.Sprintf("Runing (%d seconds past)", t))
				time.Sleep(time.Second)
				t++
			}
		}
	}(&state, &stopLookupSignal)

	// --------------------  main work here -----------------------

	wg := &sync.WaitGroup{}

	tabsForApp := tabs.Get()
	wg.Add(len(tabsForApp))
	for _, tab := range tabsForApp {
		go func(rec *binding.ExternalStringList, getRes func(string) []string, wg *sync.WaitGroup) {
			defer wg.Done()
			(*rec).Set(getRes(addres.Text))
		}(&tab.Records, tab.GetResult, wg)
	}

	wg.Wait()

	// ------------------ END main work here END ------------------

	// finish lookup process - reset UI state
	unlockInputAndHideProgres()
	stopLookupSignal <- struct{}{}
	state.Set(fmt.Sprintf("Lookup finished for %s", addres.Text))
}

func lockInputAndShowProgres() {
	progres.Show()
	progres.Start()
	addres.Disable()
	searchBtn.Disable()
}

func unlockInputAndHideProgres() {
	progres.Stop()
	progres.Hide()
	addres.Enable()
	searchBtn.Enable()
}

func ValidateAndWriteStatus(e *entry.DateEntry) error {
	err := e.Validate()
	if err != nil {
		state.Set(fmt.Sprint("Error in text field: ", err))
		return errors.New("field error")
	} else {
		state.Set("Press lookup button for start")
	}

	return nil
}

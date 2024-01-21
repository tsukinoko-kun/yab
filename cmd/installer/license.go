package main

import (
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
)

func License(w fyne.Window) {
	log.Info("Showing license")
	w.SetIcon(theme.DocumentIcon())

	licOk := false

	installBtn := widget.NewButton("Install", func() {
		if licOk {
			log.Info("Installing...")
			Install(w)
		} else {
			log.Error("You must agree to the license!")
		}
	})
	installBtn.Disable()

	w.SetContent(container.NewVBox(
		widget.NewLabel("Yab License"),
		widget.NewLabel(license()),
		widget.NewCheck("I agree", func(value bool) {
			licOk = value
			if value {
				installBtn.Enable()
			} else {
				installBtn.Disable()
			}
		}),
		installBtn,
	))
}

func license() string {
	url := "https://raw.githubusercontent.com/Frank-Mayer/yab/main/LICENSE"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

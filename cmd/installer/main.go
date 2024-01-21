package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Yab :: Yet Another Buildtool")

	License(w)

	w.ShowAndRun()
}

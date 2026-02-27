package exwidgets_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	exwidgets "github.com/grep-michael/GoWidgets/ExWidgets"
)

func Test_coloredProgressBar(t *testing.T) {
	normal := widget.NewProgressBar()
	colored := exwidgets.NewColoredProgressBar(color.RGBA{R: 255, G: 0, B: 0, A: 255}, color.RGBA{R: 0, G: 0, B: 255, A: 128})
	window := test.NewWindow(container.NewVBox(normal, colored))
	window.Resize(fyne.NewSize(100, 100))
	test.AssertRendersToImage(t, "coloredProgress/colored_progress_0.png", window.Canvas())
	normal.SetValue(.5)
	colored.SetValue(.5)
	test.AssertRendersToImage(t, "coloredProgress/colored_progress_50.png", window.Canvas())
	colored.SetForegroundColor(color.RGBA{R: 0, G: 255, B: 0, A: 128})
	test.AssertRendersToImage(t, "coloredProgress/colored_progress_50_colorchange.png", window.Canvas())
}

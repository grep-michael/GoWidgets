package exwidgets

/*

	Lots of chuncks are copied from fynes progress bar widget
	https://github.com/fyne-io/fyne/blob/master/widget/progressbar.go

	Copyright (C) 2018 Fyne.io developers (see AUTHORS)
	All rights reserved.

*/

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ColoredProgressBarRender struct {
	BaseRenderer
	Background, Foreground             canvas.Rectangle
	foreground_color, background_color color.Color
	label                              canvas.Text
	percentage                         float64
	progressBar                        *ColoredProgressBar
}

func (p *ColoredProgressBarRender) MinSize() fyne.Size {
	text := "100%"
	if format := p.progressBar.TextFormatter; format != nil {
		text = format()
	}

	th := p.progressBar.Theme()
	padding := th.Size(theme.SizeNameInnerPadding) * 2
	size := fyne.MeasureText(text, p.label.TextSize, p.label.TextStyle)
	return size.AddWidthHeight(padding, padding)
}
func (p *ColoredProgressBarRender) calculateRatio() {
	if p.progressBar.Value < p.progressBar.Min {
		p.progressBar.Value = p.progressBar.Min
	}
	if p.progressBar.Value > p.progressBar.Max {
		p.progressBar.Value = p.progressBar.Max
	}

	delta := p.progressBar.Max - p.progressBar.Min
	p.percentage = (p.progressBar.Value - p.progressBar.Min) / delta
}
func (p *ColoredProgressBarRender) updateBar() {
	p.Layout(p.progressBar.Size())
	p.Background.Hidden = p.percentage == 1.0
	p.Foreground.Hidden = p.percentage == 0.0
	if text := p.progressBar.TextFormatter; text != nil {
		p.label.Text = text()
		return
	}
	p.label.Text = strconv.Itoa(int(p.percentage*100)) + "%"
}
func (p *ColoredProgressBarRender) Layout(size fyne.Size) {
	p.calculateRatio()

	p.Foreground.Resize(fyne.NewSize(size.Width*float32(p.percentage), size.Height))
	p.Background.Resize(size)
	p.label.Resize(size)
}
func (p *ColoredProgressBarRender) applyTheme() {
	th := p.progressBar.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	inputRadius := th.Size(theme.SizeNameInputRadius)

	p.Background.FillColor = p.background_color
	p.Background.CornerRadius = inputRadius
	p.Foreground.FillColor = p.foreground_color
	p.Foreground.CornerRadius = inputRadius

	p.label.Color = th.Color(theme.ColorNameForegroundOnPrimary, v)
	p.label.TextSize = th.Size(theme.SizeNameText)
}
func (p *ColoredProgressBarRender) Refresh() {
	p.applyTheme()
	p.updateBar()
	p.Background.Refresh()
	p.Foreground.Refresh()
	p.label.Refresh()
	canvas.Refresh(p.progressBar)
}

type ColoredProgressBar struct {
	widget.BaseWidget
	renderer                         *ColoredProgressBarRender
	ForegroundColor, BackgroundColor color.Color
	Min, Max, Value                  float64
	TextFormatter                    func() string `json:"-"`
	binder                           BasicBinder
}

func NewColoredProgressBar(foreground_color, background_color color.Color) *ColoredProgressBar {
	bar := &ColoredProgressBar{
		Min: 0, Max: 1,
		ForegroundColor: foreground_color, BackgroundColor: background_color,
	}
	bar.ExtendBaseWidget(bar)
	return bar
}
func NewColoredProgressBarWithData(foreground_color, background_color color.Color, data binding.Float) *ColoredProgressBar {
	p := NewColoredProgressBar(foreground_color, background_color)
	p.Bind(data)
	return p
}

func (p *ColoredProgressBar) SetBackgroundColor(color color.Color) {
	p.renderer.Background.FillColor = color
	p.renderer.Background.Refresh()
}
func (p *ColoredProgressBar) SetForegroundColor(color color.Color) {
	p.renderer.Foreground.FillColor = color
	p.renderer.Foreground.Refresh()
}
func (p *ColoredProgressBar) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatSource, ok := data.(binding.Float)
	if !ok {
		return
	}

	val, err := floatSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	p.SetValue(val)
}
func (p *ColoredProgressBar) Bind(data binding.Float) {
	p.binder.SetCallback(p.updateFromData)
	p.binder.Bind(data)
}
func (p *ColoredProgressBar) Unbind() {
	p.binder.Unbind()
}
func (p *ColoredProgressBar) SetValue(v float64) {
	p.Value = v
	p.Refresh()
}
func (p *ColoredProgressBar) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}
func (p *ColoredProgressBar) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	if p.Min == 0 && p.Max == 0 {
		p.Max = 1.0
	}

	p.renderer = &ColoredProgressBarRender{progressBar: p, foreground_color: p.ForegroundColor, background_color: p.BackgroundColor}
	p.renderer.label.Alignment = fyne.TextAlignCenter
	p.renderer.applyTheme()
	p.renderer.updateBar()

	p.renderer.SetObjects([]fyne.CanvasObject{&p.renderer.Background, &p.renderer.Foreground, &p.renderer.label})
	return p.renderer
}

package window

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"main/internal/domain"
	"main/internal/usecase"
)

type SettingsWindow struct {
	window   fyne.Window
	settings *usecase.SettingsManager
	onSave   func()
}

func NewSettingsWindow(app fyne.App, settings *usecase.SettingsManager, onSave func()) *SettingsWindow {
	window := app.NewWindow("設定")
	return &SettingsWindow{
		window:   window,
		settings: settings,
		onSave:   onSave,
	}
}

func (w *SettingsWindow) Show() {
	currentSettings := w.settings.GetSettings()

	thresholdEntry := widget.NewEntry()
	thresholdEntry.SetText(strconv.Itoa(currentSettings.ThresholdSeconds))

	saveBtn := widget.NewButton("儲存", func() {
		threshold, err := strconv.Atoi(thresholdEntry.Text)
		if err != nil {
			// TODO: 顯示錯誤訊息
			return
		}

		newSettings := domain.Settings{
			ThresholdSeconds: threshold,
		}

		if err := w.settings.UpdateSettings(newSettings); err != nil {
			// TODO: 顯示錯誤訊息
			return
		}

		if w.onSave != nil {
			w.onSave()
		}
		w.window.Close()
	})

	content := container.NewVBox(
		widget.NewLabel("閾值設定（秒）："),
		thresholdEntry,
		saveBtn,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(300, 200))
	w.window.Show()
} 
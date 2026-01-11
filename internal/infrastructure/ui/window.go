package ui

import (
	"io"
	"log/slog"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
)

type HealthService interface {
	CreateReport(filePaths []string) (*domain.ReportInfo, []error)
	Export(reportInfo *domain.ReportInfo, writer io.WriteCloser) error
}

// Window управляет окнами приложения
type Window struct {
	app           fyne.App
	Window        fyne.Window
	healthService HealthService

	selectFilesHint *widget.Label // Текст, подсказывающий, что нужно выбрать файлы
	fileList        *widget.List
	selectedFiles   []string
	fileListData    []string

	rawDataTable      *widget.Table
	completeDataTable *widget.Table
	exportBtn         *widget.Button
}

// NewWindowManager создает новый экземпляр Window
func NewWindow(app fyne.App, healthService HealthService) *Window {
	w := &Window{app: app, healthService: healthService, selectedFiles: make([]string, 0), fileListData: make([]string, 0)}
	w.Window = app.NewWindow("Health scale analyzer")
	w.Window.Resize(fyne.NewSize(900, 700))

	slog.Info("window created")

	w.createUI()

	return w
}

func (w *Window) createUI() {
	w.selectFilesHint = widget.NewLabel("Выберите файлы для анализа")
	w.selectFilesHint.Wrapping = fyne.TextWrapWord

	selectFilesBtn := widget.NewButton("Выбрать файлы", w.selectFiles)

	clearFilesBtn := widget.NewButton("Очистить выбор", w.clearFiles)
	clearFilesBtn.Importance = widget.LowImportance

	w.exportBtn = widget.NewButton("Анализ файлов и экспорт в Excel", func() { w.handleExport() })
	w.exportBtn.Disable()

	w.fileList = widget.NewList(
		func() int {
			return len(w.fileListData)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("template"),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			container := obj.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)
			deleteBtn := container.Objects[1].(*widget.Button)

			fileName := w.fileListData[id]
			label.SetText(filepath.Base(fileName))

			deleteBtn.OnTapped = func() {
				w.removeFile(id)
			}
		},
	)

	// Создаем таблицы для данных
	w.rawDataTable = widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// Будет заполняться после анализа
		},
	)
	w.rawDataTable.ShowHeaderRow = true
	w.rawDataTable.ShowHeaderColumn = true

	// Компоновка интерфейса
	header := container.NewVBox(
		container.NewHBox(
			selectFilesBtn,
			clearFilesBtn,
		),
		w.selectFilesHint,
	)

	filesSection := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		w.fileList,
	)

	// Панель внизу с кнопкой экспорта
	bottomSection := container.NewHBox(
		layout.NewSpacer(),
		w.exportBtn,
	)

	mainContainer := container.NewBorder(
		nil,
		bottomSection,
		nil,
		nil,
		filesSection,
	)

	w.Window.SetContent(mainContainer)
}

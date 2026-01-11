package ui

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// selectFiles открывает диалог выбора нескольких файлов
func (w *Window) selectFiles() {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.Window)
			return
		}
		if reader == nil {
			return
		}

		fileURI := reader.URI()
		filePath := fileURI.Path()

		for _, existingFile := range w.selectedFiles {
			if existingFile == filePath {
				dialog.ShowInformation("Файл уже добавлен",
					"Этот файл уже есть в списке", w.Window)
				return
			}
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		if ext != ".xlsx" && ext != ".xls" {
			dialog.ShowError(fmt.Errorf("поддерживаются только файлы Excel (.xlsx, .xls)"), w.Window)
			return
		}

		w.selectedFiles = append(w.selectedFiles, filePath)
		w.fileListData = append(w.fileListData, filePath)
		w.updateFileList()
	}, w.Window)

	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", ".xls"}))

	dialog.Show()
}

// updateFileList обновляет отображение списка файлов
func (w *Window) updateFileList() {
	if len(w.selectedFiles) > 0 {
		w.selectFilesHint.Hide()
	} else {
		w.selectFilesHint.SetText("Выберите файлы для анализа")
	}

	// Активация/деактивация кнопки анализа
	if len(w.selectedFiles) > 0 {
		w.exportBtn.Enable()
	} else {
		w.exportBtn.Disable()
	}

	w.fileList.Refresh()
}

// removeFile удаляет файл из списка по индексу
func (w *Window) removeFile(index int) {
	if index < 0 || index >= len(w.selectedFiles) {
		return
	}

	w.selectedFiles = append(w.selectedFiles[:index], w.selectedFiles[index+1:]...)
	w.fileListData = append(w.fileListData[:index], w.fileListData[index+1:]...)
	w.updateFileList()
}

// clearFiles очищает все выбранные файлы
func (w *Window) clearFiles() {
	if len(w.selectedFiles) == 0 {
		return
	}

	// Подтверждение удаления
	dialog.ShowConfirm("Очистить список",
		"Вы уверены, что хотите удалить все файлы из списка?",
		func(confirmed bool) {
			if confirmed {
				w.selectedFiles = []string{}
				w.fileListData = []string{}
				w.updateFileList()
				w.exportBtn.Disable()
			}
		}, w.Window)
}

// handleExport обрабатывает нажатие кнопки Export
func (w *Window) handleExport() {
	slog.Info("export button pressed")

	reportInfo, creatingErr := w.healthService.CreateReport(w.selectedFiles)
	for _, err := range creatingErr {
		dialog.ShowError(err, w.Window)
		return
	}

	dialogSave := NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.Window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		fileURI := writer.URI()
		if fileURI == nil {
			dialog.ShowError(errors.New("cannot get file URI"), w.Window)
			return
		}

		exportErr := w.healthService.Export(reportInfo, writer)
		if exportErr != nil {
			dialog.ShowError(exportErr, w.Window)
			return
		}
	}, w.Window)

	dialogSave.Show()
}

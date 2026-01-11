package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/application"
	bmiconfig "github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/infrastructure/bmi_config"
	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/infrastructure/export"
	fileparser "github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/infrastructure/file_parser"
	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/infrastructure/ui"
)

func main() {
	a := app.NewWithID("1")

	parser := fileparser.NewExcelParser()
	configLoader := bmiconfig.NewBMIConfigLoader()
	exporter := export.NewXLSX()

	healthService := application.NewHealthService(parser, configLoader, exporter)
	window := ui.NewWindow(a, healthService)
	window.Window.Show()
	a.Run()
}

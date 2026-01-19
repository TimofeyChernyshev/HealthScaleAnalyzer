package export

import (
	"io"

	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
	"github.com/xuri/excelize/v2"
)

type XlsxExporter struct{}

func NewXLSX() *XlsxExporter {
	return &XlsxExporter{}
}

func (e XlsxExporter) Export(report *domain.ReportInfo, writer io.WriteCloser) error {
	f := excelize.NewFile()

	// Лист с ИМТ каждого человека
	sheetZero := f.GetSheetName(0)
	err := f.SetSheetName(sheetZero, "ИМТ каждого")
	if err != nil {
		return err
	}
	bmiSheet := f.GetSheetName(0)

	headers := []string{
		"ФИО",
		"Индекс массы тела",
		"Оценка индекса",
		"Класс",
	}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(bmiSheet, cell, h)
	}

	for rowIndex, item := range report.AnalyzedPersons {
		row := rowIndex + 2

		values := []any{
			item.Name,
			item.BMI,
			item.BMICategory,
			item.Group,
		}

		for colIndex, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, row)
			f.SetCellValue(bmiSheet, cell, v)
		}
	}

	// Лист с общей статистикой по группам здоровья
	_, err = f.NewSheet("Группы здоровья")
	if err != nil {
		return err
	}
	healthSheet := f.GetSheetName(1)

	f.SetCellValue(healthSheet, "A1", "Номер группы")
	f.SetCellValue(healthSheet, "A2", "Количество")

	for i, groupName := range domain.AvailableHealthGroups {
		cell, _ := excelize.CoordinatesToCellName(i+2, 1)
		f.SetCellValue(healthSheet, cell, groupName)

		cell, _ = excelize.CoordinatesToCellName(i+2, 2)
		f.SetCellValue(healthSheet, cell, report.HealthGroups[groupName])
	}

	// Лист с общей статистикой по группам физкультуры
	_, err = f.NewSheet("Группы физкультуры")
	if err != nil {
		return err
	}
	physSheet := f.GetSheetName(2)

	f.SetCellValue(physSheet, "A1", "Название группы")
	f.SetCellValue(physSheet, "A2", "Количество")

	for i, groupName := range domain.AvailablePhysicalGroups {
		cell, _ := excelize.CoordinatesToCellName(i+2, 1)
		f.SetCellValue(physSheet, cell, groupName)

		cell, _ = excelize.CoordinatesToCellName(i+2, 2)
		f.SetCellValue(physSheet, cell, report.PhysicalGroups[groupName])
	}

	// Лист с общей статистикой по оценкам ИМТ
	_, err = f.NewSheet("Оценки ИМТ")
	if err != nil {
		return err
	}

	bmiCategorySheet := f.GetSheetName(3)
	for i, category := range domain.AvailableBMICategories {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(bmiCategorySheet, cell, category)

		cell, _ = excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(bmiCategorySheet, cell, report.BMICategories[category])
	}

	return f.Write(writer)
}

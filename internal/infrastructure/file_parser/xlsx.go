package fileparser

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
	"github.com/xuri/excelize/v2"
)

// ExcelParser реализация парсера Excel
type ExcelParser struct {
	dateFormats []string
}

// NewExcelParser создает новый парсер
func NewExcelParser() *ExcelParser {
	return &ExcelParser{
		dateFormats: []string{
			"02/01/2006",
			"02.01.2006",
			"02-01-2006",
			"2006-01-02",
			"01/02/2006",
			"01.02.2006",
			"01-02-2006",
		},
	}
}

// ParseFiles парсит несколько файлов
func (p *ExcelParser) ParseFiles(filePaths []string) ([]*domain.Person, []error) {
	var allPeople []*domain.Person
	errors := make([]error, 0)

	for _, filePath := range filePaths {
		people, err := p.parse(filePath)
		if err != nil {
			slog.Warn("error parsing file '%s': %w\n", filepath.Base(filePath), err)
			errors = append(errors, err)
			continue
		}
		allPeople = append(allPeople, people...)
	}

	return allPeople, errors
}

// parse парсит один Excel файл
func (p *ExcelParser) parse(filePath string) ([]*domain.Person, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл '%s': %w", filepath.Base(filePath), err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("файл '%s' не содержит листов", filepath.Base(filePath))
	}

	var allPeople []*domain.Person

	for _, sheet := range sheets {
		people, err := p.parseSheet(f, sheet)
		if err != nil {
			return nil, fmt.Errorf("не удалось обработать лист: %w", err)
		}
		allPeople = append(allPeople, people...)
	}

	if len(allPeople) == 0 {
		return nil, fmt.Errorf("не удалось извлечь данные из файла '%s'", filepath.Base(filePath))
	}

	return allPeople, nil
}

// parseSheet парсит один лист Excel
func (p *ExcelParser) parseSheet(f *excelize.File, sheetName string) ([]*domain.Person, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения листа '%s': %w", sheetName, err)
	}

	headerRow := rows[0]
	colIndexes := p.mapColumns(headerRow)

	var people []*domain.Person
	for i := 1; i < len(rows); i++ {
		person, err := p.parseRow(rows[i], colIndexes, i+1)
		if err != nil {
			continue
		}
		people = append(people, person)
	}

	return people, nil
}

// mapColumns определяет индексы колонок по заголовкам
func (p *ExcelParser) mapColumns(headers []string) map[string]int {
	indexes := make(map[string]int)

	for i, header := range headers {
		header = normalizeHeader(header)
		fmt.Println(header)

		switch {
		case strings.Contains(header, "фио"), strings.Contains(header, "fullName"):
			indexes["fullName"] = i
		case strings.Contains(header, "датарождения"), strings.Contains(header, "др"),
			strings.Contains(header, "birthdate"), strings.Contains(header, "birthday"):
			indexes["birthDate"] = i
		case strings.Contains(header, "вес"), strings.Contains(header, "weight"),
			strings.Contains(header, "масса"):
			indexes["weight"] = i
		case header == "рост", header == "height":
			indexes["height"] = i
		case strings.Contains(header, "группаздоровья"), strings.Contains(header, "healthgroup"):
			indexes["healthGroup"] = i
		case strings.Contains(header, "группафизкультуры"), strings.Contains(header, "физкультурнаягруппа"),
			strings.Contains(header, "sportgroup"), strings.Contains(header, "physicalgroup"):
			indexes["physicalGroup"] = i
		}
	}

	return indexes
}

// normalizeHeader нормализует заголовок
func normalizeHeader(header string) string {
	header = strings.ToLower(strings.TrimSpace(header))

	header = strings.NewReplacer(
		"_", "",
		"-", "",
		"(", "",
		")", "",
		".", "",
		",", "",
		" ", "",
	).Replace(header)

	return header
}

// parseRow парсит одну строку Excel
func (p *ExcelParser) parseRow(row []string, colIndexes map[string]int, rowNum int) (*domain.Person, error) {
	person := &domain.Person{}

	// ФИО
	if i, ok := colIndexes["fullName"]; ok && i < len(row) {
		person.Name = strings.TrimSpace(row[i])
	}

	// Дата рождения
	if i, ok := colIndexes["birthDate"]; ok && i < len(row) {
		if dateStr := strings.TrimSpace(row[i]); dateStr != "" {
			birthDate, err := p.parseDate(dateStr)
			if err == nil {
				person.BirthDate = birthDate
			}
		}
	}

	// Вес
	if i, ok := colIndexes["weight"]; ok && i < len(row) {
		if weightStr := strings.TrimSpace(row[i]); weightStr != "" {
			if weight, err := strconv.ParseFloat(weightStr, 64); err == nil && weight > 0 {
				person.Weight = weight
			}
		}
	}

	// Рост
	if i, ok := colIndexes["height"]; ok && i < len(row) {
		if heightStr := strings.TrimSpace(row[i]); heightStr != "" {
			if height, err := strconv.ParseFloat(heightStr, 64); err == nil && height > 0 {
				person.Height = height
			}
		}
	}

	// Группа здоровья
	if i, ok := colIndexes["healthGroup"]; ok && i < len(row) {
		person.HealthGroup = strings.TrimSpace(row[i])
	}

	// Группа физкультуры
	if i, ok := colIndexes["physicalGroup"]; ok && i < len(row) {
		person.PhysicalGroup = strings.TrimSpace(row[i])
	}

	if person.Name == "" {
		return nil, fmt.Errorf("строка %d: отсутствуют ФИО", rowNum)
	}

	if person.Weight <= 0 || person.Height <= 0 {
		return nil, fmt.Errorf("строка %d: некорректные вес или рост", rowNum)
	}

	return person, nil
}

// parseDate парсит дату из строки
func (p *ExcelParser) parseDate(dateStr string) (time.Time, error) {
	// Пробуем все форматы
	for _, format := range p.dateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Пробуем как число Excel (серийный номер даты)
	if serial, err := strconv.ParseFloat(dateStr, 64); err == nil {
		// Конвертируем из Excel серийного формата
		excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		return excelEpoch.Add(time.Duration(serial * 24 * float64(time.Hour))), nil
	}

	return time.Time{}, fmt.Errorf("неизвестный формат даты: %s", dateStr)
}

// ValidateFile проверяет файл перед парсингом
func (p *ExcelParser) ValidateFile(filePath string) error {
	// Проверяем расширение
	ext := filepath.Ext(filePath)
	if ext != ".xlsx" && ext != ".xls" {
		return fmt.Errorf("неподдерживаемый формат файла: %s", ext)
	}

	// Пробуем открыть файл
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("файл поврежден или не является Excel: %v", err)
	}
	defer f.Close()

	// Проверяем, есть ли листы
	if sheets := f.GetSheetList(); len(sheets) == 0 {
		return fmt.Errorf("файл не содержит листов")
	}

	return nil
}

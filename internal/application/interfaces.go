package application

import (
	"io"

	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
)

// Parser интерфейс для парсинга Excel файлов
type Parser interface {
	ParseFiles(filePaths []string) ([]*domain.Person, []error)
}

// ConfigLoader интерфейс для загрузки конфигов BMI для мужского и женского полов
type ConfigLoader interface {
	LoadBMIConfig() (*domain.BMIConfig, *domain.BMIConfig, error)
}

type Exporter interface {
	Export(*domain.ReportInfo, io.WriteCloser) error
}

package application

import (
	"io"
	"strings"

	"github.com/TimofeyChernyshev/HealthScaleAnalyzer/internal/domain"
)

// HealthService представляет систему создания отчетов посещаемости
type HealthService struct {
	parser       Parser
	configLoader ConfigLoader
	exporter     Exporter
}

// NewHealthService создает новый экземпляр HealthService
func NewHealthService(p Parser, l ConfigLoader, e Exporter) *HealthService {
	return &HealthService{parser: p, configLoader: l, exporter: e}
}

func (s *HealthService) CreateReport(filePaths []string) (*domain.ReportInfo, []error) {
	persons, errors := s.parser.ParseFiles(filePaths)
	if len(errors) != 0 {
		return nil, errors
	}

	maleBMI, femaleBMI, err := s.configLoader.LoadBMIConfig()
	if err != nil {
		return nil, []error{err}
	}

	reportInfo := domain.NewReportInfo()

	for _, p := range persons {
		if strings.HasSuffix(p.Name, "ич") {
			reportInfo.AnalyzePerson(p, maleBMI)
		} else {
			reportInfo.AnalyzePerson(p, femaleBMI)
		}
	}

	return reportInfo, nil
}

func (s *HealthService) Export(reportInfo *domain.ReportInfo, writer io.WriteCloser) error {
	err := s.exporter.Export(reportInfo, writer)
	if err != nil {
		return err
	}

	return nil
}

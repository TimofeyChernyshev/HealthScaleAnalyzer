package domain

import (
	"strings"
	"time"
)

type ReportInfo struct {
	AnalyzedPersons []AnalyzedPerson
	HealthGroups    map[HealthGroup]int
	PhysicalGroups  map[PhysicalGroup]int
	BMICategories   map[BMICategory]int
}

func NewReportInfo() *ReportInfo {
	return &ReportInfo{
		HealthGroups:   map[HealthGroup]int{First: 0, Second: 0, Third: 0, Fourth: 0, Fifth: 0},
		PhysicalGroups: map[PhysicalGroup]int{Default: 0, Prepare: 0, Special: 0},
		BMICategories:  map[BMICategory]int{SeriousUnder: 0, LightUnder: 0, Normal: 0, OverWeight: 0, Obese: 0},
	}
}

type AnalyzedPerson struct {
	Name        string
	BMI         float64
	BMICategory BMICategory
}

func (r *ReportInfo) AnalyzePerson(p *Person, maleBMI, femaleBMI *BMIConfig) {
	age := time.Now().Year() - p.BirthDate.Year()
	if time.Now().YearDay() < p.BirthDate.YearDay() {
		age--
	}

	var config *BMIConfig
	// На русском языке пол можно определить по отчеству
	if strings.HasSuffix(p.Name, "ич") {
		config = maleBMI
	} else {
		config = femaleBMI
	}

	// bmi = вес/рост(м)^2, p.Height - рост в см
	bmi := p.Weight / (p.Height * p.Height * 1_0000.0)
	bmiCategory := categorizeBMI(bmi, age, config)

	analyzed := AnalyzedPerson{
		Name:        p.Name,
		BMI:         bmi,
		BMICategory: bmiCategory,
	}
	r.AnalyzedPersons = append(r.AnalyzedPersons, analyzed)

	hg := strings.TrimSpace(strings.ToLower(p.HealthGroup))
	switch hg {
	case "1", "i":
		r.HealthGroups[First]++
	case "2", "ii":
		r.HealthGroups[Second]++
	case "3", "iii":
		r.HealthGroups[Third]++
	case "4", "iv":
		r.HealthGroups[Fourth]++
	case "5", "v":
		r.HealthGroups[Fifth]++
	}

	pg := strings.TrimSpace(strings.ToLower(p.PhysicalGroup))
	switch {
	case strings.HasPrefix(pg, "о"):
		r.PhysicalGroups[Default]++
	case strings.HasPrefix(pg, "п"):
		r.PhysicalGroups[Prepare]++
	case strings.HasPrefix(pg, "с"):
		r.PhysicalGroups[Special]++
	}

	r.BMICategories[bmiCategory]++
}

func categorizeBMI(bmi float64, age int, bmiConfig *BMIConfig) BMICategory {
	for i, cfg := range bmiConfig.Ranges {
		if cfg.MinAge <= age && (cfg.MaxAge >= age || i == len(bmiConfig.Ranges)-1) {
			switch {
			case bmi >= cfg.Obese:
				return Obese
			case bmi >= cfg.OverWeight:
				return OverWeight
			case bmi >= cfg.Normal:
				return Normal
			case bmi >= cfg.LightUnder:
				return LightUnder
			default:
				return SeriousUnder
			}
		}
	}

	return BMICategory("")
}

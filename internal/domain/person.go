package domain

import "time"

type Person struct {
	Name          string
	BirthDate     time.Time
	Weight        float64
	Height        float64
	HealthGroup   string
	PhysicalGroup string
	Group         string
}

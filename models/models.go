package models

import (
	"time"

	"github.com/dadosjusbr/storage"
)

// State - Struct cotains information of a state ans its agencies
type State struct {
	Name      string
	ShortName string
	FlagURL   string
	Agency    []AgencyBasic
}

// AgencyBasic - Basic information of a agency (name e category)
type AgencyBasic struct {
	Name           string
	FullName       string
	AgencyCategory string
}

// Employee - Represents an employee and his/her salary info
type Employee struct {
	Name   string
	Wage   float64
	Perks  float64
	Others float64
	Total  float64
	Type   string
	Active bool
}

// AgencySummary - Summary of an agency
type AgencySummary struct {
	FullName          string
	TotalEmployees    int
	TotalWage         float64
	TotalPerks        float64
	MaxWage           float64
	CrawlingTime      time.Time
	AgencyName        string
	TotalMembers      int
	TotalServants     int
	TotalInactives    int
	MaxPerk           float64
	TotalRemuneration float64
}

// AgencyTotalsYear - Represents the totals of an year
type AgencyTotalsYear struct {
	Year        int
	MonthTotals []MonthTotals
}

// MonthTotals - Detailed info of a month (wage, perks, other)
type MonthTotals struct {
	Month  int
	Wage   float64
	Perks  float64
	Others float64
}

// DataForChartAtAgencyScreen - contains all necessary data to load chart
type DataForChartAtAgencyScreen struct {
	Members     map[int]int
	Servers     map[int]int
	Inactives   map[int]int
	MaxSalary   float64
	PackageURL  string
	PackageHash string
}

// ProcInfoResult - contains information of the result of the process if something went wrong during parsing or crawling process
type ProcInfoResult struct {
	ProcInfo          *storage.ProcInfo
	CrawlingTimestamp time.Time
}

// MonthAndYear - Most recent year and month with data available
type MonthAndYear struct {
	Month int
	Year  int
}

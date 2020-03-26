package models

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
	TotalEmployees int
	TotalWage      float64
	TotalPerks     float64
	MaxWage        float64
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

type DataForChart struct {
	Members   map[int]int
	Servers   map[int]int
	Inactives map[int]int
	MaxSalary float64
}

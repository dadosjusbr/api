package models

type State struct {
	Name      string
	ShortName string
	FlagURL   string
	Agency    []AgencyBasic
}

type AgencyBasic struct {
	Name           string
	AgencyCategory string
}

type Employee struct {
	Name   string
	Wage   float64
	Perks  float64
	Others float64
	Total  float64
}

type AgencySummary struct {
	TotalEmployees int
	TotalWage      float64
	TotalPerks     float64
	MaxWage        float64
}

type AgencyTotalsYear struct {
	Year        int
	MonthTotals []MonthTotals
}

type MonthTotals struct {
	Month  int
	Wage   float64
	Perks  float64
	Others float64
}

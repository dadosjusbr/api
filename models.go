package main

type state struct {
	Name      string
	ShortName string
	FlagURL   string
	Agency    []agencyBasic
}

type agencyBasic struct {
	Name           string
	AgencyCategory string
}

type employee struct {
	Name   string
	Wage   float64
	Perks  float64
	Others float64
	Total  float64
}

type agencySummary struct {
	TotalEmployees int
	TotalWage      float64
	TotalPerks     float64
	MaxWage        float64
}

type agencyTotalsYear struct {
	Year        int
	MonthTotals []monthTotals
}

type monthTotals struct {
	Month  int
	Wage   float64
	Perks  float64
	Others float64
}

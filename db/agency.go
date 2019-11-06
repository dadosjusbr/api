package storage

import "time"

// Struct containing the main descriptions of each Agency.
type Agency struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	ShortName string      // 'trt13'
	Name      string      // 'Tribunal Regional do Trabalho 13° Região'
	Type      string      // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity    string      // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF        string      // Short code for federative unity.
}

// Struct containing all necessary data to build all UI screens
type AgencyMonthlyInfo struct {
	AgencyID  string
	Storage   []StorageFile // Link to original files.
	Month     int
	Year      int
	Summary   Summary
	Employee  []Employee
	Timestamp time.Time
}

type StorageFile struct {
	Link string
	Hash string
}

// Data Summary with statistics.
type DataSummary struct {
	Max   float64
	Min   float64
	Mean  float64
	Total float64
}

type Summary struct {
	Count  int         // Number of employees
	Wage   DataSummary //  Statistics (Max, Min, Median, Total)
	Perks  DataSummary //  Statistics (Max, Min, Median, Total)
	Others DataSummary //  Statistics (Max, Min, Median, Total)
}

type Employee struct {
	Reg            string // Register number
	Name           string
	Role           string
	Type           string  // servidor, membro, pensionista or indefinido
	Workplace      string  // 'Lotacao' Like '10° Zona eleitoral'
	Active         bool    // 'Active' Or 'Inactive'
	GrossIncome    float64 // Income received without discounts applied.
	TotalDiscounts float64 // Discounts to be applied in Gross Income
	NetIncome      float64 // Final income received after discounts applied
	Income         Income  //
	Discounts      Discount
}

type Income struct {
	Wage  float64
	Perks Perks
	Other Funds // other funds that make up the total income of the employee. further details explained below
}

type Perks struct {
	Total         float64
	Food          float64 // Food Aid
	Tranportation float64
	PreSchool     float64 // Assistance provided before the child enters school.
	Health        float64
	BirthAid      float64
	HousingAid    float64 // 'Auxílio Moradia'
	Subsistence   float64 // 'Ajuda de Custo'
	Others        map[string]float64
}

type Funds struct {
	Total            float64
	PersonalBenefits float64            // Permanent Allowance, VPI, Benefits adquired thought judicial demand and others personal.
	EventualBenefits float64            // Holidays, Others Temporary Wage,  Christmas bonus and some others eventual.
	PositionOfTrust  float64            // Income given for the importance of the position held.
	Daily            float64            // Employee reimbursement for eventual expenses when working in a different location than usual.
	Gratification    float64            //
	OriginPosition   float64            // Wage received from other Agency, transfered employee.
	Others           map[string]float64 // Any other kind of income that does not have a pattern among the Agencys.
}

type Discount struct {
	PrevContribution float64 // 'Contribuição Previdenciária'
	CeilRetention    float64 // 'Retenção de teto'
	IncomeTax        float64 // 'Imposto de renda'
	Sundry           float64 // 'Diversos'
}

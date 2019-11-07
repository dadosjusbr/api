package storage

import "time"

// A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	ShortName string      `json:"short_name" bson:"short_name,omitempty"` // 'trt13'
	Name      string      `json:"name" bson:"name,omitempty"`             // 'Tribunal Regional do Trabalho 13° Região'
	Type      string      `json:"type" bson:"type,omitempty"`             // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity    string      `json:"entity" bson:"entity,omitempty"`         // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF        string      `json:"uf" bson:"uf,omitempty"`                 // Short code for federative unity.
}

//A Struct containing a snapshot of a agency in a month.
type AgencyMonthlyInfo struct {
	AgencyID string     `json:"id" bson:"_id,omitempty"`
	Storage  []Metadata `json:"storage" bson:"storage,omitempty"`
	Month    int        `json:"month" bson:"month,omitempty"`
	Year     int        `json:"year" bson:"year,omitempty"`
	Summary  Summary    `json:"sumarry" bson:"summary,omitempty"`
	Employee []Employee `json:"employee" bson:"employee,omitempty"`
	Metadata Metadata   `json:"metadata" bson:"metadata,omitempty"`
}

//A Struct containing metadatas about crawler commit
type Metadata struct {
	Timestamp      time.Time `json:"metadata" bson:"metadata,omitempty"`     // Time the crawler sent it
	CrawlerId      string    `json:"crawl_id" bson:"crawl_id,omitempty"`     // The directory of the collector's crawler
	CrawlerVersion string    `json:"crawl_vers" bson:"crawl_vers,omitempty"` // Last Commit of the repository
}

// A Struct containing URL to download a file and a hash to track if in the future will be changes in the file.
type FileBackup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
}

// A Struct containing summarized  information about a agency/month stats
type Summary struct {
	Count  int         `json:"count" bson:"count,omitempty"`   // Number of employees
	Wage   DataSummary `json:"wage" bson:"wage,omitempty"`     //  Statistics (Max, Min, Median, Total)
	Perks  DataSummary `json:"perks" bson:"perks,omitempty"`   //  Statistics (Max, Min, Median, Total)
	Others DataSummary `json:"others" bson:"others,omitempty"` //  Statistics (Max, Min, Median, Total)
}

// Data Summary with statistics.
type DataSummary struct {
	Max   float64 `json:"max" bson:"max,omitempty"`
	Min   float64 `json:"min" bson:"min,omitempty"`
	Mean  float64 `json:"mean" bson:"mean,omitempty"`
	Total float64 `json:"total" bson:"total,omitempty"`
}

// A Struct that reflets a employee snapshot, containing all relative data about a employee
type Employee struct {
	Reg       string        `json:"reg" bson:"reg,omitempty"` // Register number
	Name      string        `json:"name" bson:"name,omitempty"`
	Role      string        `json:"role" bson:"role,omitempty"`
	Type      string        `json:"type" bson:"type,omitempty"`           // servidor, membro, pensionista or indefinido
	Workplace string        `json:"workplace" bson:"workplace,omitempty"` // 'Lotacao' Like '10° Zona eleitoral'
	Active    bool          `json:"active" bson:"active,omitempty"`       // 'Active' Or 'Inactive'
	Income    IncomeDetails `json:"income" bson:"income,omitempty"`
	Discounts Discount      `json:"discounts" bson:"discounts,omitempty"`
}

// Struct that details an employee's income.
type IncomeDetails struct {
	Total float64  `json:"total" bson:"total,omitempty"`
	Wage  *float64 `json:"wage" bson:"wage,omitempty"`
	Perks Perks    `json:"perks" bson:"perks,omitempty"`
	Other Funds    `json:"other" bson:"other,omitempty"` // other funds that make up the total income of the employee. further details explained below
}

// About used pointers.
// All pointers are important to know if in the field has information and this information is 0 or if we do not have information about that field.
// For a example, a Funds Daily field with null will represent that we do not have that information, but a Dialy field with 0, represents that we have that information and the employee received 0 Reais in Daily Funds
// On the other hand, if we dont put pointer in those fields, Funds daily will be setted 0 as a float64 primitive number, and we will not be able to
// diferenciate if we have the 0 information or if we dont know about it.
// The point here is just to guarantee that what appears in the system are real collected data.
// As disavantage we add some complexity to code knowing that the final value will not be changed anyway.

// Struct that details perks that complements an employee's wage.
type Perks struct {
	Total         *float64           `json:"total" bson:"total,omitempty"`
	Food          *float64           `json:"food" bson:"food,omitempty"` // Food Aid
	Tranportation *float64           `json:"transportation" bson:"transportation,omitempty"`
	PreSchool     *float64           `json:"pre_school" bson:"pre_school,omitempty"` // Assistance provided before the child enters school.
	Health        *float64           `json:"health" bson:"health,omitempty"`
	BirthAid      *float64           `json:"birth_aid" bson:"birth_aid,omitempty"`     // 'Auxílio Natalidade'
	HousingAid    *float64           `json:"housing_aid" bson:"housing_aid,omitempty"` // 'Auxílio Moradia'
	Subsistence   *float64           `json:"subsistence" bson:"subsistence,omitempty"` // 'Ajuda de Custo'
	Others        map[string]float64 `json:"others" bson:"others,omitempty"`           // Any other kind of perk that does not have a pattern among the Agencys.
}

// A Struct that details that make up the employee income.
type Funds struct {
	Total            float64            `json:"total" bson:"total,omitempty"`
	PersonalBenefits *float64           `json:"person_benefits" bson:"person_benefits,omitempty"`     // Permanent Allowance, VPI, Benefits adquired thought judicial demand and others personal.
	EventualBenefits *float64           `json:"eventual_benefits" bson:"eventual_benefits,omitempty"` // Holidays, Others Temporary Wage,  Christmas bonus and some others eventual.
	PositionOfTrust  *float64           `json:"trust_position" bson:"trust_position,omitempty"`       // Income given for the importance of the position held.
	Daily            *float64           `json:"daily" bson:"daily,omitempty"`                         // Employee reimbursement for eventual expenses when working in a different location than usual.
	Gratification    *float64           `json:"gratific" bson:"gratific,omitempty"`                   //
	OriginPosition   *float64           `json:"origin_pos" bson:"origin_pos,omitempty"`               // Wage received from other Agency, transfered employee.
	Others           map[string]float64 `json:"others" bson:"others,omitempty"`                       // Any other kind of income that does not have a pattern among the Agencys.
}

// A Struct that details all discounts that must be applied to the employee's income.
type Discount struct {
	Total            float64  `json:"total" bson:"total,omitempty"`
	PrevContribution *float64 `json:"prev_contribution" bson:"prev_contribution,omitempty"` // 'Contribuição Previdenciária'
	CeilRetention    *float64 `json:"ceil_retention" bson:"ceil_retention,omitempty"`       // 'Retenção de teto'
	IncomeTax        *float64 `json:"income_tax" bson:"income_tax,omitempty"`               // 'Imposto de renda'
	Sundry           float64  `json:"sundry" bson:"sundry,omitempty"`                       // 'Diversos'
}

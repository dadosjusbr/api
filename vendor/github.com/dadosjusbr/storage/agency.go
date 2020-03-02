// About used pointers.
// All pointers are important to know if in the field has information and this information is 0 or if we do not have information about that field.
// This is justified because of the use of omitempty. If a collected float64 is 0, it will not appear in the json fields, cause that's it's zero value.
// Any application consuming this data might not know if the field is really 0 or data is unavailable.
// For a example, a Funds Daily field with null will represent that we do not have that information, but a Dialy field with 0, represents that we have that information and the employee received 0 Reais in Daily Funds
// On the other hand, if we dont put pointer in those fields, Funds daily will be setted 0 as a float64 primitive number, and we will not be able to
// diferenciate if we have the 0 information or if we dont know about it.
// The point here is just to guarantee that what appears in the system are real collected data.
// As disavantage we add some complexity to code knowing that the final value will not be changed anyway.
// Use Case:
// Pointers                                 No Pointers
// daily: nil                              daily: 0
// perks: nil							   perks: 0
// total: 0								   total: 0

package storage

import "time"

// Crawler keeps information about the crawler.
type Crawler struct {
	CrawlerID      string `json:"id" bson:"id,omitempty"`           // Convention: crawler the directory
	CrawlerVersion string `json:"version" bson:"version,omitempty"` // Convention: crawler commit id
}

// CrawlingResult stores the result of a crawler-parser ("coletor") run.
type CrawlingResult struct {
	AgencyID  string     `json:"aid"`
	Month     int        `json:"month"`
	Year      int        `json:"year"`
	Crawler   Crawler    `json:"crawler"`
	Files     []string   `json:"files"`
	Employees []Employee `json:"employees"`
	Timestamp time.Time  `json:"timestamp"`
}

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID      string `json:"aid" bson:"aid,omitempty"`       // 'trt13'
	Name    string `json:"name" bson:"name,omitempty"`     // 'Tribunal Regional do Trabalho 13° Região'
	Type    string `json:"type" bson:"type,omitempty"`     // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity  string `json:"entity" bson:"entity,omitempty"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF      string `json:"uf" bson:"uf,omitempty"`         // Short code for federative unity.
	FlagURL string `json:"url" bson:"url,omitempty"`       //Link for state url
}

// AgencyMonthlyInfo A Struct containing a snapshot of a agency in a month.
type AgencyMonthlyInfo struct {
	AgencyID          string     `json:"aid,omitempty" bson:"aid,omitempty"`
	Month             int        `json:"month,omitempty" bson:"month,omitempty"`
	Year              int        `json:"year,omitempty" bson:"year,omitempty"`
	Backups           []Backup   `json:"backups,omitempty" bson:"backups,omitempty"`
	Summary           Summaries  `json:"summary,omitempty" bson:"summary,omitempty"`
	Employee          []Employee `json:"employee,omitempty" bson:"employee,omitempty"`
	Crawler           Crawler    `json:"crawler,omitempty" bson:"crawler,omitempty"`
	CrawlingTimestamp time.Time  `json:"ts,omitempty" bson:"ts,omitempty"` // Crawling moment (always UTC)
}

// Backup contains the URL to download a file and a hash to track if in the future will be changes in the file.
type Backup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
}

// Summaries contains all summary detailed information
type Summaries struct {
	General         Summary
	MemberActive    Summary
	MemberInactive  Summary
	ServantActive   Summary
	ServantInactive Summary
}

// Summary A Struct containing summarized  information about a agency/month stats
type Summary struct {
	Count  int         `json:"count" bson:"count,omitempty"`   // Number of employees
	Wage   DataSummary `json:"wage" bson:"wage,omitempty"`     //  Statistics (Max, Min, Median, Total)
	Perks  DataSummary `json:"perks" bson:"perks,omitempty"`   //  Statistics (Max, Min, Median, Total)
	Others DataSummary `json:"others" bson:"others,omitempty"` //  Statistics (Max, Min, Median, Total)
}

// DataSummary A Struct containing data summary with statistics.
type DataSummary struct {
	Max     float64 `json:"max" bson:"max,omitempty"`
	Min     float64 `json:"min" bson:"min,omitempty"`
	Average float64 `json:"avg" bson:"avg,omitempty"`
	Total   float64 `json:"total" bson:"total,omitempty"`
}

// Employee a Struct that reflets a employee snapshot, containing all relative data about a employee
type Employee struct {
	Reg       string         `json:"reg" bson:"reg,omitempty"` // Register number
	Name      string         `json:"name" bson:"name,omitempty"`
	Role      string         `json:"role" bson:"role,omitempty"`
	Type      string         `json:"type" bson:"type,omitempty"`           // servidor, membro, pensionista or indefinido
	Workplace string         `json:"workplace" bson:"workplace,omitempty"` // 'Lotacao' Like '10° Zona eleitoral'
	Active    bool           `json:"active" bson:"active,omitempty"`       // 'Active' Or 'Inactive'
	Income    *IncomeDetails `json:"income" bson:"income,omitempty"`
	Discounts *Discount      `json:"discounts" bson:"discounts,omitempty"`
}

// IncomeDetails a Struct that details an employee's income.
type IncomeDetails struct {
	Total float64  `json:"total" bson:"total,omitempty"`
	Wage  *float64 `json:"wage" bson:"wage,omitempty"`
	Perks *Perks   `json:"perks" bson:"perks,omitempty"`
	Other *Funds   `json:"other" bson:"other,omitempty"` // other funds that make up the total income of the employee. further details explained below
}

// Perks a Struct that details perks that complements an employee's wage.
type Perks struct {
	Total          float64            `json:"total" bson:"total,omitempty"`
	Food           *float64           `json:"food" bson:"food,omitempty"` // Food Aid
	Transportation *float64           `json:"transportation" bson:"transportation,omitempty"`
	PreSchool      *float64           `json:"pre_school" bson:"pre_school,omitempty"` // Assistance provided before the child enters school.
	Health         *float64           `json:"health" bson:"health,omitempty"`
	BirthAid       *float64           `json:"birth_aid" bson:"birth_aid,omitempty"`     // 'Auxílio Natalidade'
	HousingAid     *float64           `json:"housing_aid" bson:"housing_aid,omitempty"` // 'Auxílio Moradia'
	Subsistence    *float64           `json:"subsistence" bson:"subsistence,omitempty"` // 'Ajuda de Custo'
	Others         map[string]float64 `json:"others" bson:"others,omitempty"`           // Any other kind of perk that does not have a pattern among the Agencys.
}

// Funds a Struct that details that make up the employee income.
type Funds struct {
	Total            float64            `json:"total" bson:"total,omitempty"`
	PersonalBenefits *float64           `json:"person_benefits" bson:"person_benefits,omitempty"`     // Permanent Allowance, VPI, Benefits adquired thought judicial demand and others personal.
	EventualBenefits *float64           `json:"eventual_benefits" bson:"eventual_benefits,omitempty"` // Holidays, Others Temporary Wage,  Christmas bonus and some others eventual.
	PositionOfTrust  *float64           `json:"trust_position" bson:"trust_position,omitempty"`       // Income given for the importance of the position held.
	Daily            *float64           `json:"daily" bson:"daily,omitempty"`                         // Employee reimbursement for eventual expenses when working in a different location than usual.
	Gratification    *float64           `json:"gratification" bson:"gratification,omitempty"`         //
	OriginPosition   *float64           `json:"origin_pos" bson:"origin_pos,omitempty"`               // Wage received from other Agency, transfered employee.
	Others           map[string]float64 `json:"others" bson:"others,omitempty"`                       // Any other kind of income that does not have a pattern among the Agencys.
}

// Discount a Struct that details all discounts that must be applied to the employee's income.
type Discount struct {
	Total            float64            `json:"total" bson:"total,omitempty"`
	PrevContribution *float64           `json:"prev_contribution" bson:"prev_contribution,omitempty"` // 'Contribuição Previdenciária'
	CeilRetention    *float64           `json:"ceil_retention" bson:"ceil_retention,omitempty"`       // 'Retenção de teto'
	IncomeTax        *float64           `json:"income_tax" bson:"income_tax,omitempty"`               // 'Imposto de renda'
	Others           map[string]float64 `json:"other" bson:"other,omitempty"`                         // Any other kind of discount that does not have a pattern among the Agencys.
}

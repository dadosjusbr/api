package structs

// Struct containing the main descriptions of each Agency.
type Agency struct {
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	Short_name string
	Name       string
	Type       string
	Entity     string
	UF         string
	State      string
}

// Struct containing all necessary data to build all UI screens
type DataMonthAgency struct {
	Agency           string
	Storage          string
	Month            int
	Year             int
	MaxWage          float32
	MinWage          float32
	TotalRestitution float32 // Total em indenizações
	Members          map[string]DataSummary
	Servers          map[string]DataSummary
	Pensioners       ShortData
	Employee         []Employee
}

type DataSummary struct {
	Resume   ShortData
	Active   ShortData
	Inactive ShortData
}

type ShortData struct {
	TotalQte         int
	MaxWage          float32
	MinWage          float32
	TotalWage        float32
	TotalRestitution float32
}

type Employee struct {
	Reg            string
	Name           string
	Role           string
	Type           string
	Lotacao        string
	Ativo          bool
	GrossIncome    float32
	TotalDiscounts float32
	NetIncome      float32
	Income         Income
	Discounts      Discount
}

type Income struct {
	Wage        float32
	Restitution float32
	OthersFunds Funds
}

type Funds struct {
	PersonalBenefits   float32
	EventualBenefits   float32
	PositionOfTrust    float32
	ChristmasBonus     float32
	vacation           float32
	PermanentAllowance float32
	OthersTemporary    float32
	Daily              float32
	OriginPosition     float32
}

type Discount struct {
	PrevContribution float32
	CeilRetention    float32
	IncomeTax        float32
	sundry           float32
}

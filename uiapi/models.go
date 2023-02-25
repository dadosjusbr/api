package uiapi

import (
	"time"

	"github.com/dadosjusbr/proto/coleta"
	"github.com/dadosjusbr/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// dataForChartAtAgencyScreen - contains all necessary data to load chart
type dataForChartAtAgencyScreen struct {
	Members     map[int]int
	Servers     map[int]int
	MaxSalary   float64
	PackageURL  string
	PackageHash string
	PackageSize int64
}

type agencySalary struct {
	MaxSalary float64     `json:"max_salario"`
	Histogram map[int]int `json:"histograma"`
	Package   *backup     `json:"package"`
}

type backup struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

// generalTotals - contains the summary from all DadosJusBr data
type generalTotals struct {
	AgencyAmount             int
	MonthlyTotalsAmount      int
	StartDate                time.Time
	EndDate                  time.Time
	RemunerationRecordsCount int
	GeneralRemunerationValue float64
}

// State - Struct cotains information of a state ans its agencies
type state struct {
	Name      string
	ShortName string
	FlagURL   string
	Agency    []agencyBasic
}

// AgencyBasic - Basic information of a agency (name e category)
type agencyBasic struct {
	Name           string
	FullName       string
	AgencyCategory string
}

// Employee - Represents an employee and his/her salary info
type employee struct {
	Name   string
	Wage   float64
	Perks  float64
	Others float64
	Total  float64
	Type   string
	Active bool
}

// AgencySummary - Summary of an agency
type agencySummary struct {
	FullName          string
	TotalEmployees    int
	TotalWage         float64
	TotalPerks        float64
	MaxWage           float64
	CrawlingTime      *timestamppb.Timestamp
	AgencyName        string
	TotalMembers      int
	TotalServants     int
	TotalInactives    int
	MaxPerk           float64
	TotalRemuneration float64
	HasNext           bool
	HasPrevious       bool
}

type v2AgencySummary struct {
	Agency             string    `json:"orgao"`
	BaseRemuneration   float64   `json:"remuneracao_base"`
	MaxBase            float64   `json:"max_remuneracao_base"`
	OtherRemunerations float64   `json:"outras_remuneracoes"`
	MaxOther           float64   `json:"max_outras_remuneracoes"`
	CrawlingTime       timestamp `json:"timestamp"`
	TotalMembers       int       `json:"total_membros"`
	TotalRemuneration  float64   `json:"total_remuneracao"`
	HasNext            bool      `json:"tem_proximo"`
	HasPrevious        bool      `json:"tem_anterior"`
}

// AgencyTotalsYear - Represents the totals of an year
type agencyTotalsYear struct {
	Year           int
	Agency         *models.Agency
	MonthTotals    []monthTotals
	AgencyFullName string
	SummaryPackage *models.Package `json:"SummaryPackage,omitempty"`
}

type procError struct {
	Stdout string `protobuf:"bytes,2,opt,name=stdout,proto3" json:"stdout,omitempty"` // String containing the standard output of the process.
	Stderr string `protobuf:"bytes,3,opt,name=stderr,proto3" json:"stderr,omitempty"` // String containing the standard error of the process.
}

// MonthTotals - Detailed info of a month (wage, perks, other)
type monthTotals struct {
	Error              *procError
	Month              int
	TotalMembers       int
	BaseRemuneration   float64
	OtherRemunerations float64
	CrawlingTimestamp  *timestamppb.Timestamp
}

// ProcInfoResult - contains information of the result of the process if something went wrong during parsing or crawling process
type procInfoResult struct {
	ProcInfo          *coleta.ProcInfo
	CrawlingTimestamp *timestamppb.Timestamp
}

type v2ProcInfoResult struct {
	ProcInfo  *procInfo  `json:"proc_info"`
	Timestamp *timestamp `json:"timestamp"`
}

type procInfo struct {
	Stdin  string   `json:"stdin,omitempty"`
	Stdout string   `json:"stdout,omitempty"`
	Stderr string   `json:"stderr,omitempty"`
	Cmd    string   `json:"cmd,omitempty"`
	CmdDir string   `json:"cmd_dir,omitempty"`
	Status int32    `json:"status,omitempty"`
	Env    []string `json:"env,omitempty"`
}

type timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

// Os campos que serão trazido pela query de pesquisa
type searchDetails struct {
	Descontos int    `db:"descontos" json:"descontos"`
	Base      int    `db:"base" json:"base"`
	Outras    int    `db:"outras" json:"outras"`
	Orgao     string `db:"orgao" json:"orgao"`
	Mes       int    `db:"mes" json:"mes"`
	Ano       int    `db:"ano" json:"ano"`
	ZipUrl    string `db:"zip_url" json:"zip_url"`
}

type searchResult struct {
	Orgao                    string  `db:"orgao" json:"orgao" csv:"orgao" tableheader:"orgao"`
	Mes                      int     `db:"mes" json:"mes" csv:"mes" tableheader:"mes"`
	Ano                      int     `db:"ano" json:"ano" csv:"ano" tableheader:"ano"`
	Matricula                *string `db:"matricula" json:"matricula" csv:"matricula" tableheader:"matricula"`
	Nome                     string  `db:"nome" json:"nome" csv:"nome" tableheader:"nome"`
	Cargo                    *string `db:"cargo" json:"cargo" csv:"cargo" tableheader:"cargo"`
	Lotacao                  *string `db:"lotacao" json:"lotacao" csv:"lotacao" tableheader:"lotacao"`
	CategoriaContracheque    string  `db:"categoria_contracheque" json:"categoria_contracheque" csv:"categoria_contracheque" tableheader:"categoria_contracheque"`
	DetalhamentoContracheque string  `db:"detalhamento_contracheque" json:"detalhamento_contracheque" csv:"detalhamento_contracheque" tableheader:"detalhamento_contracheque"`
	Valor                    float64 `db:"valor" json:"valor" csv:"valor" tableheader:"valor"`
}

// A resposta que será enviada pela rota de pesquisa
type searchResponse struct {
	DownloadAvailable  bool           `json:"download_available"`
	NumRowsIfAvailable int            `json:"num_rows_if_available"`
	SearchLimit        int            `json:"search_limit"`
	DownloadLimit      int            `json:"download_limit"`
	Results            []searchResult `json:"result"`
}

type agency struct {
	ID            string       `json:"id_orgao,omitempty"`   // 'trt13'
	Name          string       `json:"nome,omitempty"`       // 'Tribunal Regional do Trabalho 13° Região'
	Type          string       `json:"jurisdicao,omitempty"` // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity        string       `json:"entidade,omitempty"`   // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF            string       `json:"uf,omitempty"`         // Short code for federative unity.
	URL           string       `json:"url,omitempty"`        // Link for state url
	Collecting    []collecting `json:"coletando,omitempty"`
	TwitterHandle string       `json:"twitter_handle,omitempty"` // Agency's twitter handle
	OmbudsmanURL  string       `json:"ouvidoria,omitempty"`      //Agencys's ombudsman url
}

// Collecting A Struct containing the day we checked the status of the data and the reasons why we didn't collected it.
type collecting struct {
	Timestamp   *int64   `json:"timestamp,omitempty"` // Day(unix) we checked the status of the data
	Description []string `json:"descricao,omitempty"` // Reasons why we didn't collect the data
}

type annualSummary struct {
	Agency *agency             `json:"orgao,omitempty"`
	Data   []annualSummaryData `json:"dados_anuais,omitempty"`
}

type annualSummaryData struct {
	Year               int            `json:"ano,omitempty"`
	Count              int            `json:"num_membros,omitempty"`
	BaseRemuneration   float64        `json:"remuneracao_base"`
	OtherRemunerations float64        `json:"outras_remuneracoes"`
	Package            *models.Backup `json:"package,omitempty"`
}

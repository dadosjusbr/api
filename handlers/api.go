package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"github.com/labstack/echo"
)

type ApiHandler struct {
	Client         storage.Client
	DadosJusURL    string
	PackageRepoURL string
}

func (a ApiHandler) GetAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	agency, err := a.Client.Db.GetAgency(agencyName)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Agency not found")
	}
	host := c.Request().Host
	agency.URL = fmt.Sprintf("%s/v1/orgao/%s", host, agency.ID)
	return c.JSON(http.StatusFound, agency)
}

func (a ApiHandler) GetAllAgencies(c echo.Context) error {
	agencies, err := a.Client.Db.GetAllAgencies()
	if err != nil {
		fmt.Println("Error while listing agencies: %w", err)
		return c.JSON(http.StatusInternalServerError, "Error while listing agencies")
	}
	host := c.Request().Host
	for i := range agencies {
		agencies[i].URL = fmt.Sprintf("%s/v1/orgao/%s", host, agencies[i].ID)
	}
	return c.JSON(http.StatusOK, agencies)
}

func (a ApiHandler) GetMonthlyInfo(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	agencyName := strings.ToLower(c.Param("orgao"))
	var monthlyInfo map[string][]models.AgencyMonthlyInfo
	month := c.Param("mes")
	if month != "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mes=%d inválido", m))
		}
		oma, _, err := a.Client.GetOMA(m, year, agencyName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error getting OMA data"))
		}
		monthlyInfo = map[string][]models.AgencyMonthlyInfo{
			agencyName: {*oma},
		}
	} else {
		monthlyInfo, err = a.Client.Db.GetMonthlyInfo([]models.Agency{{ID: agencyName}}, year)
	}
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, agencyName))
	}

	if len(monthlyInfo[agencyName]) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	// we recreate all struct to controll the serialization better
	type Backup struct {
		URL  string `json:"url,omitempty"`
		Hash string `json:"hash,omitempty"`
		Size int64  `json:"size,omitempty"`
	}
	type DataSummary struct {
		Max     float64 `json:"max,omitempty"`
		Min     float64 `json:"min,omitempty"`
		Average float64 `json:"media,omitempty"`
		Total   float64 `json:"total,omitempty"`
	}
	type Summary struct {
		Count              int         `json:"quantidade,omitempty"`
		BaseRemuneration   DataSummary `json:"remuneracao_base,omitempty"`
		OtherRemunerations DataSummary `json:"outras_remuneracoes,omitempty"`
	}
	type Summaries struct {
		MemberActive Summary `json:"membros_ativos,omitempty"`
	}
	type Metadata struct {
		OpenFormat       bool   `json:"formato_aberto"`
		Access           string `json:"acesso,omitempty"`
		Extension        string `json:"extensao,omitempty"`
		StrictlyTabular  bool   `json:"dados_estritamente_tabulares"`
		ConsistentFormat bool   `json:"manteve_consistencia_no_formato"`
		HasEnrollment    bool   `json:"tem_matricula"`
		HasCapacity      bool   `json:"tem_lotacao"`
		HasPosition      bool   `json:"tem_cargo"`
		BaseRevenue      string `json:"remuneracao_basica,omitempty"`
		OtherRecipes     string `json:"outras_receitas,omitempty"`
		Expenditure      string `json:"despesas,omitempty"`
	}
	type Score struct {
		Score             float64 `json:"indice_transparencia"`
		CompletenessScore float64 `json:"indice_completude"`
		EasinessScore     float64 `json:"indice_facilidade"`
	}
	type Collect struct {
		Duration       float64 `json:"duracao_segundos,omitempty"`
		CrawlerRepo    string  `json:"repositorio_coletor,omitempty"`
		CrawlerVersion string  `json:"versao_coletor,omitempty"`
		ParserRepo     string  `json:"repositorio_parser,omitempty"`
		ParserVersion  string  `json:"versao_parser,omitempty"`
	}

	type MIError struct {
		ErrorMessage string `json:"err_msg,omitempty"`
		Status       int32  `json:"status,omitempty"`
		Cmd          string `json:"cmd,omitempty"`
	}
	type SummaryzedMI struct {
		AgencyID string     `json:"id_orgao,omitempty"`
		Month    int        `json:"mes,omitempty"`
		Year     int        `json:"ano,omitempty"`
		Summary  *Summaries `json:"sumarios,omitempty"`
		Package  *Backup    `json:"pacote_de_dados,omitempty"`
		Meta     *Metadata  `json:"metadados,omitempty`
		Score    *Score     `json:"indice_transparencia,omitempty`
		Collect  *Collect   `json:"dados_coleta,omitempty`
		Error    *MIError   `json:"error,omitempty"`
	}
	var summaryzedMI []SummaryzedMI
	for i := range monthlyInfo {
		for _, mi := range monthlyInfo[i] {
			if mi.ProcInfo.String() == "" {
				summaryzedMI = append(
					summaryzedMI,
					SummaryzedMI{
						AgencyID: mi.AgencyID,
						Error:    nil,
						Month:    mi.Month,
						Year:     mi.Year,
						Package: &Backup{
							URL:  a.formatDownloadUrl(mi.Package.URL),
							Hash: mi.Package.Hash,
							Size: mi.Package.Size,
						},
						Summary: &Summaries{
							MemberActive: Summary{
								Count: mi.Summary.Count,
								BaseRemuneration: DataSummary{
									Max:     mi.Summary.BaseRemuneration.Max,
									Min:     mi.Summary.BaseRemuneration.Min,
									Average: mi.Summary.BaseRemuneration.Average,
									Total:   mi.Summary.BaseRemuneration.Total,
								},
								OtherRemunerations: DataSummary{
									Max:     mi.Summary.OtherRemunerations.Max,
									Min:     mi.Summary.OtherRemunerations.Min,
									Average: mi.Summary.OtherRemunerations.Average,
									Total:   mi.Summary.OtherRemunerations.Total,
								},
							},
						},
						Meta: &Metadata{
							OpenFormat:       mi.Meta.OpenFormat,
							Access:           mi.Meta.Access,
							Extension:        mi.Meta.Extension,
							StrictlyTabular:  mi.Meta.StrictlyTabular,
							ConsistentFormat: mi.Meta.ConsistentFormat,
							HasEnrollment:    mi.Meta.HaveEnrollment,
							HasCapacity:      mi.Meta.ThereIsACapacity,
							HasPosition:      mi.Meta.HasPosition,
							BaseRevenue:      mi.Meta.BaseRevenue,
							OtherRecipes:     mi.Meta.OtherRecipes,
							Expenditure:      mi.Meta.Expenditure,
						},
						Score: &Score{
							Score:             mi.Score.Score,
							CompletenessScore: mi.Score.CompletenessScore,
							EasinessScore:     mi.Score.EasinessScore,
						},
						Collect: &Collect{
							Duration:       mi.Duration,
							CrawlerRepo:    mi.CrawlerRepo,
							CrawlerVersion: mi.CrawlerVersion,
							ParserRepo:     mi.ParserRepo,
							ParserVersion:  mi.ParserVersion,
						}})
				// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
			} else if mi.ProcInfo.Status != 4 {
				summaryzedMI = append(
					summaryzedMI,
					SummaryzedMI{
						AgencyID: mi.AgencyID,
						Error: &MIError{
							ErrorMessage: mi.ProcInfo.Stderr,
							Status:       mi.ProcInfo.Status,
							Cmd:          mi.ProcInfo.Cmd,
						},
						Month:   mi.Month,
						Year:    mi.Year,
						Package: nil,
						Summary: nil,
						Meta:    nil})
			}
		}
	}
	return c.JSON(http.StatusOK, summaryzedMI)
}

func (a ApiHandler) formatDownloadUrl(url string) string {
	return strings.Replace(url, a.PackageRepoURL, a.DadosJusURL, -1)
}

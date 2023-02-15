package papi

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

type handler struct {
	client         *storage.Client
	dadosJusURL    string
	packageRepoURL string
}

func NewHandler(client *storage.Client, dadosJusURL, packageRepoURL string) *handler {
	return &handler{
		client:         client,
		dadosJusURL:    dadosJusURL,
		packageRepoURL: packageRepoURL,
	}
}

func (h handler) V1GetAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	agency, err := h.client.Db.GetAgency(agencyName)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Agency not found")
	}
	host := c.Request().Host
	agency.URL = fmt.Sprintf("%s/v1/orgao/%s", host, agency.ID)
	return c.JSON(http.StatusOK, agency)
}

func (h handler) V2GetAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	strAgency, err := h.client.Db.GetAgency(agencyName)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Agency not found")
	}
	var collect []collecting
	for _, c := range strAgency.Collecting {
		collect = append(collect, collecting{
			Timestamp:   c.Timestamp,
			Description: c.Description,
		})
	}
	host := c.Request().Host
	url := fmt.Sprintf("%s/v2/orgao/%s", host, strAgency.ID)
	agency := &agency{
		ID:            strAgency.ID,
		Name:          strAgency.Name,
		Type:          strAgency.Type,
		Entity:        strAgency.Entity,
		UF:            strAgency.UF,
		URL:           url,
		Collecting:    collect,
		TwitterHandle: strAgency.TwitterHandle,
		OmbudsmanURL:  strAgency.OmbudsmanURL,
	}
	return c.JSON(http.StatusOK, agency)
}

func (h handler) V1GetAllAgencies(c echo.Context) error {
	agencies, err := h.client.Db.GetAllAgencies()
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

func (h handler) V2GetAllAgencies(c echo.Context) error {
	strAgencies, err := h.client.Db.GetAllAgencies()
	if err != nil {
		fmt.Println("Error while listing agencies: %w", err)
		return c.JSON(http.StatusInternalServerError, "Error while listing agencies")
	}
	agencies := []agency{}
	host := c.Request().Host
	for _, a := range strAgencies {
		var collect []collecting
		for _, c := range a.Collecting {
			collect = append(collect, collecting{
				Timestamp:   c.Timestamp,
				Description: c.Description,
			})
		}
		url := fmt.Sprintf("%s/v2/orgao/%s", host, a.ID)
		agency := agency{
			ID:            a.ID,
			Name:          a.Name,
			Type:          a.Type,
			Entity:        a.Entity,
			UF:            a.UF,
			URL:           url,
			Collecting:    collect,
			TwitterHandle: a.TwitterHandle,
			OmbudsmanURL:  a.OmbudsmanURL,
		}
		agencies = append(agencies, agency)
	}
	return c.JSON(http.StatusOK, agencies)
}

func (h handler) GetMonthlyInfo(c echo.Context) error {
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
		oma, _, err := h.client.GetOMA(m, year, agencyName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Error getting OMA data")
		}
		monthlyInfo = map[string][]models.AgencyMonthlyInfo{
			agencyName: {*oma},
		}
	} else {
		monthlyInfo, err = h.client.Db.GetMonthlyInfo([]models.Agency{{ID: agencyName}}, year)
	}
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, agencyName))
	}

	if len(monthlyInfo[agencyName]) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	var sumMI []summaryzedMI
	for i := range monthlyInfo {
		for _, mi := range monthlyInfo[i] {
			// Fazemos duas checagens no formato do ProcInfo para saber se ele é vazio pois alguns dados diferem, no banco de dados, quando o procinfo é nulo.
			if mi.ProcInfo == nil || mi.ProcInfo.String() == "" {
				sumMI = append(
					sumMI,
					summaryzedMI{
						AgencyID: mi.AgencyID,
						Error:    nil,
						Month:    mi.Month,
						Year:     mi.Year,
						Package: &backup{
							URL:  h.formatDownloadUrl(mi.Package.URL),
							Hash: mi.Package.Hash,
							Size: mi.Package.Size,
						},
						Summary: &summaries{
							MemberActive: summary{
								Count: mi.Summary.Count,
								BaseRemuneration: dataSummary{
									Max:     mi.Summary.BaseRemuneration.Max,
									Min:     mi.Summary.BaseRemuneration.Min,
									Average: mi.Summary.BaseRemuneration.Average,
									Total:   mi.Summary.BaseRemuneration.Total,
								},
								OtherRemunerations: dataSummary{
									Max:     mi.Summary.OtherRemunerations.Max,
									Min:     mi.Summary.OtherRemunerations.Min,
									Average: mi.Summary.OtherRemunerations.Average,
									Total:   mi.Summary.OtherRemunerations.Total,
								},
							},
						},
						Metadata: &metadata{
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
						Score: &score{
							Score:             mi.Score.Score,
							CompletenessScore: mi.Score.CompletenessScore,
							EasinessScore:     mi.Score.EasinessScore,
						},
						Collect: &collect{
							Duration:       mi.Duration,
							CrawlerRepo:    mi.CrawlerRepo,
							CrawlerVersion: mi.CrawlerVersion,
							ParserRepo:     mi.ParserRepo,
							ParserVersion:  mi.ParserVersion,
						}})
				// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
			} else if mi.ProcInfo.Status != 4 {
				sumMI = append(
					sumMI,
					summaryzedMI{
						AgencyID: mi.AgencyID,
						Error: &miError{
							ErrorMessage: mi.ProcInfo.Stderr,
							Status:       mi.ProcInfo.Status,
							Cmd:          mi.ProcInfo.Cmd,
						},
						Month:    mi.Month,
						Year:     mi.Year,
						Package:  nil,
						Summary:  nil,
						Metadata: nil})
			}
		}
	}
	return c.JSON(http.StatusOK, sumMI)
}

func (h handler) formatDownloadUrl(url string) string {
	return strings.Replace(url, h.packageRepoURL, h.dadosJusURL, -1)
}

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

func (h handler) GetAllAgencies(c echo.Context) error {
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

func (h handler) V1GetMonthlyInfo(c echo.Context) error {
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
		oma, _, err := h.client.Db.GetOMA(m, year, agencyName)
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

func (h handler) V2GetMonthlyInfo(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}

	agencyName := strings.ToLower(c.Param("orgao"))
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mes=%d inválido", month))
	}

	var monthlyInfo *models.AgencyMonthlyInfo
	monthlyInfo, _, err = h.client.Db.GetOMA(month, year, agencyName)
	if err != nil {
		if err.Error() == "there is no data with this parameters" {
			return c.JSON(http.StatusNotFound, "Não existem dados para os parâmetros informados")
		}
		return c.JSON(http.StatusBadRequest, "Error getting OMA data")
	}

	var sumMI summaryzedMI
	if monthlyInfo.ProcInfo.String() == "" {
		sumMI =
			summaryzedMI{
				AgencyID: monthlyInfo.AgencyID,
				Error:    nil,
				Month:    monthlyInfo.Month,
				Year:     monthlyInfo.Year,
				Package: &backup{
					URL:  h.formatDownloadUrl(monthlyInfo.Package.URL),
					Hash: monthlyInfo.Package.Hash,
					Size: monthlyInfo.Package.Size,
				},
				Summary: &summaries{
					MemberActive: summary{
						Count: monthlyInfo.Summary.Count,
						BaseRemuneration: dataSummary{
							Max:     monthlyInfo.Summary.BaseRemuneration.Max,
							Min:     monthlyInfo.Summary.BaseRemuneration.Min,
							Average: monthlyInfo.Summary.BaseRemuneration.Average,
							Total:   monthlyInfo.Summary.BaseRemuneration.Total,
						},
						OtherRemunerations: dataSummary{
							Max:     monthlyInfo.Summary.OtherRemunerations.Max,
							Min:     monthlyInfo.Summary.OtherRemunerations.Min,
							Average: monthlyInfo.Summary.OtherRemunerations.Average,
							Total:   monthlyInfo.Summary.OtherRemunerations.Total,
						},
					},
				},
				Metadata: &metadata{
					OpenFormat:       monthlyInfo.Meta.OpenFormat,
					Access:           monthlyInfo.Meta.Access,
					Extension:        monthlyInfo.Meta.Extension,
					StrictlyTabular:  monthlyInfo.Meta.StrictlyTabular,
					ConsistentFormat: monthlyInfo.Meta.ConsistentFormat,
					HasEnrollment:    monthlyInfo.Meta.HaveEnrollment,
					HasCapacity:      monthlyInfo.Meta.ThereIsACapacity,
					HasPosition:      monthlyInfo.Meta.HasPosition,
					BaseRevenue:      monthlyInfo.Meta.BaseRevenue,
					OtherRecipes:     monthlyInfo.Meta.OtherRecipes,
					Expenditure:      monthlyInfo.Meta.Expenditure,
				},
				Score: &score{
					Score:             monthlyInfo.Score.Score,
					CompletenessScore: monthlyInfo.Score.CompletenessScore,
					EasinessScore:     monthlyInfo.Score.EasinessScore,
				},
				Collect: &collect{
					Duration:       monthlyInfo.Duration,
					CrawlerRepo:    monthlyInfo.CrawlerRepo,
					CrawlerVersion: monthlyInfo.CrawlerVersion,
					ParserRepo:     monthlyInfo.ParserRepo,
					ParserVersion:  monthlyInfo.ParserVersion,
				},
			}
		//O status 4 informa que os dados estão indisponíveis. Ao removê-los dos resultados da API, garantimos que eles sejam exibidos como se não houvesse dados.
	} else if monthlyInfo.ProcInfo.Status != 4 {
		sumMI = summaryzedMI{
			AgencyID: monthlyInfo.AgencyID,
			Error: &miError{
				ErrorMessage: monthlyInfo.ProcInfo.Stderr,
				Status:       monthlyInfo.ProcInfo.Status,
				Cmd:          monthlyInfo.ProcInfo.Cmd,
			},
			Month:    monthlyInfo.Month,
			Year:     monthlyInfo.Year,
			Package:  nil,
			Summary:  nil,
			Metadata: nil,
		}
	} else {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, sumMI)
}

func (h handler) GetMonthlyInfosByYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}

	agencyName := strings.ToLower(c.Param("orgao"))
	var monthlyInfo map[string][]models.AgencyMonthlyInfo
	monthlyInfo, err = h.client.Db.GetMonthlyInfo([]models.Agency{{ID: agencyName}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, agencyName))
	}

	if len(monthlyInfo[agencyName]) == 0 {
		return c.JSON(http.StatusNotFound, "Não existem dados para os parâmetros informados")
	}

	var sumMI []summaryzedMI
	for i := range monthlyInfo {
		for _, mi := range monthlyInfo[i] {
			if mi.ProcInfo.String() == "" {
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
				//O status 4 informa que os dados estão indisponíveis. Ao removê-los dos resultados da API, garantimos que eles sejam exibidos como se não houvesse dados.
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

package papi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"github.com/labstack/echo/v4"
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

//	@ID				GetAgencyById
//	@Tags			public_api
//	@Description	Busca um órgão específico utilizando seu ID.
//	@Produce		json
//	@Param			orgao				path		string	true	"ID do órgão. Exemplos: tjal, tjba, mppb."
//	@Success		200					{object}	agency	"Requisição bem sucedida."
//	@Failure		404					{string}	string	"Órgão não encontrado."
//	@Router			/v2/orgao/{orgao} 	[get]
func (h handler) V2GetAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	strAgency, err := h.client.Db.GetAgency(agencyName)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Agency not found")
	}
	var collect []collecting
	var hasData bool
	for _, c := range strAgency.Collecting {
		collect = append(collect, collecting{
			Timestamp:   c.Timestamp,
			Description: c.Description,
		})
		hasData = c.Collecting
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
		HasData:       hasData,
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

//	@ID				GetAllAgencies
//	@Tags			public_api
//	@Description	Busca todos os órgãos disponíveis.
//	@Produce		json
//	@Success		200			{object}	[]agency	"Requisição bem sucedida."
//	@Failure		500			{string}	string		"Erro interno do servidor."
//	@Router			/v2/orgaos 	[get]
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
		var hasData bool
		for _, c := range a.Collecting {
			collect = append(collect, collecting{
				Timestamp:   c.Timestamp,
				Description: c.Description,
			})
			hasData = c.Collecting
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
			HasData:       hasData,
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
								Discounts: dataSummary{
									Max:     mi.Summary.Discounts.Max,
									Min:     mi.Summary.Discounts.Min,
									Average: mi.Summary.Discounts.Average,
									Total:   mi.Summary.Discounts.Total,
								},
								Remunerations: dataSummary{
									Max:     mi.Summary.Remunerations.Max,
									Min:     mi.Summary.Remunerations.Min,
									Average: mi.Summary.Remunerations.Average,
									Total:   mi.Summary.Remunerations.Total,
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

//	@ID				GetMonthlyInfo
//	@Tags			public_api
//	@Description	Busca um dado mensal de um órgão
//	@Produce		json
//	@Success		200		{object}	summaryzedMI	"Requisição bem sucedida"
//	@Failure		400		{string}	string			"Parâmetros inválidos"
//	@Failure		404		{string}	string			"Não existem dados para os parâmetros informados"
//	@Param			ano		path		int				true	"Ano"
//	@Param			orgao	path		string			true	"Órgão"
//	@Param			mes		path		int				true	"Mês"
//	@Router			/v2/dados/{orgao}/{ano}/{mes} [get]
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
						Discounts: dataSummary{
							Max:     monthlyInfo.Summary.Discounts.Max,
							Min:     monthlyInfo.Summary.Discounts.Min,
							Average: monthlyInfo.Summary.Discounts.Average,
							Total:   monthlyInfo.Summary.Discounts.Total,
						},
						Remunerations: dataSummary{
							Max:     monthlyInfo.Summary.Remunerations.Max,
							Min:     monthlyInfo.Summary.Remunerations.Min,
							Average: monthlyInfo.Summary.Remunerations.Average,
							Total:   monthlyInfo.Summary.Remunerations.Total,
						},
						ItemSummary: itemSummary{
							FoodAllowance:        monthlyInfo.Summary.ItemSummary.FoodAllowance,
							BonusLicense:         monthlyInfo.Summary.ItemSummary.BonusLicense,
							VacationCompensation: monthlyInfo.Summary.ItemSummary.VacationCompensation,
							ChristmasBonus:       monthlyInfo.Summary.ItemSummary.ChristmasBonus,
							CompensatoryLicense:  monthlyInfo.Summary.ItemSummary.CompensatoryLicense,
							Others:               monthlyInfo.Summary.ItemSummary.Others,
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

//	@ID				GetMonthlyInfosByYear
//	@Tags			public_api
//	@Description	Busca os dados mensais de um órgão por ano
//	@Produce		json
//	@Success		200		{object}	[]summaryzedMI	"Requisição bem sucedida"
//	@Failure		400		{string}	string			"Parâmetros inválidos"
//	@Failure		404		{string}	string			"Não existem dados para os parâmetros informados"
//	@Param			ano		path		int				true	"Ano"
//	@Param			orgao	path		string			true	"Órgão"
//	@Router			/v2/dados/{orgao}/{ano} [get]
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
								Discounts: dataSummary{
									Max:     mi.Summary.Discounts.Max,
									Min:     mi.Summary.Discounts.Min,
									Average: mi.Summary.Discounts.Average,
									Total:   mi.Summary.Discounts.Total,
								},
								Remunerations: dataSummary{
									Max:     mi.Summary.Remunerations.Max,
									Min:     mi.Summary.Remunerations.Min,
									Average: mi.Summary.Remunerations.Average,
									Total:   mi.Summary.Remunerations.Total,
								},
								ItemSummary: itemSummary{
									FoodAllowance:        mi.Summary.ItemSummary.FoodAllowance,
									BonusLicense:         mi.Summary.ItemSummary.BonusLicense,
									VacationCompensation: mi.Summary.ItemSummary.VacationCompensation,
									ChristmasBonus:       mi.Summary.ItemSummary.ChristmasBonus,
									CompensatoryLicense:  mi.Summary.ItemSummary.CompensatoryLicense,
									Others:               mi.Summary.ItemSummary.Others,
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

//	@ID				GetAggregateIndexesWithParams
//	@Tags			public_api
//	@Description	Busca as informações de índices de um grupo ou órgão específico.
//	@Produce		json
//	@Success		200							{object}	[]aggregateIndexes	"Requisição bem sucedida."
//	@Failure		400							{string}	string				"Requisição inválida."
//	@Failure		500							{string}	string				"Erro interno do servidor."
//	@Param			param						path		string				true	"'grupo' ou 'orgao'"
//	@Param			valor						path		string				true	"Jurisdição ou ID do órgao"
//	@Router			/v2/indice/{param}/{valor} 	[get]
func (h handler) V2GetAggregateIndexesWithParams(c echo.Context) error {
	param := c.Param("param")
	valor := c.Param("valor")
	ano := c.Param("ano")
	mes := c.Param("mes")
	agregado := c.QueryParam("agregado")
	detalhe := c.QueryParam("detalhe")

	groupMap := map[string]string{
		"justica-eleitoral":    "Eleitoral",
		"ministerios-publicos": "Ministério",
		"justica-estadual":     "Estadual",
		"justica-do-trabalho":  "Trabalho",
		"justica-federal":      "Federal",
		"justica-militar":      "Militar",
		"justica-superior":     "Superior",
		"conselhos-de-justica": "Conselho",
	}

	// porJurisdicao tbm será usada para verificar a possibilidade de uma BadRequest
	var porJurisdicao bool

	// Verificamos se o parâmetro é válido.
	if param == "grupo" {
		if _, porJurisdicao = groupMap[valor]; porJurisdicao {
			valor = groupMap[valor]
		} else {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Jurisdição inválida: %s.", valor))
		}
	} else if param != "orgao" && param != "" {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro inválido: %s.", param))
	}

	// Personalizando a mensagem de saída de acordo com o parâmetro.
	msg := map[string]string{
		"grupo": fmt.Sprintf("para o grupo: %s.", valor),
		"orgao": fmt.Sprintf("para o órgão: %s.", valor),
		"":      "para todos os órgãos.",
	}

	var anoInt, mesInt int
	var err error

	// Verificamos se ano e mês foram informados e se são válidos (convertemos para inteiro)
	if ano != "" {
		anoInt, err = strconv.Atoi(ano)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ANO inválido: %s.", ano))
		}
		if mes != "" {
			mesInt, err = strconv.Atoi(mes)
			if err != nil {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro MÊS inválido: %s.", mes))
			}
		}
	}

	var indexes map[string][]models.IndexInformation

	if ano != "" && mes != "" {
		// Caso o ano e o mês sejam informados
		indexes, err = h.client.Db.GetIndexInformation(valor, mesInt, anoInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro consultando os índices de %d/%d %s", mesInt, anoInt, msg[param]))
		}
	} else if ano != "" {
		// Caso apenas o ano seja informado
		indexes, err = h.client.Db.GetIndexInformation(valor, 0, anoInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro consultando os índices de %d %s", anoInt, msg[param]))
		}
	} else {
		// Caso nem ano ou mês tenham sido informados
		indexes, err = h.client.Db.GetIndexInformation(valor, 0, 0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro consultando os índices %s", msg[param]))
		}
	}
	if _, ok := indexes[valor]; valor != "" && !ok && !porJurisdicao {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro consultando os índices. Órgão/grupo inválido: %s", valor))
	}

	indexMap := make(map[string][]indexInformation)
	aggregateScore := make(map[string]float64)
	aggregateEasinessScore := make(map[string]float64)
	aggregateCompletenessScore := make(map[string]float64)

	for id, index := range indexes {
		for _, a := range index {
			var meta *metadata
			// Se "detalhe=true" não for passada na URL, os metadados não serão passados
			if detalhe == "true" {
				meta = &metadata{
					OpenFormat:       a.Meta.OpenFormat,
					Access:           a.Meta.Access,
					Extension:        a.Meta.Extension,
					StrictlyTabular:  a.Meta.StrictlyTabular,
					ConsistentFormat: a.Meta.ConsistentFormat,
					HasEnrollment:    a.Meta.HaveEnrollment,
					HasCapacity:      a.Meta.ThereIsACapacity,
					HasPosition:      a.Meta.HasPosition,
					BaseRevenue:      a.Meta.BaseRevenue,
					OtherRecipes:     a.Meta.OtherRecipes,
					Expenditure:      a.Meta.Expenditure,
				}
			}
			if agregado != "true" {
				indexMap[id] = append(indexMap[id], indexInformation{
					Month: a.Month,
					Year:  a.Year,
					Score: &score{
						Score:             a.Score.Score,
						EasinessScore:     a.Score.EasinessScore,
						CompletenessScore: a.Score.CompletenessScore,
					},
					Metadata: meta,
				})
			}
			aggregateScore[id] += a.Score.Score
			aggregateCompletenessScore[id] += a.Score.CompletenessScore
			aggregateEasinessScore[id] += a.Score.EasinessScore
		}
	}
	var aggregate []aggregateIndexes
	for id, index := range indexes {
		aggregateScore[id] = aggregateScore[id] / float64(len(index))
		aggregateEasinessScore[id] = aggregateEasinessScore[id] / float64(len(index))
		aggregateCompletenessScore[id] = aggregateCompletenessScore[id] / float64(len(index))

		agg := aggregateIndexes{
			ID: id,
			Aggregate: &score{
				Score:             aggregateScore[id],
				EasinessScore:     aggregateEasinessScore[id],
				CompletenessScore: aggregateCompletenessScore[id],
			},
		}
		// Se "agregado=true" não estiver presente na URL, será listado também o detalhamento dos índices do órgão
		if agregado != "true" {
			agg.IndexInformation = indexMap[id]
		}

		aggregate = append(aggregate, agg)
	}
	return c.JSON(http.StatusOK, aggregate)
}

//	@ID				GetAggregateIndexes
//	@Tags			public_api
//	@Description	Busca as informações de índices de todos os órgãos.
//	@Produce		json
//	@Success		200			{object}	[]aggregateIndexesByGroup	"Requisição bem sucedida."
//	@Failure		500			{string}	string						"Erro interno do servidor."
//	@Router			/v2/indice 																																																													[get]
func (h handler) V2GetAggregateIndexes(c echo.Context) error {
	agregado := c.QueryParam("agregado")
	detalhe := c.QueryParam("detalhe")

	indexes, err := h.client.Db.GetIndexInformation("", 0, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Erro consultando os índices para todos os órgãos.")
	}

	indexMap := make(map[string][]indexInformation)
	aggregateScore := make(map[string]float64)
	aggregateEasinessScore := make(map[string]float64)
	aggregateCompletenessScore := make(map[string]float64)

	// Mapearemos jurisdição e respectivos orgãos
	grupos := make(map[string][]string)

	for id, index := range indexes {
		for _, a := range index {
			var meta *metadata
			// Se "detalhe=true" não for passada na URL, os metadados não serão passados
			if detalhe == "true" {
				meta = &metadata{
					OpenFormat:       a.Meta.OpenFormat,
					Access:           a.Meta.Access,
					Extension:        a.Meta.Extension,
					StrictlyTabular:  a.Meta.StrictlyTabular,
					ConsistentFormat: a.Meta.ConsistentFormat,
					HasEnrollment:    a.Meta.HaveEnrollment,
					HasCapacity:      a.Meta.ThereIsACapacity,
					HasPosition:      a.Meta.HasPosition,
					BaseRevenue:      a.Meta.BaseRevenue,
					OtherRecipes:     a.Meta.OtherRecipes,
					Expenditure:      a.Meta.Expenditure,
				}
			}
			if agregado != "true" {
				indexMap[id] = append(indexMap[id], indexInformation{
					Month: a.Month,
					Year:  a.Year,
					Score: &score{
						Score:             a.Score.Score,
						EasinessScore:     a.Score.EasinessScore,
						CompletenessScore: a.Score.CompletenessScore,
					},
					Metadata: meta,
				})
			}

			aggregateScore[id] += a.Score.Score
			aggregateCompletenessScore[id] += a.Score.CompletenessScore
			aggregateEasinessScore[id] += a.Score.EasinessScore

			// Criando a lista de órgãos por jurisdição para filtrar posteriormente
			if !slices.Contains(grupos[a.Type], id) {
				grupos[a.Type] = append(grupos[a.Type], id)
			}
		}
	}
	aggregate := make(map[string]aggregateIndexes)
	for id, index := range indexes {
		aggregateScore[id] = aggregateScore[id] / float64(len(index))
		aggregateEasinessScore[id] = aggregateEasinessScore[id] / float64(len(index))
		aggregateCompletenessScore[id] = aggregateCompletenessScore[id] / float64(len(index))

		agg := aggregateIndexes{
			ID: id,
			Aggregate: &score{
				Score:             aggregateScore[id],
				EasinessScore:     aggregateEasinessScore[id],
				CompletenessScore: aggregateCompletenessScore[id],
			},
		}
		// Se "agregado=true" não estiver presente na URL, será listado também o detalhamento dos índices do órgão
		if agregado != "true" {
			agg.IndexInformation = indexMap[id]
		}

		aggregate[id] = agg
	}

	// Aqui realizamos o filtro, adicionando o agregado de cada órgão ao seu respectivo grupo.
	dadosGrupo := make(map[string][]aggregateIndexes)
	for grupo, orgaos := range grupos {
		for _, orgao := range orgaos {
			dadosGrupo[grupo] = append(dadosGrupo[grupo], aggregate[orgao])
		}
	}

	dados := aggregateIndexesByGroup{
		JusticaEstadual:  dadosGrupo["Estadual"],
		JusticaMilitar:   dadosGrupo["Militar"],
		JusticaFederal:   dadosGrupo["Federal"],
		JusticaEleitoral: dadosGrupo["Eleitoral"],
		JusticaSuperior:  dadosGrupo["Superior"],
		Ministerios:      dadosGrupo["Ministério"],
		Conselhos:        dadosGrupo["Conselho"],
		JusticaTrabalho:  dadosGrupo["Trabalho"],
	}

	return c.JSON(http.StatusOK, dados)
}

//	@ID				GetAllAgencyInformation
//	@Tags			public_api
//	@Description	Busca todas as informações de um órgão específico.
//	@Produce		json
//	@Success		200					{object}	allAgencyInformation	"Requisição bem sucedida."
//	@Failure		400					{string}	string					"Requisição inválida."
//	@Param			orgao				path		string					true	"órgão"
//	@Router			/v2/dados/{orgao} 	[get]
func (h handler) V2GetAllAgencyInformation(c echo.Context) error {
	agency := strings.ToLower(c.Param("orgao"))

	ag, err := h.client.Db.GetAgency(agency)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("Órgão não encontrado: %s", strings.ToUpper(agency)))
	}
	collections, err := h.client.Db.GetAllAgencyCollection(agency)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("Não encontramos dados para o órgão %s", strings.ToUpper(agency)))
	}

	aggregateScore := 0.0
	aggregateEasinessScore := 0.0
	aggregateCompletenessScore := 0.0
	numMonthsWithData := 0
	var result []summaryzedMI

	for _, c := range collections {
		if c.ProcInfo == nil || c.ProcInfo.String() == "" {
			result = append(result, summaryzedMI{
				Error: nil,
				Month: c.Month,
				Year:  c.Year,
				Summary: &summaries{
					MemberActive: summary{
						Count: c.Summary.Count,
						BaseRemuneration: dataSummary{
							Max:     c.Summary.BaseRemuneration.Max,
							Min:     c.Summary.BaseRemuneration.Min,
							Average: c.Summary.BaseRemuneration.Average,
							Total:   c.Summary.BaseRemuneration.Total,
						},
						OtherRemunerations: dataSummary{
							Max:     c.Summary.OtherRemunerations.Max,
							Min:     c.Summary.OtherRemunerations.Min,
							Average: c.Summary.OtherRemunerations.Average,
							Total:   c.Summary.OtherRemunerations.Total,
						},
						Discounts: dataSummary{
							Max:     c.Summary.Discounts.Max,
							Min:     c.Summary.Discounts.Min,
							Average: c.Summary.Discounts.Average,
							Total:   c.Summary.Discounts.Total,
						},
						Remunerations: dataSummary{
							Max:     c.Summary.Remunerations.Max,
							Min:     c.Summary.Remunerations.Min,
							Average: c.Summary.Remunerations.Average,
							Total:   c.Summary.Remunerations.Total,
						},
						ItemSummary: itemSummary{
							FoodAllowance:        c.Summary.ItemSummary.FoodAllowance,
							BonusLicense:         c.Summary.ItemSummary.BonusLicense,
							VacationCompensation: c.Summary.ItemSummary.VacationCompensation,
							ChristmasBonus:       c.Summary.ItemSummary.ChristmasBonus,
							CompensatoryLicense:  c.Summary.ItemSummary.CompensatoryLicense,
							Others:               c.Summary.ItemSummary.Others,
						},
					},
				},
				Metadata: &metadata{
					OpenFormat:       c.Meta.OpenFormat,
					Access:           c.Meta.Access,
					Extension:        c.Meta.Extension,
					StrictlyTabular:  c.Meta.StrictlyTabular,
					ConsistentFormat: c.Meta.ConsistentFormat,
					HasEnrollment:    c.Meta.HaveEnrollment,
					HasCapacity:      c.Meta.ThereIsACapacity,
					HasPosition:      c.Meta.HasPosition,
					BaseRevenue:      c.Meta.BaseRevenue,
					OtherRecipes:     c.Meta.OtherRecipes,
					Expenditure:      c.Meta.Expenditure,
				},
				Score: &score{
					Score:             c.Score.Score,
					CompletenessScore: c.Score.CompletenessScore,
					EasinessScore:     c.Score.EasinessScore,
				}})
			numMonthsWithData++
		}
		aggregateScore += c.Score.Score
		aggregateCompletenessScore += c.Score.CompletenessScore
		aggregateEasinessScore += c.Score.EasinessScore
	}

	var collect []collecting
	for _, c := range ag.Collecting {
		collect = append(collect, collecting{
			Timestamp:   c.Timestamp,
			Description: c.Description,
		})
	}

	agencyInfo := allAgencyInformation{
		ID:                ag.ID,
		Name:              ag.Name,
		Type:              ag.Type,
		Entity:            ag.Entity,
		UF:                ag.UF,
		URL:               ag.URL,
		Collecting:        collect,
		TwitterHandle:     ag.TwitterHandle,
		OmbudsmanURL:      ag.OmbudsmanURL,
		TotalCollections:  len(collections),
		NumMonthsWithData: numMonthsWithData,
		Score: &score{
			Score:             aggregateScore / float64(len(collections)),
			EasinessScore:     aggregateEasinessScore / float64(len(collections)),
			CompletenessScore: aggregateCompletenessScore / float64(len(collections)),
		},
		Collections: result,
	}
	return c.JSON(http.StatusOK, agencyInfo)
}

func (h handler) formatDownloadUrl(url string) string {
	return strings.Replace(url, h.packageRepoURL, h.dadosJusURL, -1)
}

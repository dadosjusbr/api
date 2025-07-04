package uiapi

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/storage"
	strModels "github.com/dadosjusbr/storage/models"
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type handler struct {
	client           *storage.Client
	db               *postgresDB
	sess             *awsSession
	s3Bucket         string
	loc              *time.Location
	envOmittedFields []string
	searchLimit      int
	downloadLimit    int
}

func NewHandler(client *storage.Client, conn *gorm.DB, newrelic *newrelic.Application, awsRegion string, s3Bucket string, loc *time.Location, envOmittedFields []string, searchLimit, downloadLimit int) (*handler, error) {
	db := &postgresDB{
		conn:     conn,
		newrelic: newrelic,
	}
	sess, err := newAwsSession(awsRegion)
	if err != nil {
		return nil, err
	}
	return &handler{
		db:               db,
		sess:             sess,
		s3Bucket:         s3Bucket,
		client:           client,
		loc:              loc,
		envOmittedFields: envOmittedFields,
		searchLimit:      searchLimit,
		downloadLimit:    downloadLimit,
	}, nil
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GetSummaryOfAgency(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%d", month))
	}
	agencyName := c.Param("orgao")
	agencyMonthlyInfo, agency, err := h.client.GetOMA(month, year, agencyName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	agencySummary := agencySummary{
		FullName:   agency.Name,
		TotalWage:  agencyMonthlyInfo.Summary.BaseRemuneration.Total,
		MaxWage:    agencyMonthlyInfo.Summary.BaseRemuneration.Max,
		TotalPerks: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
		MaxPerk:    agencyMonthlyInfo.Summary.OtherRemunerations.Max,
		TotalRemuneration: agencyMonthlyInfo.Summary.BaseRemuneration.Total +
			agencyMonthlyInfo.Summary.OtherRemunerations.Total,
		TotalMembers: agencyMonthlyInfo.Summary.Count,
		CrawlingTime: agencyMonthlyInfo.CrawlingTimestamp,
		HasNext:      time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(h.loc).Before(time.Now().AddDate(0, 1, 0)),
		HasPrevious:  time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(h.loc).After(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).In(h.loc)),
	}
	return c.JSON(http.StatusOK, agencySummary)
}

// @ID				GetSummaryOfAgency
// @Tags			ui_api
// @Description	Endpoint de resumo das remunerações mensal de um órgão. Fornece uma análise financeira abrangente para um órgão específico em um mês e ano determinados
// @Description
// @Description	Informações Financeiras Detalhadas:
// @Description	- Remuneração base total e máxima
// @Description	- Outras remunerações e benefícios
// @Description	- Valor total de descontos
// @Description	- Contagem de membros
// @Description	- Análise de rubricas (penduricalhos) específicas
// @Description
// @Description	Contexto Adicional:
// @Description	- Marcadores de existência de dados anteriores/posteriores ao ano/mês consultados
// @Description	- Timestamp da coleta de dados
// @Description	- Detalhamento de diferentes tipos de remuneração (remuneração base, outras remunerações e descontos)
// @Produce		json
// @Param			orgao	path		string			true	"Sigla do órgão. Ex.: tjal, mppb, tjmmg"
// @Param			mes		path		string			true	"Mês de referência (1-12)"
// @Param			ano		path		string			true	"Ano de referência"
// @Success		200		{object}	v2AgencySummary	"Resumo financeiro do órgão processado com sucesso"
// @Failure		400		{string}	string			"Erro de validação dos parâmetros de entrada"
// @Failure		404		{string}	string			"Órgão ou dados não encontrados"
// @Failure		500		{string}	string			"Erro interno durante processamento da consulta"
// @Router			/uiapi/v2/orgao/resumo/{orgao}/{ano}/{mes} [get]
func (h handler) V2GetSummaryOfAgency(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%s inválido", c.Param("ano")))
	}
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%s inválido", c.Param("mes")))
	}
	agencyName := c.Param("orgao")
	agencyMonthlyInfo, agency, err := h.client.Db.GetOMA(month, year, agencyName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	agencySummary := v2AgencySummary{
		Agency:             agency.Name,
		BaseRemuneration:   agencyMonthlyInfo.Summary.BaseRemuneration.Total,
		MaxBase:            agencyMonthlyInfo.Summary.BaseRemuneration.Max,
		OtherRemunerations: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
		MaxOther:           agencyMonthlyInfo.Summary.OtherRemunerations.Max,
		Discounts:          agencyMonthlyInfo.Summary.Discounts.Total,
		MaxDiscounts:       agencyMonthlyInfo.Summary.Discounts.Max,
		MaxRemuneration:    agencyMonthlyInfo.Summary.Remunerations.Max,
		TotalRemuneration:  agencyMonthlyInfo.Summary.Remunerations.Total,
		TotalMembers:       agencyMonthlyInfo.Summary.Count,
		CrawlingTime: timestamp{
			Seconds: agencyMonthlyInfo.CrawlingTimestamp.GetSeconds(),
			Nanos:   agencyMonthlyInfo.CrawlingTimestamp.GetNanos(),
		},
		HasNext:     time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(h.loc).Before(time.Now().AddDate(0, 1, 0)),
		HasPrevious: time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(h.loc).After(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).In(h.loc)),
		ItemSummary: itemSummary(agencyMonthlyInfo.Summary.ItemSummary),
	}
	return c.JSON(http.StatusOK, agencySummary)
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GetSalaryOfAgencyMonthYear(c echo.Context) error {
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%d", month))
	}
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d", year))
	}
	agencyName := strings.ToLower(c.Param("orgao"))
	agencyMonthlyInfo, _, err := h.client.GetOMA(month, year, agencyName)
	if err != nil {
		log.Printf("[salary agency month year] error getting data for second screen(mes:%d ano:%d, orgao:%s):%q", month, year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	if agencyMonthlyInfo.ProcInfo.String() != "" {
		var newEnv = agencyMonthlyInfo.ProcInfo.Env
		for _, omittedField := range h.envOmittedFields {
			for i, field := range newEnv {
				if strings.Contains(field, omittedField) {
					newEnv[i] = omittedField + "= ##omitida##"
					break
				}
			}
		}
		agencyMonthlyInfo.ProcInfo.Env = newEnv
		return c.JSON(http.StatusPartialContent, procInfoResult{
			ProcInfo:          agencyMonthlyInfo.ProcInfo,
			CrawlingTimestamp: agencyMonthlyInfo.CrawlingTimestamp,
		})
	}
	return c.JSON(http.StatusOK, dataForChartAtAgencyScreen{
		Members:     agencyMonthlyInfo.Summary.IncomeHistogram,
		MaxSalary:   agencyMonthlyInfo.Summary.BaseRemuneration.Max,
		PackageURL:  agencyMonthlyInfo.Package.URL,
		PackageHash: agencyMonthlyInfo.Package.Hash,
		PackageSize: agencyMonthlyInfo.Package.Size,
	})
}

// @ID				GetSalaryOfAgencyMonthYear
// @Tags			ui_api
// @Description	Endpoint de análise detalhada de remunerações dos órgãos. Recupera informações financeiras granulares para um órgão específico em um mês e ano determinados
// @Description
// @Description	Informações Fornecidos:
// @Description	- Remuneração máxima no período
// @Description	- Histograma de distribuição de rendimentos
// @Description	- Metadados do pacote de dados
// @Description
// @Description	Histograma de Distribuição de Rendimentos:
// @Description
// @Description	O histograma apresenta a quantidade de membros que receberam diferentes faixas de remunerações em um determinado mês.
// @Description	As faixas de remuneração são definidas da seguinte forma:
// @Description	- 10000: quantidade de membros que receberam até R$ 10.000,00
// @Description	- 20000: quantidade de membros que receberam entre R$ 10.000,01 e R$ 20.000,00
// @Description	- 30000: quantidade de membros que receberam entre R$ 20.000,01 e R$ 30.000,00
// @Description	- 40000: quantidade de membros que receberam entre R$ 30.000,01 e R$ 40.000,00
// @Description	- 50000: quantidade de membros que receberam entre R$ 40.000,01 e R$ 50.000,00
// @Description	- -1: quantidade de membros que receberam acima de R$ 50.000,01
// @Produce		json
// @Param			orgao	path		string				true	"Sigla do órgão. Ex.: tjal, mppb, tjmmg"
// @Param			mes		path		string				true	"Mês de referência (1-12)"
// @Param			ano		path		string				true	"Ano de referência"
// @Success		200		{object}	agencyRemuneration	"Dados de remuneração processados com sucesso"
// @Success		206		{object}	v2ProcInfoResult	"Dados coletados com informações de processamento"
// @Failure		400		{string}	string				"Erro de validação dos parâmetros de entrada"
// @Failure		500		{string}	string				"Erro interno durante processamento da consulta"
// @Router			/uiapi/v2/orgao/salario/{orgao}/{ano}/{mes} [get]
func (h handler) V2GetSalaryOfAgencyMonthYear(c echo.Context) error {
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%s inválido", c.Param("mes")))
	}
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%s inválido", c.Param("ano")))
	}
	agencyName := strings.ToLower(c.Param("orgao"))
	agencyMonthlyInfo, _, err := h.client.Db.GetOMA(month, year, agencyName)
	if err != nil {
		log.Printf("[salary agency month year] error getting data for second screen(mes:%d ano:%d, orgao:%s):%q", month, year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	if agencyMonthlyInfo.ProcInfo.String() != "" {
		var newEnv = agencyMonthlyInfo.ProcInfo.Env
		for _, omittedField := range h.envOmittedFields {
			for i, field := range newEnv {
				if strings.Contains(field, omittedField) {
					newEnv[i] = omittedField + "= ##omitida##"
					break
				}
			}
		}
		agencyMonthlyInfo.ProcInfo.Env = newEnv
		return c.JSON(http.StatusPartialContent, v2ProcInfoResult{
			ProcInfo: &procInfo{
				Stdin:  agencyMonthlyInfo.ProcInfo.Stdin,
				Stdout: agencyMonthlyInfo.ProcInfo.Stdout,
				Stderr: agencyMonthlyInfo.ProcInfo.Stderr,
				Env:    agencyMonthlyInfo.ProcInfo.Env,
				Cmd:    agencyMonthlyInfo.ProcInfo.Cmd,
				CmdDir: agencyMonthlyInfo.ProcInfo.CmdDir,
				Status: agencyMonthlyInfo.ProcInfo.Status,
			},
			Timestamp: &timestamp{
				Seconds: agencyMonthlyInfo.CrawlingTimestamp.GetSeconds(),
				Nanos:   agencyMonthlyInfo.CrawlingTimestamp.GetNanos(),
			},
		})
	}
	return c.JSON(http.StatusOK, agencyRemuneration{
		MaxRemuneration: agencyMonthlyInfo.Summary.Remunerations.Max,
		Histogram:       agencyMonthlyInfo.Summary.IncomeHistogram,
		Package: &backup{
			URL:  agencyMonthlyInfo.Package.URL,
			Hash: agencyMonthlyInfo.Package.Hash,
			Size: agencyMonthlyInfo.Package.Size,
		},
	})
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GetTotalsOfAgencyYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	aID := c.Param("orgao")
	agenciesMonthlyInfo, err := h.client.Db.GetMonthlyInfo([]strModels.Agency{{ID: aID}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, aID))
	}
	var monthTotalsOfYear []monthTotals
	agency, err := h.client.Db.GetAgency(aID)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(estado:%s):%q", aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", aID))
	}
	host := c.Request().Host
	agency.URL = fmt.Sprintf("%s/v1/orgao/%s", host, agency.ID)
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[aID] {
		if agencyMonthlyInfo.Summary != nil && agencyMonthlyInfo.Summary.BaseRemuneration.Total+agencyMonthlyInfo.Summary.OtherRemunerations.Total > 0 {
			monthTotals := monthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:   agencyMonthlyInfo.Summary.BaseRemuneration.Total,
				OtherRemunerations: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
				CrawlingTimestamp:  agencyMonthlyInfo.CrawlingTimestamp,
				TotalMembers:       agencyMonthlyInfo.Summary.Count,
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)

			// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
		} else if agencyMonthlyInfo.ProcInfo.String() != "" && agencyMonthlyInfo.ProcInfo.Status != 4 {
			monthTotals := monthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:   0,
				OtherRemunerations: 0,
				CrawlingTimestamp:  agencyMonthlyInfo.CrawlingTimestamp,
				Error:              &procError{Stdout: agencyMonthlyInfo.ProcInfo.Stdout, Stderr: agencyMonthlyInfo.ProcInfo.Stderr},
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
		}
	}
	sort.Slice(monthTotalsOfYear, func(i, j int) bool {
		return monthTotalsOfYear[i].Month < monthTotalsOfYear[j].Month
	})
	destKey := fmt.Sprintf("%s/datapackage/%s-%d.zip", aID, aID, year)
	bkp, _ := h.client.Cloud.GetFile(destKey)
	var pkg *strModels.Package
	if bkp != nil {
		pkg = &strModels.Package{
			AgencyID: &aID,
			Year:     &year,
			Package:  *bkp,
		}
	}

	agencyTotalsYear := agencyTotalsYear{Year: year, Agency: agency, MonthTotals: monthTotalsOfYear, AgencyFullName: agency.Name, SummaryPackage: pkg}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

// @ID				GetTotalsOfAgencyYear
// @Tags			ui_api
// @Description	Recupera dados financeiros detalhados para um órgão em um ano específico
// @Description
// @Description	Dados Financeiros Mensais:
// @Description	- Remuneração base
// @Description	- Outras remunerações e benefícios
// @Description	- Descontos
// @Description	- Contagem de membros
// @Description
// @Description	Métricas Adicionais:
// @Description	- Médias per capita
// @Description	- Detalhamento de rubricas (auxílios, férias, gratificações)
// @Description	- Informações gerais sobre o órgão pesquisado
// @Description	- Informações sobre o pacote de dados (URL para download, hash, tamanho)
// @Produce		json
// @Param			orgao	path		string				true	"Identificador do órgão público"	example:"tjal"
// @Param			ano		path		int					true	"Ano de referência para a consulta"	example:"2022"
// @Success		200		{object}	v2AgencyTotalsYear	"Dados financeiros completos do órgão no ano especificado"
// @Failure		400		{string}	string				"Erro de validação: parâmetros de órgão ou ano inválidos"
// @Failure		500		{string}	string				"Erro interno durante processamento da consulta"
// @Router			/uiapi/v2/orgao/totais/{orgao}/{ano} [get]
func (h handler) V2GetTotalsOfAgencyYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%s inválido", c.Param("ano")))
	}
	aID := c.Param("orgao")
	agenciesMonthlyInfo, err := h.client.Db.GetMonthlyInfo([]strModels.Agency{{ID: aID}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, aID))
	}
	var monthTotalsOfYear []v2MonthTotals
	strAgency, err := h.client.Db.GetAgency(aID)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(estado:%s):%q", aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", aID))
	}
	strAveragePerCapita, err := h.client.Db.GetAveragePerCapita(aID, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting average per capita (estado:%s):%q", aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", aID))
	}
	host := c.Request().Host
	strAgency.URL = fmt.Sprintf("%s/v2/orgao/%s", host, strAgency.ID)
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[aID] {
		if agencyMonthlyInfo.Summary != nil && agencyMonthlyInfo.Summary.BaseRemuneration.Total+agencyMonthlyInfo.Summary.OtherRemunerations.Total > 0 {
			monthTotals := v2MonthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:            agencyMonthlyInfo.Summary.BaseRemuneration.Total,
				OtherRemunerations:          agencyMonthlyInfo.Summary.OtherRemunerations.Total,
				Remunerations:               agencyMonthlyInfo.Summary.Remunerations.Total,
				Discounts:                   agencyMonthlyInfo.Summary.Discounts.Total,
				BaseRemunerationPerCapita:   agencyMonthlyInfo.Summary.BaseRemuneration.Average,
				OtherRemunerationsPerCapita: agencyMonthlyInfo.Summary.OtherRemunerations.Average,
				RemunerationsPerCapita:      agencyMonthlyInfo.Summary.Remunerations.Average,
				DiscountsPerCapita:          agencyMonthlyInfo.Summary.Discounts.Average,
				CrawlingTimestamp: timestamp{
					Seconds: agencyMonthlyInfo.CrawlingTimestamp.GetSeconds(),
					Nanos:   agencyMonthlyInfo.CrawlingTimestamp.GetNanos(),
				},
				MemberCount: agencyMonthlyInfo.Summary.Count,
				ItemSummary: itemSummary(agencyMonthlyInfo.Summary.ItemSummary),
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)

			// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
		} else if agencyMonthlyInfo.ProcInfo.String() != "" && agencyMonthlyInfo.ProcInfo.Status != 4 {
			monthTotals := v2MonthTotals{Month: agencyMonthlyInfo.Month,
				CrawlingTimestamp: timestamp{
					Seconds: agencyMonthlyInfo.CrawlingTimestamp.GetSeconds(),
					Nanos:   agencyMonthlyInfo.CrawlingTimestamp.GetNanos(),
				},
				Error: &procError{Stdout: agencyMonthlyInfo.ProcInfo.Stdout, Stderr: agencyMonthlyInfo.ProcInfo.Stderr},
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
		}
	}
	sort.Slice(monthTotalsOfYear, func(i, j int) bool {
		return monthTotalsOfYear[i].Month < monthTotalsOfYear[j].Month
	})
	destKey := fmt.Sprintf("%s/datapackage/%s-%d.zip", aID, aID, year)
	bkp, _ := h.client.Cloud.GetFile(destKey)
	var pkg *backup
	if bkp != nil {
		pkg = &backup{
			URL:  bkp.URL,
			Hash: bkp.Hash,
			Size: bkp.Size,
		}
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
	agencyTotalsYear := v2AgencyTotalsYear{
		Year: year,
		Agency: &agency{
			ID:            strAgency.ID,
			Name:          strAgency.Name,
			Type:          strAgency.Type,
			Entity:        strAgency.Entity,
			UF:            strAgency.UF,
			URL:           strAgency.URL,
			Collecting:    collect,
			TwitterHandle: strAgency.TwitterHandle,
			OmbudsmanURL:  strAgency.OmbudsmanURL,
			HasData:       hasData,
		},
		MonthTotals:    monthTotalsOfYear,
		SummaryPackage: pkg,
		AveragePerCapita: &perCapitaData{
			BaseRemuneration:   strAveragePerCapita.BaseRemuneration,
			OtherRemunerations: strAveragePerCapita.OtherRemunerations,
			Discounts:          strAveragePerCapita.Discounts,
			Remunerations:      strAveragePerCapita.Remunerations,
		},
	}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GetBasicInfoOfType(c echo.Context) error {
	yearOfConsult := time.Now().Year()
	groupName := c.Param("grupo")
	var agencies []strModels.Agency
	var err error
	var estadual bool
	var exists bool
	jurisdicao := map[string]string{
		"justica-eleitoral":    "Eleitoral",
		"ministerios-publicos": "Ministério",
		"justica-estadual":     "Estadual",
		"justica-do-trabalho":  "Trabalho",
		"justica-federal":      "Federal",
		"justica-militar":      "Militar",
		"justica-superior":     "Superior",
		"conselhos-de-justica": "Conselho",
	}

	// Adaptando as URLs do site com o banco de dados
	// Primeiro consultamos entre as chaves do mapa.
	if jurisdicao[groupName] != "" {
		groupName = jurisdicao[groupName]
	} else {
		// Caso não encontremos entre as chaves, verificamos entre os valores do mapa.
		// Isso pois, até a consolidação ser finalizada, o front consulta a api com /Eleitoral, /Trabalho, etc.
		for _, value := range jurisdicao {
			if groupName == value {
				exists = true
				break
			}
		}
		// Se a jurisdição não existir no mapa, verificamos se trata-se de um estado
		if !exists {
			values := map[string]struct{}{"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {}, "GO": {}, "MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {}, "PE": {}, "PI": {}, "RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {}, "SP": {}, "SE": {}, "TO": {}}
			if _, estadual = values[groupName]; estadual {
				exists = true
			}
		}
		// Se o parâmetro dado não for encontrado de forma alguma, retornamos um NOT FOUND (404)
		if !exists {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("Grupo não encontrado: %s.", groupName))
		}
	}

	if estadual {
		agencies, err = h.client.GetStateAgencies(groupName)
	} else {
		agencies, err = h.client.GetOPJ(groupName)
	}
	if err != nil {
		// That happens when there is no information on that year.
		log.Printf("[basic info type] first error getting data for first screen(ano:%d, grupo:%s). Going to try again with last year:%q", yearOfConsult, groupName, err)
		yearOfConsult = yearOfConsult - 1

		if estadual {
			agencies, err = h.client.GetStateAgencies(groupName)
		} else {
			agencies, err = h.client.GetOPJ(groupName)
		}
		if err != nil {
			log.Printf("[basic info type] error getting data for first screen(ano:%d, grupo:%s):%q", yearOfConsult, groupName, err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetros ano=%d ou grupo=%s são inválidos", yearOfConsult, groupName))
		}
	}
	var agenciesBasic []agencyBasic
	for k := range agencies {
		agenciesBasic = append(agenciesBasic, agencyBasic{Name: agencies[k].ID, FullName: agencies[k].Name, AgencyCategory: agencies[k].Entity})
	}
	state := state{Name: c.Param("grupo"), ShortName: "", FlagURL: "", Agency: agenciesBasic}
	return c.JSON(http.StatusOK, state)
}

// @ID				GetBasicInfoOfType
// @Tags			ui_api
// @Description	Busca informações de id (sigla), nome e entidade (jurisdiçãso) de um determinado grupo de órgãos. Ao Selecionar um grupo de órgãos por estado (Ex.: RJ, SP, etc.), retorna as informações dos tribunais de justiça desse estado (entidade=Tribunal).
// @Produce		json
// @Param			grupo						path		string	false	"Grupo de órgãos"	Enums(justica-eleitoral, ministerios-publicos, justica-estadual, justica-do-trabalho, justica-federal, justica-militar, justica-superior, conselhos-de-justica, AC, AL, AP, AM, BA, CE, DF, ES, GO, MA, MT, MS, MG, PA, PB, PR, PE, PI, RJ, RN, RS, RO, RR, SC, SP, SE, TO)
// @Success		200							{object}	state	"Órgãos do grupo"
// @Failure		400							{object}	string	"Parâmetro inválido"
// @Failure		404							{object}	string	"Grupo não encontrado"
// @Router			/uiapi/v2/orgao/{grupo} 	[get]
func (h handler) V2GetBasicInfoOfType(c echo.Context) error {
	groupName := strings.ToLower(c.Param("grupo"))
	var strAgencies []strModels.Agency
	var err error
	var estadual bool
	var exists bool
	jurisdicao := map[string]string{
		"justica-eleitoral":    "Eleitoral",
		"ministerios-publicos": "Ministério",
		"justica-estadual":     "Estadual",
		"justica-do-trabalho":  "Trabalho",
		"justica-federal":      "Federal",
		"justica-militar":      "Militar",
		"justica-superior":     "Superior",
		"conselhos-de-justica": "Conselho",
	}

	// Adaptando as URLs do site com o banco de dados
	// Primeiro consultamos entre as chaves do mapa.
	if jurisdicao[groupName] != "" {
		groupName = jurisdicao[groupName]
	} else {
		// Caso não encontremos entre as chaves, verificamos entre os valores do mapa.
		// Isso pois, até a consolidação ser finalizada, o front consulta a api com /Eleitoral, /Trabalho, etc.
		for _, value := range jurisdicao {
			if groupName == value {
				exists = true
				break
			}
		}
		// Se a jurisdição não existir no mapa, verificamos se trata-se de um estado
		if !exists {
			values := map[string]struct{}{"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {}, "GO": {}, "MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {}, "PE": {}, "PI": {}, "RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {}, "SP": {}, "SE": {}, "TO": {}}
			if _, estadual = values[strings.ToUpper(groupName)]; estadual {
				exists = true
			}
		}
		// Se o parâmetro dado não for encontrado de forma alguma, retornamos um NOT FOUND (404)
		if !exists {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("Grupo não encontrado: '%s'", c.Param("grupo")))
		}
	}

	if estadual {
		strAgencies, err = h.client.Db.GetStateAgencies(strings.ToUpper(groupName))
	} else {
		strAgencies, err = h.client.Db.GetOPJ(groupName)
	}
	if err != nil {
		// That happens when there is no information on that year.
		log.Printf("[basic info type] error getting agencies by type='%s': %q", c.Param("grupo"), err)

		if estadual {
			strAgencies, err = h.client.Db.GetStateAgencies(groupName)
		} else {
			strAgencies, err = h.client.Db.GetOPJ(groupName)
		}
		if err != nil {
			log.Printf("[basic info type] error getting data by type='%s': %q", c.Param("grupo"), err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro grupo=%s inválido", c.Param("grupo")))
		}
	}
	var agencies []v2AgencyBasic
	for k := range strAgencies {
		agencies = append(agencies, v2AgencyBasic{Id: strAgencies[k].ID, Name: strAgencies[k].Name, Entity: strAgencies[k].Entity})
	}
	group := group{Name: strings.ToUpper(c.Param("grupo")), Agencies: agencies}
	return c.JSON(http.StatusOK, group)
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GetGeneralRemunerationFromYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	data, err := h.client.Db.GetGeneralMonthlyInfosFromYear(year)
	if err != nil {
		fmt.Println("Error searching for monthly info from year: %w", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	return c.JSON(http.StatusOK, data)
}

// @ID				GetGeneralRemunerationFromYear
// @Tags			ui_api
// @Description	Busca os dados, das remunerações (remuneração base/salário, outras remunerações/benefícios, descontos) e benefícios identificados (rubricas/penduricalhos) de um ano inteiro, agrupados por mês.
// @Produce		json
// @Param			ano									path		string					true	"Ano da remuneração. Ex.: 2018, 2019, 2020..."
// @Success		200									{object}	[]mensalRemuneration	"Requisição bem sucedida."
// @Failure		400									{string}	string					"Parâmetro ano inválido."
// @Failure		500									{string}	string					"Erro interno."
// @Router			/uiapi/v2/geral/remuneracao/{ano} 	[get]
func (h handler) V2GetGeneralRemunerationFromYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%s inválido", c.Param("ano")))
	}
	data, err := h.client.Db.GetGeneralMonthlyInfosFromYear(year)
	if err != nil {
		fmt.Println("Error searching for monthly info from year: %w", err)
		return c.JSON(http.StatusInternalServerError, "error buscando dados")
	}
	annualRemu := []mensalRemuneration{}
	for _, d := range data {
		annualRemu = append(annualRemu, mensalRemuneration{
			Month:              d.Month,
			Members:            d.Count,
			BaseRemuneration:   d.BaseRemuneration,
			OtherRemunerations: d.OtherRemunerations,
			Discounts:          d.Discounts,
			Remunerations:      d.Remunerations,
			ItemSummary:        itemSummary(d.ItemSummary),
		})
	}
	return c.JSON(http.StatusOK, annualRemu)
}

// TODO: Remover quando o site tiver migrado para o novo endpoint
func (h handler) GeneralSummaryHandler(c echo.Context) error {
	agencies, err := h.client.GetAgenciesCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar orgãos: %q", err))
	}
	collections, err := h.client.GetNumberOfMonthsCollected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar registros: %q", err))
	}
	fmonth, fyear, err := h.client.Db.GetFirstDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetFirstDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando primeiro registro de remuneração: %q", err))
	}
	fdate := time.Date(fyear, time.Month(fmonth), 2, 0, 0, 0, 0, time.UTC).In(h.loc)
	lmonth, lyear, err := h.client.GetLastDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetLastDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando último registro de remuneração: %q", err))
	}
	ldate := time.Date(lyear, time.Month(lmonth), 2, 0, 0, 0, 0, time.UTC).In(h.loc)
	remuValue, err := h.client.Db.GetGeneralMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetGeneralRemunerationValue: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando valor total de remuneração: %q", err))
	}
	return c.JSON(http.StatusOK, generalTotals{
		AgencyAmount:             int(agencies),
		MonthlyTotalsAmount:      int(collections),
		StartDate:                fdate,
		EndDate:                  ldate,
		RemunerationRecordsCount: int(collections),
		GeneralRemunerationValue: remuValue,
	})
}

// @ID				GetGeneralSummary
// @Tags			ui_api
// @Description	Busca e resume os dados das remunerações de todos os anos, trazendo o número de órgão participantes das tetativas de coleta (inclue os órgãos da coleta automatizada, coleta manual e órgãos não coletados, mas que armazenamos alguma informação), número total de meses coletados, data do primeiro e último mês coletado e o valor total de remuneração considerando todos os meses.
// @Produce		json
// @Success		200						{object}	generalSummary	"Requisição bem sucedida."
// @Failure		500						{string}	string			"Erro interno do servidor."
// @Router			/uiapi/v2/geral/resumo 	[get]
func (h handler) GetGeneralSummary(c echo.Context) error {
	agencies, err := h.client.Db.GetAgenciesCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar orgãos: %q", err))
	}
	collections, err := h.client.Db.GetNumberOfMonthsCollected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar registros de meses coletados: %q", err))
	}
	paychecks, err := h.client.Db.GetNumberOfPaychecksCollected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar registros de contracheques coletados: %q", err))
	}
	fmonth, fyear, err := h.client.Db.GetFirstDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetFirstDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando primeiro registro de remuneração: %q", err))
	}
	fdate := time.Date(fyear, time.Month(fmonth), 2, 0, 0, 0, 0, time.UTC).In(h.loc)
	lmonth, lyear, err := h.client.Db.GetLastDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetLastDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando último registro de remuneração: %q", err))
	}
	ldate := time.Date(lyear, time.Month(lmonth), 2, 0, 0, 0, 0, time.UTC).In(h.loc)
	remuValue, err := h.client.Db.GetGeneralMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetGeneralRemunerationValue: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando valor total de remuneração: %q", err))
	}
	return c.JSON(http.StatusOK, generalSummary{
		Agencies:                 int(agencies),
		MonthlyInfos:             int(collections),
		Paychecks:                int(paychecks),
		StartDate:                fdate,
		EndDate:                  ldate,
		GeneralRemunerationValue: remuValue,
	})
}

// @ID				SearchByUrl
// @Tags			ui_api
// @Description	Endpoint de busca avançada para remunerações de servidores públicos
// @Description
// @Description	Permite realizar pesquisas detalhadas nas remunerações de servidores públicos com múltiplos filtros flexíveis:
// @Description
// @Description	- Filtragem por anos específicos
// @Description	- Seleção de meses específicos
// @Description	- Pesquisa por órgãos públicos de diferentes esferas
// @Description	- Categorias de remuneração
// @Description
// @Description	Características principais:
// @Description	- Suporta múltiplas seleções em cada filtro
// @Description	- Permite combinações complexas de busca
// @Description	- Retorna dados consolidados de remuneração dos contracheques por membros
// @Description
// @Description	Casos de uso:
// @Description	- Análise comparativa de remunerações entre diferentes órgãos
// @Description	- Análise análise granular das remunerações por membros dos órgãos
// @Produce		json
// @Param			anos		query		string			false	"Lista de anos a serem pesquisados, separados por virgula. Exemplo: 2018,2019,2020"
// @Param			meses		query		string			false	"Lista de meses a serem pesquisados, separados por virgula. Exemplo: 1,2,3"
// @Param			orgaos		query		string			false	"Lista de órgãos a serem pesquisados, separados por virgula. Exemplo: tjal,mpal,mppb"
// @Param			categorias	query		string			false	"Categorias a serem pesquisadas. Remuneração base (salário), outras remunerações (benefícios) e descontos"	Enums(base,outras,descontos)
// @Success		200			{object}	searchResponse	"Requisição bem-sucedida com dados de remuneração"
// @Failure		400			{string}	string			"Erro de validação dos parâmetros de busca"
// @Failure		500			{string}	string			"Erro interno do servidor durante processamento da pesquisa"
// @Router			/uiapi/v2/pesquisar [get]
func (h handler) SearchByUrl(c echo.Context) error {
	//Pegando os query params
	years := c.QueryParam("anos")
	months := c.QueryParam("meses")
	agencies := c.QueryParam("orgaos")
	categories := c.QueryParam("categorias")
	types := c.QueryParam("tipos")
	//Criando os filtros a partir dos query params e validando eles
	searchParams, err := newSearchParams(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var category string
	if searchParams != nil {
		category = searchParams.Category
	}
	// Pegando os resultados da pesquisa a partir dos filtros;
	results, err := h.db.filter(h.db.remunerationQuery(searchParams), h.db.arguments(searchParams))
	if err != nil {
		log.Printf("Error querying BD (searchParams or counter):%q", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	remunerations, numRows, err := h.getSearchResults(h.searchLimit, category, results)
	if err != nil {
		log.Printf("Error getting search results: %q", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := searchResponse{
		DownloadAvailable:  numRows > 0 && numRows <= h.downloadLimit,
		NumRowsIfAvailable: numRows,
		DownloadLimit:      h.downloadLimit,
		SearchLimit:        h.searchLimit,
		Results:            remunerations, // retornando os SearchLimit primeiros elementos.
	}
	return c.JSON(http.StatusOK, response)
}

// @ID				DownloadByUrl
// @Tags			ui_api
// @Description	Baixa um arquivo csv referentes a remunerações a partir de filtros. O arquivo tem um limite de 10 mil linhas. Para cada parâmetro, é possível passar múltiplos valores separados por vírgula. As colunas do arquivo csv são:
// @Description
// @Description	- Nome do órgão
// @Description	- Mês de referência do contracheque
// @Description	- Ano de referência do contracheque
// @Description	- Matrícula do membro (identificador único do membro no órgão)
// @Description	- Nome do membro
// @Description	- Cargo que o membro exerce no órgão
// @Description	- Lotação (unidade na qual o membro do órgão desenvolve suas atividades)
// @Description	- Categoria do contracheque (base, outras remunerações ou descontos)
// @Description	- Detalhamento do contracheque (ex: subsídio, desconto, benefício, etc)
// @Description	- Valor do contracheque em reais, não corrigido pela inflação
// @Produce		json
// @Param			anos		query		string	false	"Anos a serem pesquisados, separados por virgula. Exemplo: 2018,2019,2020"
// @Param			meses		query		string	false	"Meses a serem pesquisados, separados por virgula. Exemplo: 1,2,3"
// @Param			orgaos		query		string	false	"Orgãos a serem pesquisados, separados por virgula. Exemplo: tjal,mpal,mppb"
// @Param			categorias	query		string	false	"Categorias a serem pesquisadas. Se nada for informado, todas as categorias serão baixadas"	Enums(base,outras,descontos)
// @Success		200			{file}		file	"Arquivo CSV com os dados."
// @Failure		400			{string}	string	"Erro de validação dos parâmetros."
// @Failure		500			{string}	string	"Erro interno do servidor."
// @Router			/uiapi/v2/download [get]
func (h handler) DownloadByUrl(c echo.Context) error {
	//Pegando os query params
	years := c.QueryParam("anos")
	months := c.QueryParam("meses")
	agencies := c.QueryParam("orgaos")
	categories := c.QueryParam("categorias")
	types := c.QueryParam("tipos")

	//Criando os filtros a partir dos query params e validando eles
	searchParams, err := newSearchParams(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	results, err := h.db.filter(h.db.remunerationQuery(searchParams), h.db.arguments(searchParams))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	searchResults, _, err := h.getSearchResults(h.downloadLimit, searchParams.Category, results)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=dadosjusbr-remuneracoes.csv")
	c.Response().Header().Set("Content-Type", c.Response().Header().Get("Content-Type"))
	err = gocsv.Marshal(searchResults, c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("erro tentando fazer download do csv: %q", err))
	}
	return nil
}

//go:embed readme_content.txt
var readmeContent []byte

// @ID				DownloadReadme
// @Tags			ui_api
// @Description	Recupera o arquivo README com informações sobre o conjunto de dados. Permite filtrar o README com base em parâmetros opcionais de ano, mês e órgão
// @Description
// @Description	Comportamentos:
// @Description	- Se nenhum filtro for aplicado, retorna o README original
// @Description	- Com filtro de órgão, gera um README com observações específicas sobre potenciais falhas nos dados
// @Produce		text/plain
// @Param			ano		query		string	false	"Ano para filtragem dos dados"							default(2024)
// @Param			mes		query		string	false	"Mês para filtragem dos dados"							default(12)
// @Param			orgao	query		string	false	"Sigla do órgão para filtragem. Ex.: tjal, mppb, mpdft"	default(tjrr)
// @Success		200		{string}	string	"README.txt com conteúdo detalhado"
// @Failure		400		{string}	string	"Erro de validação de parâmetros (ano/mês inválidos)"
// @Failure		500		{string}	string	"Erro interno ao processar o README"
// @Router			/uiapi/v2/readme [get]
func (h handler) DownloadReadme(c echo.Context) error {
	originalContent := readmeContent
	year := c.QueryParam("ano")
	month := c.QueryParam("mes")
	agency := c.QueryParam("orgao")
	yearInt := 0
	monthInt := 0
	var err error

	if agency != "" {
		// Verificamos se ano e mês foram informados e se são válidos (convertemos para inteiro)
		if year != "" {
			yearInt, err = strconv.Atoi(year)
			if err != nil {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ANO inválido: %s.", year))
			}
			if month != "" {
				monthInt, err = strconv.Atoi(month)
				if err != nil {
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro MÊS inválido: %s.", month))
				}
			}
		}
		var results []*string
		if results, err = h.client.Db.GetNotices(agency, yearInt, monthInt); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("erro coletando os avisos: %q", err))
		}
		// Filtrar nulos e converter para []string
		var cleanResults []string
		for _, s := range results {
			if s != nil {
				cleanResults = append(cleanResults, *s)
			}
		}

		var newLines string
		if len(cleanResults) != 0 {
			newLines = strings.Join(cleanResults, "\n")
		} else {
			newLines = "Não identificamos potenciais falhas na origem destes dados."
		}

		// Adicionar as novas linhas ao conteúdo original
		updatedContent := "\n**Observações sobre este conjunto de dados**:\n\n" +
			newLines +
			" Em sua análise, esteja atento a possíveis valores estranhos.\n\n"

		// Organizando para que seja o 2º tópico
		newContent := append(originalContent[:1733], append([]byte(updatedContent), originalContent[1734:]...)...)
		originalContent = newContent
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=README.txt")
	c.Response().Header().Set("Content-Type", "text/plain")
	c.Response().WriteHeader(http.StatusOK)
	_, err = c.Response().Write(originalContent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("erro tentando fazer download do readme: %q", err))
	}
	return nil
}

// @ID				GetAnnualSummary
// @Tags			ui_api
// @Description	Retorna os dados de remuneração de todos os anos disponíveis para um órgão específico, incluindo:
// @Description	- Remuneração base/salário, outras remunerações/benefícios, descontos e remuneração líquida (salário+benefícios-descontos). Dados brutos, agrupados por mês e per capita
// @Description	- Quantidade de meses com dados no determinado ano
// @Description	- Quantidade média de membros do órgão naquele ano
// @Description	- Resumo dos benefícios identificados (rubricas/penduricalhos) e seus respectivos valores no ano
// @Description	- Informações do pacote de dados, URL do pacote de dados para download, seu hash e tamanho do pacote de dados (em bytes)
// @Produce		json
// @Param			orgao	path		string			true	"Nome do orgão"
// @Success		200		{object}	[]annualSummary	"Requisição bem sucedida."
// @Failure		400		{string}	string			"Parâmetro orgao inválido"
// @Failure		500		{string}	string			"Algo deu errado ao tentar coletar os dados anuais do orgao"
// @Router			/uiapi/v1/orgao/resumo/{orgao} [get]
func (h handler) GetAnnualSummary(c echo.Context) error {
	agencyName := c.Param("orgao")
	strAgency, err := h.client.Db.GetAgency(agencyName)
	if err != nil {
		log.Printf("error getting agency '%s' :%q", agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", agencyName))
	}
	host := c.Request().Host
	strAgency.URL = fmt.Sprintf("%s/v2/orgao/%s", host, strAgency.ID)
	summaries, err := h.client.GetAnnualSummary(agencyName)
	if err != nil {
		log.Printf("error getting annual data of '%s' :%q", agencyName, err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Algo deu errado ao tentar coletar os dados anuais do orgao=%s", agencyName))
	}
	var annualData []annualSummaryData
	for _, s := range summaries {
		baseRemPerMonth := s.BaseRemuneration / float64(s.NumMonthsWithData)
		baseRemPerCapita := s.BaseRemunerationPerCapita
		otherRemPerMonth := s.OtherRemunerations / float64(s.NumMonthsWithData)
		otherRemPerCapita := s.OtherRemunerationsPerCapita
		remPerMonth := s.Remunerations / float64(s.NumMonthsWithData)
		remPerCapita := s.RemunerationsPerCapita
		discountsRemPerMonth := s.Discounts / float64(s.NumMonthsWithData)
		discountsRemPerCapita := s.DiscountsPerCapita

		annualData = append(annualData, annualSummaryData{
			Year:                        s.Year,
			AverageMemberCount:          s.AverageCount,
			BaseRemuneration:            s.BaseRemuneration,
			BaseRemunerationPerMonth:    baseRemPerMonth,
			BaseRemunerationPerCapita:   baseRemPerCapita,
			OtherRemunerations:          s.OtherRemunerations,
			OtherRemunerationsPerMonth:  otherRemPerMonth,
			OtherRemunerationsPerCapita: otherRemPerCapita,
			Discounts:                   s.Discounts,
			DiscountsPerMonth:           discountsRemPerMonth,
			DiscountsPerCapita:          discountsRemPerCapita,
			Remunerations:               s.Remunerations,
			RemunerationsPerMonth:       remPerMonth,
			RemunerationsPerCapita:      remPerCapita,
			NumMonthsWithData:           s.NumMonthsWithData,
			Package: &backup{
				URL:  s.Package.URL,
				Hash: s.Package.Hash,
				Size: s.Package.Size,
			},
			ItemSummary:  itemSummary(s.ItemSummary),
			Inconsistent: s.Inconsistent,
		})
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
	annualSum := annualSummary{
		Agency: &agency{
			ID:            strAgency.ID,
			Name:          strAgency.Name,
			URL:           strAgency.URL,
			Type:          strAgency.Type,
			Entity:        strAgency.Entity,
			UF:            strAgency.UF,
			Collecting:    collect,
			TwitterHandle: strAgency.TwitterHandle,
			OmbudsmanURL:  strAgency.OmbudsmanURL,
			HasData:       hasData,
		},
		Data: annualData,
	}
	return c.JSON(http.StatusOK, annualSum)
}

func (h handler) getSearchResults(limit int, category string, results []searchDetails) ([]searchResult, int, error) {
	searchResults := []searchResult{}
	numRows := 0
	if len(results) == 0 {
		return searchResults, numRows, nil
	} else {
		// A razão para essa ordenação é que quando o usuário escolhe diversos órgãos
		// provavelmente ele prefere ver dados de todos eles. Dessa forma, aumentamos
		// as chances do preview limitado retornar dados de diversos órgãos.
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].Ano < results[j].Ano || results[i].Mes < results[j].Mes
		})
		searchResults, numRows, err := h.sess.getRemunerationsFromS3(limit, h.downloadLimit, category, h.s3Bucket, results)
		if err != nil {
			return nil, numRows, fmt.Errorf("failed to get remunerations from s3 %q", err)
		}
		return searchResults, numRows, nil
	}
}

// @ID				GetAveragePerAgency
// @Tags			ui_api
// @Description	Busca médias (remuneração base, outras remunerações, descontos e remuneração total) de cada órgão em um ano especificado.
// @Produce		json
// @Param			ano	path		int					true	"Ano para filtrar os dados"
// @Success		200	{array}		averagePerAgency	"Lista de dados de médias dos órgãos"
// @Failure		400	{string}	string				"Parâmetro ANO inválido"
// @Failure		500	{string}	string				"Erro ao buscar dados"
// @Router			/uiapi/v2/orgao/media/{ano} [get]
func (h handler) GetAveragePerAgency(c echo.Context) error {
	year := c.Param("ano")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ANO inválido: %s.", year))
	}

	// Busca os as médias por membro do banco de dados
	data, err := h.client.Db.GetAveragePerAgency(yearInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao buscar dados: %q", err))
	}

	var avgPerAgency []averagePerAgency
	for _, d := range data {
		avgPerAgency = append(avgPerAgency, averagePerAgency{
			ID: d.AgencyID,
			AveragePerMember: &perCapitaData{
				BaseRemuneration:   d.BaseRemuneration,
				OtherRemunerations: d.OtherRemunerations,
				Discounts:          d.Discounts,
				Remunerations:      d.Remunerations},
		})
	}

	return c.JSON(http.StatusOK, avgPerAgency)
}

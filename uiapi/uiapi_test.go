package uiapi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dadosjusbr/proto/coleta"
	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var loc *time.Location
var hand handler

func TestMain(m *testing.M) {
	l, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err.Error())
	}
	loc = l
	exitValue := m.Run()
	os.Exit(exitValue)
}

func TestGetSummaryOfAgency(t *testing.T) {
	tests := getSummaryOfAgency{}
	t.Run("Test GetSummaryOfAgency when data exists", tests.testWhenDataExists)
	t.Run("Test GetSummaryOfAgency when data does not exist", tests.testWhenDataDoesNotExist)
	t.Run("Test GetSummaryOfAgency when year is invalid", tests.testWhenYearIsInvalid)
	t.Run("Test GetSummaryOfAgency when month is invalid", tests.testWhenMonthIsInvalid)
}

type getSummaryOfAgency struct{}

func (g getSummaryOfAgency) testWhenDataExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agmi := agencyMonthlyInfos()
	agency := &models.Agency{
		ID:            "tjal",
		Name:          "Tribunal de Justiça do Estado de Alagoas",
		Type:          "Estadual",
		Entity:        "Tribunal",
		UF:            "AL",
		TwitterHandle: "TJALagoas",
		OmbudsmanURL:  "TJALagoas.com.br",
	}
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetOMA(1, 2020, "tjal").Return(&agmi[0], agency, nil)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/resumo/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSummaryOfAgency(ctx)

	expectedCode := http.StatusOK
	expectedJson := `
		{
			"orgao": "Tribunal de Justiça do Estado de Alagoas",
			"remuneracao_base": 7.099024400000013e+06,
			"max_remuneracao_base": 35462.22,
			"outras_remuneracoes": 1.9515865600000022e+06,
			"max_outras_remuneracoes": 45200.05,
			"timestamp": {
				"seconds": 1,
				"nanos": 1
			},
			"total_membros": 214,
			"total_remuneracao": 9.050610960000016e+06,
			"tem_proximo": true,
			"tem_anterior": true
		}
	`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, recorder.Body.String())
}

func (g getSummaryOfAgency) testWhenDataDoesNotExist(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetOMA(1, 2020, "tjal").Return(nil, nil, fmt.Errorf("there is no data with this parameters"))

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/resumo/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSummaryOfAgency(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro ano=2020, mês=1 ou nome do orgão=tjal são inválidos"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSummaryOfAgency) testWhenYearIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/resumo/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020a", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSummaryOfAgency(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro ano=2020a inválido"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSummaryOfAgency) testWhenMonthIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/resumo/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1a")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSummaryOfAgency(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro mês=1a inválido"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func TestGetSalaryOfAgencyMonthYear(t *testing.T) {
	g := getSalaryOfAgencyMonthYear{}
	t.Run("Test when data exists", g.testWhenDataExists)
	t.Run("Test when data does not exist", g.testWhenDataDoesNotExist)
	t.Run("Test when year is invalid", g.testWhenYearIsInvalid)
	t.Run("Test when month is invalid", g.testWhenMonthIsInvalid)
	t.Run("Test when procinfo is not null", g.testWhenProcInfoIsNotNull)
}

type getSalaryOfAgencyMonthYear struct{}

func (g getSalaryOfAgencyMonthYear) testWhenDataExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agmi := agencyMonthlyInfos()
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetOMA(1, 2020, "tjal").Return(&agmi[0], nil, nil)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/salario/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSalaryOfAgencyMonthYear(ctx)

	expectedCode := http.StatusOK
	expectedJson := `
		{
			"max_salario": 35462.22,
			"histograma": {
				"-1": 0,
				"10000": 1,
				"20000": 0,
				"30000": 3,
				"40000": 210,
				"50000": 0
			},
			"package": {
				"url": "https://dadosjusbr.org/download/tjal/datapackage/tjal-2020-1.zip",
				"hash": "4d7ca8986101673aea060ac1d8e5a529",
				"size": 30195
			}
		}
	`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSalaryOfAgencyMonthYear) testWhenDataDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetOMA(1, 2020, "tjal").Return(nil, nil, fmt.Errorf("there is no data with this parameters"))

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/salario/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSalaryOfAgencyMonthYear(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro ano=2020, mês=1 ou nome do orgão=tjal são inválidos"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSalaryOfAgencyMonthYear) testWhenYearIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/salario/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020a", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSalaryOfAgencyMonthYear(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro ano=2020a inválido"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSalaryOfAgencyMonthYear) testWhenMonthIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/salario/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1a")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSalaryOfAgencyMonthYear(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro mês=1a inválido"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recorder.Body.String(), "\n"))
}

func (g getSalaryOfAgencyMonthYear) testWhenProcInfoIsNotNull(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agmi := agencyMonthlyInfos()
	agmi[0].ProcInfo = &coleta.ProcInfo{
		Stderr: "stderr",
		Cmd:    "docker run ...",
		CmdDir: "/tmp",
		Status: 4,
		Env:    []string{"VAR1=1", "VAR2=2"},
	}
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetOMA(1, 2020, "tjal").Return(&agmi[0], nil, nil)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/salario/:orgao/:ano/:mes",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano", "mes")
	ctx.SetParamValues("tjal", "2020", "1")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetSalaryOfAgencyMonthYear(ctx)

	expectedCode := http.StatusPartialContent
	expectedJson := `
		{
			"proc_info": {
				"stderr": "stderr",
				"cmd": "docker run ...",
				"cmd_dir": "/tmp",
				"status": 4,
				"env": [
					"VAR1=1",
					"VAR2=2"
				]
			},
			"timestamp": {
				"seconds": 1,
				"nanos": 1
			}
		}
	`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, recorder.Body.String())
}

func TestGetTotalsOfAgencyYear(t *testing.T) {
	tests := getTotalsOfAgencyYear{}
	t.Run("test when data exists", tests.testWhenDataExists)
	t.Run("test when monthly info does not exist", getTotalsOfAgencyYear{}.testWhenMonthlyInfoDoesNotExist)
	t.Run("test when year is invalid", tests.testWhenYearIsInvalid)
}

type getTotalsOfAgencyYear struct{}

func (g getTotalsOfAgencyYear) testWhenDataExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agmi := agencyMonthlyInfos()
	agencies := []models.Agency{
		{
			ID:            "tjal",
			Name:          "Tribunal de Justiça do Estado de Alagoas",
			Type:          "Estadual",
			Entity:        "Tribunal",
			UF:            "AL",
			TwitterHandle: "tjaloficial",
			OmbudsmanURL:  "http://www.tjal.jus.br/ombudsman",
		},
	}
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAgency("tjal").Return(&agencies[0], nil).Times(1)
	dbMock.EXPECT().GetMonthlyInfo([]models.Agency{{ID: "tjal"}}, 2020).Return(map[string][]models.AgencyMonthlyInfo{"tjal": agmi}, nil).Times(1)
	fsMock.EXPECT().GetFile("tjal/datapackage/tjal-2020.zip").Return(agmi[0].Package, nil)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/uiapi/v2/orgao/totais/:orgao/:ano",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano")
	ctx.SetParamValues("tjal", "2020")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetTotalsOfAgencyYear(ctx)

	expectedCode := http.StatusOK
	expectedJson := `
		{
			"ano": 2020,
			"orgao": {
				"id_orgao": "tjal",
				"nome": "Tribunal de Justiça do Estado de Alagoas",
				"jurisdicao": "Estadual",
				"entidade": "Tribunal",
				"uf": "AL",
				"twitter_handle": "tjaloficial",
				"ouvidoria": "http://www.tjal.jus.br/ombudsman",
				"url": "example.com/v1/orgao/tjal"
			},
			"meses": [
				{
					"mes": 1,
					"outras_remuneracoes":1.9515865600000022e+06,
					"remuneracao_base":7.099024400000013e+06,
					"timestamp": {
						"seconds": 1,
						"nanos": 1
					},
					"total_membros": 214
				}
			],
			"package": {
				"url": "https://dadosjusbr.org/download/tjal/datapackage/tjal-2020-1.zip",
				"hash": "4d7ca8986101673aea060ac1d8e5a529",
				"size": 30195
			}
		}
	`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, recorder.Body.String())
}

func (g getTotalsOfAgencyYear) testWhenMonthlyInfoDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agency := models.Agency{
		ID:            "tjal",
		Name:          "Tribunal de Justiça do Estado de Alagoas",
		Type:          "Estadual",
		Entity:        "Tribunal",
		UF:            "AL",
		TwitterHandle: "tjaloficial",
		OmbudsmanURL:  "http://www.tjal.jus.br/ombudsman",
	}
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAgency("tjal").Return(&agency, nil).Times(1)
	dbMock.EXPECT().GetMonthlyInfo([]models.Agency{{ID: "tjal"}}, 2020).Return(nil, nil).Times(1)
	fsMock.EXPECT().GetFile("tjal/datapackage/tjal-2020.zip").Return(nil, nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/uiapi/v2/orgao/totais/:orgao/:ano",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano")
	ctx.SetParamValues("tjal", "2020")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetTotalsOfAgencyYear(ctx)

	expectedCode := http.StatusOK
	expectedJson := `
		{
				"ano": 2020,
				"orgao": {
					"id_orgao": "tjal",
					"nome": "Tribunal de Justiça do Estado de Alagoas",
					"jurisdicao": "Estadual",
					"entidade": "Tribunal",
					"uf": "AL",
					"twitter_handle": "tjaloficial",
					"ouvidoria": "http://www.tjal.jus.br/ombudsman",
					"url": "example.com/v1/orgao/tjal"
				}
		}
	`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, recorder.Body.String())
}

func (g getTotalsOfAgencyYear) testWhenYearIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/uiapi/v2/orgao/totais/:orgao/:ano",
		nil,
	)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.SetParamNames("orgao", "ano")
	ctx.SetParamValues("tjal", "2020a")

	client, _ := storage.NewClient(dbMock, fsMock)
	handler, err := NewHandler(client, nil, nil, "us-east-1", "dadosjusbr_public", loc, []string{}, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	handler.V2GetTotalsOfAgencyYear(ctx)

	expectedCode := http.StatusBadRequest
	expectedJson := `"Parâmetro ano=2020a inválido"`

	assert.Equal(t, expectedCode, recorder.Code)
	assert.JSONEq(t, expectedJson, recorder.Body.String())
}

func agencyMonthlyInfos() []models.AgencyMonthlyInfo {
	return []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjal",
			Month:    1,
			Year:     2020,
			Summary: &models.Summary{
				Count: 214,
				BaseRemuneration: models.DataSummary{
					Max:     35462.22,
					Min:     7473.09,
					Average: 33173.01121495333,
					Total:   7099024.400000013,
				},
				OtherRemunerations: models.DataSummary{
					Max:     45200.05,
					Average: 9119.563364485992,
					Total:   1951586.5600000022,
				},
				IncomeHistogram: map[int]int{
					-1:    0,
					10000: 1,
					20000: 0,
					30000: 3,
					40000: 210,
					50000: 0,
				},
			},
			CrawlerVersion:    "unspecified",
			CrawlerRepo:       "https://github.com/dadosjusbr/coletor-cnj",
			CrawlingTimestamp: timestamppb.New(time.Unix(1, 1)),
			Package: &models.Backup{
				URL:  "https://dadosjusbr.org/download/tjal/datapackage/tjal-2020-1.zip",
				Hash: "4d7ca8986101673aea060ac1d8e5a529",
				Size: 30195,
			},
			Meta: &models.Meta{
				OpenFormat:       false,
				Access:           "NECESSITA_SIMULACAO_USUARIO",
				Extension:        "XLS",
				StrictlyTabular:  true,
				ConsistentFormat: true,
				HaveEnrollment:   false,
				ThereIsACapacity: false,
				HasPosition:      false,
				BaseRevenue:      "DETALHADO",
				OtherRecipes:     "DETALHADO",
				Expenditure:      "DETALHADO",
			},
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
			ProcInfo: &coleta.ProcInfo{},
		},
	}
}

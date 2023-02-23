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
					1: 10,
					2: 20,
					3: 30,
					4: 40,
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
		},
	}
}

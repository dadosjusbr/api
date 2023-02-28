package papi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAgencyById(t *testing.T) {
	tests := getAgencyById{}
	t.Run("Test GetAgencyById when agency exists", tests.testGetAgencyByIdWhenAgencyExists)
	t.Run("Test GetAgencyById when agency does not exist", tests.testGetAgencyByIdWhenAgencyDoesNotExist)
}

type getAgencyById struct{}

func (g getAgencyById) testGetAgencyByIdWhenAgencyExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agency := &models.Agency{
		ID:            "tjal",
		Name:          "Tribunal de Justiça do Estado de Alagoas",
		Type:          "Estadual",
		Entity:        "Tribunal",
		UF:            "AL",
		TwitterHandle: "TJALagoas",
		OmbudsmanURL:  "TJALagoas.com.br",
	}
	agencyId := "tjal"
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAgency(agencyId).Return(agency, nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/:orgao",
		nil,
	)
	recoder := httptest.NewRecorder()
	ctx := e.NewContext(request, recoder)
	ctx.SetParamNames("orgao")
	ctx.SetParamValues(agencyId)

	client, _ := storage.NewClient(dbMock, fsMock)
	handler := NewHandler(client, "", "")
	handler.V2GetAgencyById(ctx)

	expectedHttpCode := 200
	expectedJson := `
		{
			"id_orgao": "tjal",
			"nome": "Tribunal de Justiça do Estado de Alagoas",
			"jurisdicao": "Estadual",
			"entidade": "Tribunal",
			"uf": "AL",
			"url": "example.com/v2/orgao/tjal",
			"twitter_handle": "TJALagoas",
			"ouvidoria": "TJALagoas.com.br"
		}
	`

	assert.Equal(t, expectedHttpCode, recoder.Code)
	assert.JSONEq(t, expectedJson, recoder.Body.String())
}

func (g getAgencyById) testGetAgencyByIdWhenAgencyDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agencyId := "tjal"
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAgency(agencyId).Return(nil, fmt.Errorf("agency not found")).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgao/:orgao",
		nil,
	)
	recoder := httptest.NewRecorder()
	ctx := e.NewContext(request, recoder)
	ctx.SetParamNames("orgao")
	ctx.SetParamValues(agencyId)

	client, _ := storage.NewClient(dbMock, fsMock)
	handler := NewHandler(client, "", "")
	handler.V2GetAgencyById(ctx)

	expectedHttpCode := 404
	expectedJson := `"Agency not found"`

	assert.Equal(t, expectedHttpCode, recoder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recoder.Body.String(), "\n"))
}

func TestGetAllAgencies(t *testing.T) {
	tests := getAllAgencies{}
	t.Run("Test GetAllAgencies when agencies exists", tests.testGetAllAgenciesWhenAgenciesExists)
	t.Run("Test GetAllAgencies when agencies does not exist", tests.testGetAllAgenciesWhenAgenciesDoesNotExist)
}

type getAllAgencies struct{}

func (g getAllAgencies) testGetAllAgenciesWhenAgenciesExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	agencies := []models.Agency{
		{
			ID:            "tjal",
			Name:          "Tribunal de Justiça do Estado de Alagoas",
			Type:          "Estadual",
			Entity:        "Tribunal",
			UF:            "AL",
			TwitterHandle: "TJALagoas",
			OmbudsmanURL:  "TJALagoas.com.br",
		},
		{
			ID:            "tjba",
			Name:          "Tribunal de Justiça do Estado da Bahia",
			Type:          "Estadual",
			Entity:        "Tribunal",
			UF:            "BA",
			TwitterHandle: "TJBA",
			OmbudsmanURL:  "TJBA.com.br",
		},
	}
	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAllAgencies().Return(agencies, nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgaos",
		nil,
	)
	recoder := httptest.NewRecorder()
	ctx := e.NewContext(request, recoder)

	client, _ := storage.NewClient(dbMock, fsMock)
	handler := NewHandler(client, "", "")
	handler.V2GetAllAgencies(ctx)

	expectedHttpCode := 200
	expectedJson := `
		[
			{
				"id_orgao": "tjal",
				"nome": "Tribunal de Justiça do Estado de Alagoas",
				"jurisdicao": "Estadual",
				"entidade": "Tribunal",
				"uf": "AL",
				"url": "example.com/v2/orgao/tjal",
				"twitter_handle": "TJALagoas",
				"ouvidoria": "TJALagoas.com.br"
			},
			{
				"id_orgao": "tjba",
				"nome": "Tribunal de Justiça do Estado da Bahia",
				"jurisdicao": "Estadual",
				"entidade": "Tribunal",
				"uf": "BA",
				"url": "example.com/v2/orgao/tjba",
				"twitter_handle": "TJBA",
				"ouvidoria": "TJBA.com.br"
			}
		]
	`
	assert.Equal(t, expectedHttpCode, recoder.Code)
	assert.JSONEq(t, expectedJson, recoder.Body.String())
}

func (g getAllAgencies) testGetAllAgenciesWhenAgenciesDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCtrl)
	fsMock := file_storage.NewMockInterface(mockCtrl)

	dbMock.EXPECT().Connect().Return(nil).Times(1)
	dbMock.EXPECT().GetAllAgencies().Return([]models.Agency{}, nil).Times(1)

	e := echo.New()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v2/orgaos",
		nil,
	)
	recoder := httptest.NewRecorder()
	ctx := e.NewContext(request, recoder)

	client, _ := storage.NewClient(dbMock, fsMock)
	handler := NewHandler(client, "", "")
	handler.V2GetAllAgencies(ctx)

	expectedHttpCode := 200
	expectedJson := "[]"

	assert.Equal(t, expectedHttpCode, recoder.Code)
	assert.Equal(t, expectedJson, strings.Trim(recoder.Body.String(), "\n"))
}

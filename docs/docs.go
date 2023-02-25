// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "DadosJusBr",
            "url": "https://dadosjusbr.org"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/uiapi/v2/orgao/resumo/{orgao}/{ano}/{mes}": {
            "get": {
                "description": "Resume os dados de remuneração mensal de um órgão.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ui_api"
                ],
                "operationId": "GetSummaryOfAgency",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do órgão. Exemplos: tjal, tjba, mppb.",
                        "name": "orgao",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Ano da remuneração. Exemplo: 2018.",
                        "name": "ano",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Mês da remuneração. Exemplo: 1.",
                        "name": "mes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Requisição bem sucedida.",
                        "schema": {
                            "$ref": "#/definitions/uiapi.v2AgencySummary"
                        }
                    },
                    "400": {
                        "description": "Parâmetro ano, mês ou nome do órgão são inválidos.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Órgão não encontrado.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/uiapi/v2/orgao/salario/{orgao}/{ano}/{mes}": {
            "get": {
                "description": "Busca dados das remunerações mensais de um órgão.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ui_api"
                ],
                "operationId": "GetSalaryOfAgencyMonthYear",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do órgão. Exemplos: tjal, tjba, mppb.",
                        "name": "orgao",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Mês da remuneração. Exemplos: 01, 02, 03...",
                        "name": "mes",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Ano da remuneração. Exemplos: 2018, 2019, 2020...",
                        "name": "ano",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Requisição bem sucedida.",
                        "schema": {
                            "$ref": "#/definitions/uiapi.agencySalary"
                        }
                    },
                    "206": {
                        "description": "Requisição bem sucedida, mas os dados do órgão não foram bem processados",
                        "schema": {
                            "$ref": "#/definitions/uiapi.v2ProcInfoResult"
                        }
                    },
                    "400": {
                        "description": "Parâmetros inválidos.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/uiapi/v2/orgao/totais/{orgao}/{ano}": {
            "get": {
                "description": "Busca os dados de remuneração de um órgão em um ano específico.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ui_api"
                ],
                "operationId": "GetTotalsOfAgencyYear",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do órgão. Exemplos: tjal, tjba, mppb.",
                        "name": "orgao",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Ano. Exemplo: 2018.",
                        "name": "ano",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Requisição bem sucedida.",
                        "schema": {
                            "$ref": "#/definitions/uiapi.v2AgencyTotalsYear"
                        }
                    },
                    "400": {
                        "description": "Parâmetro ano ou orgao inválido.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/orgao/{orgao}": {
            "get": {
                "description": "Busca um órgão específico utilizando seu ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public_api"
                ],
                "operationId": "GetAgencyById",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do órgão. Exemplos: tjal, tjba, mppb.",
                        "name": "orgao",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Requisição bem sucedida.",
                        "schema": {
                            "$ref": "#/definitions/papi.agency"
                        }
                    },
                    "404": {
                        "description": "Órgão não encontrado.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "papi.agency": {
            "type": "object",
            "properties": {
                "coletando": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/papi.collecting"
                    }
                },
                "entidade": {
                    "description": "\"J\" For Judiciário, \"M\" for Ministério Público, \"P\" for Procuradorias and \"D\" for Defensorias.",
                    "type": "string"
                },
                "id_orgao": {
                    "description": "'trt13'",
                    "type": "string"
                },
                "jurisdicao": {
                    "description": "\"R\" for Regional, \"M\" for Municipal, \"F\" for Federal, \"E\" for State.",
                    "type": "string"
                },
                "nome": {
                    "description": "'Tribunal Regional do Trabalho 13° Região'",
                    "type": "string"
                },
                "ouvidoria": {
                    "description": "Agencys's ombudsman url",
                    "type": "string"
                },
                "twitter_handle": {
                    "description": "Agency's twitter handle",
                    "type": "string"
                },
                "uf": {
                    "description": "Short code for federative unity.",
                    "type": "string"
                },
                "url": {
                    "description": "Link for state url",
                    "type": "string"
                }
            }
        },
        "papi.collecting": {
            "type": "object",
            "properties": {
                "descricao": {
                    "description": "Reasons why we didn't collect the data",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "timestamp": {
                    "description": "Day(unix) we checked the status of the data",
                    "type": "integer"
                }
            }
        },
        "uiapi.agency": {
            "type": "object",
            "properties": {
                "coletando": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/uiapi.collecting"
                    }
                },
                "entidade": {
                    "description": "\"J\" For Judiciário, \"M\" for Ministério Público, \"P\" for Procuradorias and \"D\" for Defensorias.",
                    "type": "string"
                },
                "id_orgao": {
                    "description": "'trt13'",
                    "type": "string"
                },
                "jurisdicao": {
                    "description": "\"R\" for Regional, \"M\" for Municipal, \"F\" for Federal, \"E\" for State.",
                    "type": "string"
                },
                "nome": {
                    "description": "'Tribunal Regional do Trabalho 13° Região'",
                    "type": "string"
                },
                "ouvidoria": {
                    "description": "Agencys's ombudsman url",
                    "type": "string"
                },
                "twitter_handle": {
                    "description": "Agency's twitter handle",
                    "type": "string"
                },
                "uf": {
                    "description": "Short code for federative unity.",
                    "type": "string"
                },
                "url": {
                    "description": "Link for state url",
                    "type": "string"
                }
            }
        },
        "uiapi.agencySalary": {
            "type": "object",
            "properties": {
                "histograma": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "max_salario": {
                    "type": "number"
                },
                "package": {
                    "$ref": "#/definitions/uiapi.backup"
                }
            }
        },
        "uiapi.backup": {
            "type": "object",
            "properties": {
                "hash": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "uiapi.collecting": {
            "type": "object",
            "properties": {
                "descricao": {
                    "description": "Reasons why we didn't collect the data",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "timestamp": {
                    "description": "Day(unix) we checked the status of the data",
                    "type": "integer"
                }
            }
        },
        "uiapi.procError": {
            "type": "object",
            "properties": {
                "stderr": {
                    "description": "String containing the standard error of the process.",
                    "type": "string"
                },
                "stdout": {
                    "description": "String containing the standard output of the process.",
                    "type": "string"
                }
            }
        },
        "uiapi.procInfo": {
            "type": "object",
            "properties": {
                "cmd": {
                    "type": "string"
                },
                "cmd_dir": {
                    "type": "string"
                },
                "env": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "integer"
                },
                "stderr": {
                    "type": "string"
                },
                "stdin": {
                    "type": "string"
                },
                "stdout": {
                    "type": "string"
                }
            }
        },
        "uiapi.timestamp": {
            "type": "object",
            "properties": {
                "nanos": {
                    "type": "integer"
                },
                "seconds": {
                    "type": "integer"
                }
            }
        },
        "uiapi.v2AgencySummary": {
            "type": "object",
            "properties": {
                "max_outras_remuneracoes": {
                    "type": "number"
                },
                "max_remuneracao_base": {
                    "type": "number"
                },
                "orgao": {
                    "type": "string"
                },
                "outras_remuneracoes": {
                    "type": "number"
                },
                "remuneracao_base": {
                    "type": "number"
                },
                "tem_anterior": {
                    "type": "boolean"
                },
                "tem_proximo": {
                    "type": "boolean"
                },
                "timestamp": {
                    "$ref": "#/definitions/uiapi.timestamp"
                },
                "total_membros": {
                    "type": "integer"
                },
                "total_remuneracao": {
                    "type": "number"
                }
            }
        },
        "uiapi.v2AgencyTotalsYear": {
            "type": "object",
            "properties": {
                "ano": {
                    "type": "integer"
                },
                "meses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/uiapi.v2MonthTotals"
                    }
                },
                "orgao": {
                    "$ref": "#/definitions/uiapi.agency"
                },
                "package": {
                    "$ref": "#/definitions/uiapi.backup"
                }
            }
        },
        "uiapi.v2MonthTotals": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/uiapi.procError"
                },
                "mes": {
                    "type": "integer"
                },
                "outras_remuneracoes": {
                    "type": "number"
                },
                "remuneracao_base": {
                    "type": "number"
                },
                "timestamp": {
                    "$ref": "#/definitions/uiapi.timestamp"
                },
                "total_membros": {
                    "type": "integer"
                }
            }
        },
        "uiapi.v2ProcInfoResult": {
            "type": "object",
            "properties": {
                "proc_info": {
                    "$ref": "#/definitions/uiapi.procInfo"
                },
                "timestamp": {
                    "$ref": "#/definitions/uiapi.timestamp"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API do dadosjusbr.org",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

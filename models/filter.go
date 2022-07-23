package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter struct {
	years      []string
	months     []string
	agencies   []string
	categories []string
	types      []string
}

func NewFilter(yearsQp, monthsQp, agenciesQp, categoriesQp, typesQp string) (*Filter, error) {
	var years []string
	var months []string
	var agencies []string
	var categories []string
	var types []string

	if yearsQp == "" && monthsQp == "" && agenciesQp == "" && categoriesQp == "" && typesQp == "" {
		return nil, nil
	}
	if yearsQp != "" {
		years = strings.Split(yearsQp, ",")
		for _, y := range years {
			if _, err := strconv.Atoi(y); err != nil {
				return nil, fmt.Errorf("parâmetro ano '%s' é inválido!", y)
			}
		}
	}
	if monthsQp != "" {
		months = strings.Split(monthsQp, ",")
		for _, m := range months {
			if _, err := strconv.Atoi(m); err != nil {
				return nil, fmt.Errorf("parâmetro mês '%s' é inválido!", m)
			}
		}
	}
	if agenciesQp != "" {
		agencies = strings.Split(agenciesQp, ",")
	}
	if categoriesQp != "" {
		categories = strings.Split(categoriesQp, ",")
	}
	if typesQp != "" {
		types = strings.Split(typesQp, ",")
	}

	return &Filter{
		years:      years,
		months:     months,
		agencies:   agencies,
		categories: categories,
		types:      types,
	}, nil
}

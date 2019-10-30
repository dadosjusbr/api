package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test if solution gives right answers for questions.
func TestSolution(t *testing.T) {
	data := []struct {
		desc     string
		question string
		answer   string
	}{
		{"Sequence question", "Que número completa a seguinte sequência? 92, 93, 94, 95", "96"},
		{"Even or Odd question - odd", "O número 5 é par ou ímpar?", "impar"},
		{"Even or Odd question - even", "O número 16 é par ou ímpar?", "par"},
		{"Arithmetic question - sum", "Quanto é 8 + 6?", "14"},
		{"Arithmetic question - sub", "Quanto é 17 - 6?", "11"},
		{"After month question", "Depois de outubro vem qual mês?", "novembro"},
		{"After month question - corner case", "Depois de dezembro vem qual mês?", "janeiro"},
		{"Before month question", "Antes de outubro vem qual mês?", "setembro"},
		{"Before month question - corner case", "Antes de janeiro vem qual mês?", "dezembro"},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			ans, err := solution(d.question)
			assert.NoError(t, err)
			assert.Equal(t, d.answer, ans)
		})
	}
}

// Test if error is thrown if question can't be solved.
func TestSolution_Error(t *testing.T) {
	_, err := solution("any question")
	assert.Error(t, err)
}

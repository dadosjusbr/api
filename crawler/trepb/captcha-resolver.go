package main

import (
	"fmt"
	"strconv"
	"strings"
)

var monthStr = []string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}

// solution solves captcha questions
func solution(question string) (string, error) {
	if strings.Contains(question, "mês") {
		ans, err := monthCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	}
	if strings.Contains(question, "sequência") {
		ans, err := sequenceCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	}
	if strings.Contains(question, "Quanto") {
		ans, err := arithmeticCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	}
	if strings.Contains(question, "par ou") {
		ans, err := evenOrOddCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	}

	return "", fmt.Errorf("Couldn't fit question in any algorithm: %s", question)
}

// monthCaptcha solves month questions.
func monthCaptcha(question string) (string, error) {
	monthIndex := -1
	for i, month := range monthStr {
		if strings.Contains(question, month) {
			monthIndex = i
			break
		}
	}
	if monthIndex == -1 {
		return "", fmt.Errorf("detected as month question: no month found - %s", question)
	}

	if strings.Contains(question, "Antes") {
		if monthIndex == 0 {
			return monthStr[11], nil
		}
		return monthStr[monthIndex-1], nil
	}

	if strings.Contains(question, "Depois") {
		if monthIndex == 11 {
			return monthStr[0], nil
		}
		return monthStr[monthIndex+1], nil
	}

	return "", fmt.Errorf("detected as month question: couldn't solve question - %s", question)
}

// sequenceCaptcha solves sequence questions.
func sequenceCaptcha(question string) (string, error) {
	args := strings.Split(question, ", ")
	num, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		return "", fmt.Errorf("sequence question error (%s): couldn't convert string(%s) to int: %q", question, args[len(args)-1], err)
	}
	return strconv.Itoa(num + 1), nil
}

// arithmeticCaptcha solves arithmetic questions.
func arithmeticCaptcha(question string) (string, error) {
	args := strings.Split(question, " ")

	arg1 := args[len(args)-3]
	op := args[len(args)-2]
	arg2 := strings.TrimSuffix(args[len(args)-1], "?")

	num1, err := strconv.Atoi(arg1)
	if err != nil {
		return "", fmt.Errorf("arithmetic question error (%s): couldn't convert string(%s) to int: %q", question, arg1, err)
	}

	num2, err := strconv.Atoi(arg2)
	if err != nil {
		return "", fmt.Errorf("arithmetic question error (%s): couldn't convert string(%s) to int: %q", question, arg2, err)
	}

	if op == "+" {
		return strconv.Itoa(num1 + num2), nil
	}
	return strconv.Itoa(num1 - num2), nil
}

//evenOrOddCaptcha solves even or odd questions
func evenOrOddCaptcha(question string) (string, error) {
	args := strings.Split(question, " ")

	num, err := strconv.Atoi(args[2])
	if err != nil {
		return "", fmt.Errorf("even/odd question(%s): couldn't convert string(%s) to int: %q", question, args[2], err)
	}

	if num%2 == 0 {
		return "par", nil
	}
	return "impar", nil
}

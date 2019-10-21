package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Gives solution to question.
func solution(question string) (string, error) {
	if strings.Contains(question, "mês") {
		ans, err := monthCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	} else if strings.Contains(question, "sequência") {
		ans, err := sequenceCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	} else if strings.Contains(question, "Quanto") {
		ans, err := arithmeticCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	} else if strings.Contains(question, "par ou") {
		ans, err := evenOrOddCaptcha(question)
		if err != nil {
			return "", err
		}
		return ans, nil
	}

	return "", fmt.Errorf("Couldn't fit question in any algorithm: %s", question)
}

// Solve month questions.
func monthCaptcha(question string) (string, error) {
	monthIndex := -1
	for i, month := range monthStr {
		if strings.Contains(question, month) {
			monthIndex = i
			break
		}
	}

	if strings.Contains(question, "Antes") {
		if monthIndex == 0 {
			return monthStr[11], nil
		}
		return monthStr[monthIndex-1], nil
	}
	if monthIndex == 11 {
		return monthStr[0], nil
	}
	return monthStr[monthIndex+1], nil
}

// Solve sequence questions.
func sequenceCaptcha(question string) (string, error) {
	args := strings.Split(question, ", ")
	num, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		return "", fmt.Errorf("Couldn't convert string(%s) to int: %q", args[len(args)-1], err)
	}
	return strconv.Itoa(num + 1), nil
}

// Solve arithmetic questions.
func arithmeticCaptcha(question string) (string, error) {
	args := strings.Split(question, " ")

	arg1 := args[len(args)-3]
	op := args[len(args)-2]
	arg2 := args[len(args)-1]
	arg2 = arg2[:len(arg2)-1] // Removing question mark

	num1, err := strconv.Atoi(arg1)
	if err != nil {
		return "", fmt.Errorf("Couldn't convert string(%s) to int: %q", arg1, err)
	}

	num2, err := strconv.Atoi(arg2) //Removing question mark
	if err != nil {
		return "", fmt.Errorf("Couldn't convert string(%s) to int: %q", arg2, err)
	}

	if op == "+" {
		return strconv.Itoa(num1 + num2), nil
	}
	return strconv.Itoa(num1 - num2), nil
}

//Solve even or odd questions
func evenOrOddCaptcha(question string) (string, error) {
	args := strings.Split(question, " ")

	num, err := strconv.Atoi(args[2])
	if err != nil {
		return "", fmt.Errorf("Couldn't convert string(%s) to int: %q", args[2], err)
	}

	if num%2 == 0 {
		return "par", nil
	}
	return "impar", nil
}

package utils

import (
	"encoding/base64"
	"strings"
)

func GenerateColorByFirstLetter(input string) string {
	// Convert the input string to uppercase for case-insensitivity
	firstLetter := strings.ToUpper(input[:1])

	switch firstLetter {
	case "J", "Z", "A":
		return "365397"
	case "B":
		return "00a9f1"
	case "C":
		return "006db3"
	case "D":
		return "737373"
	case "E":
		return "4285f4"
	case "F":
		return "e0452c"
	case "G":
		return "ff3333"
	case "W", "H":
		return "48b6ed"
	case "I":
		return "ce1a19"
	case "K":
		return "ed4584"
	case "L":
		return "ff9700"
	case "X", "M":
		return "083790"
	case "N":
		return "00acf4"
	case "O":
		return "396d9a"
	case "P":
		return "0d84de"
	case "Q":
		return "ea0066"
	case "R":
		return "2f2f2f"
	case "S":
		return "6bbd6d"
	case "T":
		return "304c68"
	case "Y", "U":
		return "207dc5"
	case "V":
		return "1277bc"
	default:
		return "365397"
	}
}

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func CapitalizeEachWord(input string) string {
	// Split the input string into words
	words := strings.Fields(input)

	// Capitalize the first letter of each word
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	// Join the words back into a single string
	result := strings.Join(words, " ")

	return result
}

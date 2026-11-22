package erasmus

import (
	"strconv"
	"strings"
)

func entero(texto string) int {
	numero, _ := strconv.Atoi(texto)
	return numero
}

func decimal(texto string) float64 {
	numero, _ := strconv.ParseFloat(strings.Replace(texto, ",", ".", 1), 64)
	return numero
}

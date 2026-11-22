package erasmus

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	lineaCursos = regexp.MustCompile(`Cursos .*\n *(.+?) {2,}`)
	digito      = regexp.MustCompile(`\d`)
)

func entero(texto string) int {
	numero, _ := strconv.Atoi(texto)
	return numero
}

func decimal(texto string) float64 {
	numero, _ := strconv.ParseFloat(strings.Replace(texto, ",", ".", 1), 64)
	return numero
}

func leerCursos(bloque string) []int {
	cursos := []int{}
	for _, curso := range digito.FindAllString(lineaCursos.FindStringSubmatch(bloque)[1], -1) {
		cursos = append(cursos, entero(curso))
	}
	return cursos
}

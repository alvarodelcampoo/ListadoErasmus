package erasmus

import (
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	lineaDestino   = regexp.MustCompile(`^(\d{3}) +\S+ +\S+ +(.+?) {2,}(\d+) +(\d+)\b.*?(?: {3,}(\d+(?:[.,]\d+)?))? *\n`)
	lineaCursos    = regexp.MustCompile(`Cursos .*\n *(.+?) {2,}`)
	digito         = regexp.MustCompile(`\d`)
	lineaIdiomas   = regexp.MustCompile(`Idiomas *\n(.*)`)
	idiomaConNivel = regexp.MustCompile(`(\p{Lu}+) \((B1|B2|C1|C2)\)`)
	nivelPorNombre = map[string]NivelIdioma{"B1": NivelB1, "B2": NivelB2, "C1": NivelC1, "C2": NivelC2}
)

func LeerDestinos(texto io.Reader) (map[int]Destino, error) {
	contenido, err := io.ReadAll(texto)
	destinos := map[int]Destino{}
	for _, bloque := range strings.Split(string(contenido), "\n 015-")[1:] {
		numeroOrden, destino := leerDestino(bloque)
		destinos[numeroOrden] = destino
	}
	return destinos, err
}

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
	for _, curso := range digito.FindAllString(buscar(lineaCursos, bloque)[1], -1) {
		cursos = append(cursos, entero(curso))
	}
	return cursos
}

func leerIdiomas(bloque string) map[string]NivelIdioma {
	idiomas := map[string]NivelIdioma{}
	for _, idioma := range idiomaConNivel.FindAllStringSubmatch(buscar(lineaIdiomas, bloque)[1], -1) {
		idiomas[idioma[1]] = nivelPorNombre[idioma[2]]
	}
	return idiomas
}

func leerDestino(bloque string) (int, Destino) {
	campos := buscar(lineaDestino, bloque)
	return entero(campos[1]), Destino{
		institucion: campos[2],
		plazas:      entero(campos[3]),
		duracion:    entero(campos[4]),
		notaMinima:  decimal(campos[5]),
		cursos:      leerCursos(bloque),
		idiomas:     leerIdiomas(bloque),
	}
}

func buscar(expresion *regexp.Regexp, bloque string) []string {
	campos := expresion.FindStringSubmatch(bloque)
	return append(campos, make([]string, expresion.NumSubexp()+1-len(campos))...)
}

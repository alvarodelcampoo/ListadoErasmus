package erasmus

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	lineaDestino        = regexp.MustCompile(`^(\d{3}) +\S+ +\S+ +(.+?) {2,}(\d+) +(\d+)\b.*?(?: {3,}(\S+(?: \S+)?))? *\n`)
	tablaCursos         = regexp.MustCompile(`Cursos .*\n((?: *\d.*\n)+)`)
	filaCursos          = regexp.MustCompile(`(?m)^ *(.+?) {2,}(.+)$`)
	codigoTitulacion    = regexp.MustCompile(`(\w+)/\d{2}`)
	numeroDecimal       = regexp.MustCompile(`^(\d+([.,]\d+)?)?$`)
	digito              = regexp.MustCompile(`\d`)
	lineaIdiomas        = regexp.MustCompile(`Idiomas *\n(.*)`)
	idiomaConNivel      = regexp.MustCompile(`(\p{Lu}+) \((B1|B2|C1|C2)\)`)
	nivelPorNombre      = map[string]NivelIdioma{"B1": NivelB1, "B2": NivelB2, "C1": NivelC1, "C2": NivelC2}
	titulacionPorCodigo = map[string]Titulacion{
		"296": IngenieriaInformatica,
		"216": InformaticaYADE,
		"297": InformaticaYMatematicas,
		"221": Telecomunicaciones,
	}
)

type LineaNoEncontradaError struct {
	Bloque int
	Linea  string
}

func (e *LineaNoEncontradaError) Error() string {
	return fmt.Sprintf("el bloque %d del listado no tiene la línea de %s", e.Bloque, e.Linea)
}

type NotaMinimaIlegibleError struct {
	Bloque int
	Texto  string
}

func (e *NotaMinimaIlegibleError) Error() string {
	return fmt.Sprintf("el bloque %d del listado tiene una nota mínima que no es un número: %q", e.Bloque, e.Texto)
}

func LeerDestinos(texto io.Reader) (map[int]Destino, error) {
	contenido, err := io.ReadAll(texto)
	errores := []error{err}
	destinos := map[int]Destino{}
	for posicion, bloque := range strings.Split(string(contenido), "\n 015-")[1:] {
		numeroOrden, destino, errBloque := leerDestino(posicion+1, bloque)
		destinos[numeroOrden] = destino
		errores = append(errores, errBloque)
	}
	return destinos, errors.Join(errores...)
}

func entero(texto string) int {
	numero, _ := strconv.Atoi(texto)
	return numero
}

func decimal(texto string) float64 {
	valor, _ := strconv.ParseFloat(strings.Replace(texto, ",", ".", 1), 64)
	return valor
}

func leerCursos(bloque string) map[Titulacion][]int {
	cursos := map[Titulacion][]int{}
	for _, fila := range filaCursos.FindAllStringSubmatch(buscar(tablaCursos, bloque)[1], -1) {
		for _, titulacion := range leerTitulaciones(fila[2]) {
			cursos[titulacion] = leerNumeros(fila[1])
		}
	}
	return cursos
}

func leerTitulaciones(texto string) []Titulacion {
	titulaciones := []Titulacion{}
	for _, codigo := range codigoTitulacion.FindAllStringSubmatch(texto, -1) {
		titulacion, deGrado := titulacionPorCodigo[codigo[1]]
		if deGrado {
			titulaciones = append(titulaciones, titulacion)
		}
	}
	return titulaciones
}

func leerNumeros(texto string) []int {
	numeros := []int{}
	for _, numero := range digito.FindAllString(texto, -1) {
		numeros = append(numeros, entero(numero))
	}
	return numeros
}

func leerIdiomas(bloque string) map[string]NivelIdioma {
	idiomas := map[string]NivelIdioma{}
	for _, idioma := range idiomaConNivel.FindAllStringSubmatch(buscar(lineaIdiomas, bloque)[1], -1) {
		idiomas[idioma[1]] = nivelPorNombre[idioma[2]]
	}
	return idiomas
}

func leerDestino(posicion int, bloque string) (int, Destino, error) {
	campos := buscar(lineaDestino, bloque)
	return entero(campos[1]), Destino{
		institucion: campos[2],
		plazas:      entero(campos[3]),
		duracion:    entero(campos[4]),
		notaMinima:  decimal(campos[5]),
		cursos:      leerCursos(bloque),
		idiomas:     leerIdiomas(bloque),
	}, errors.Join(
		errorSi(!numeroDecimal.MatchString(campos[5]), &NotaMinimaIlegibleError{Bloque: posicion, Texto: campos[5]}),
		errorSi(!lineaDestino.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "destino"}),
		errorSi(!tablaCursos.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "cursos"}),
		errorSi(!lineaIdiomas.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "idiomas"}),
	)
}

func buscar(linea *regexp.Regexp, bloque string) []string {
	campos := linea.FindStringSubmatch(bloque)
	return append(campos, make([]string, linea.NumSubexp()+1-len(campos))...)
}

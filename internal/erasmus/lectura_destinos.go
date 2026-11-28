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
	lineaDestino        = regexp.MustCompile(`^(?P<numeroOrden>\d{3}) +\S+ +\S+ +(?P<institucion>.+?) {2,}(?P<plazas>\d+) +(?P<duracion>\d+)\b.*?(?: {3,}(?P<notaMinima>\S+(?: \S+)?))? *\n`)
	tablaCursos         = regexp.MustCompile(`Cursos .*\n(?P<filas>(?: *\d.*\n)+)`)
	filaCursos          = regexp.MustCompile(`(?m)^ *(?P<cursos>.+?) {2,}(?P<titulaciones>.+)$`)
	codigoTitulacion    = regexp.MustCompile(`(?P<codigo>\w+)/\d{2}`)
	numeroDecimal       = regexp.MustCompile(`^(\d+([.,]\d+)?)?$`)
	digito              = regexp.MustCompile(`(?P<numero>\d)`)
	lineaIdiomas        = regexp.MustCompile(`Idiomas *\n(?P<idiomas>.*)`)
	idiomaConNivel      = regexp.MustCompile(`(?P<idioma>\p{Lu}+) \((?P<nivel>B1|B2|C1|C2)\)`)
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

type NumeroOrdenRepetidoError struct {
	NumeroOrden int
}

func (e *NumeroOrdenRepetidoError) Error() string {
	return fmt.Sprintf("el nº de orden %d aparece más de una vez en el listado", e.NumeroOrden)
}

func LeerDestinos(texto io.Reader) (map[int]Destino, error) {
	contenido, err := io.ReadAll(texto)
	errores := []error{err}
	destinos := map[int]Destino{}
	for posicion, bloque := range strings.Split(string(contenido), "\n 015-")[1:] {
		numeroOrden, destino, errBloque := leerDestino(posicion+1, bloque)
		_, repetido := destinos[numeroOrden]
		destinos[numeroOrden] = destino
		errores = append(errores, errBloque, errorSi(repetido, &NumeroOrdenRepetidoError{NumeroOrden: numeroOrden}))
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
	for _, fila := range buscarTodas(filaCursos, buscar(tablaCursos, bloque)["filas"]) {
		for _, titulacion := range leerTitulaciones(fila["titulaciones"]) {
			cursos[titulacion] = leerNumeros(fila["cursos"])
		}
	}
	return cursos
}

func leerTitulaciones(texto string) []Titulacion {
	titulaciones := []Titulacion{}
	for _, codigo := range buscarTodas(codigoTitulacion, texto) {
		titulacion, deGrado := titulacionPorCodigo[codigo["codigo"]]
		if deGrado {
			titulaciones = append(titulaciones, titulacion)
		}
	}
	return titulaciones
}

func leerNumeros(texto string) []int {
	numeros := []int{}
	for _, numero := range buscarTodas(digito, texto) {
		numeros = append(numeros, entero(numero["numero"]))
	}
	return numeros
}

func leerIdiomas(bloque string) map[string]NivelIdioma {
	idiomas := map[string]NivelIdioma{}
	for _, idioma := range buscarTodas(idiomaConNivel, buscar(lineaIdiomas, bloque)["idiomas"]) {
		idiomas[idioma["idioma"]] = nivelPorNombre[idioma["nivel"]]
	}
	return idiomas
}

func leerDestino(posicion int, bloque string) (int, Destino, error) {
	campos := buscar(lineaDestino, bloque)
	return entero(campos["numeroOrden"]), Destino{
		institucion: campos["institucion"],
		plazas:      entero(campos["plazas"]),
		duracion:    entero(campos["duracion"]),
		notaMinima:  decimal(campos["notaMinima"]),
		cursos:      leerCursos(bloque),
		idiomas:     leerIdiomas(bloque),
	}, comprobarBloque(posicion, bloque, campos["notaMinima"])
}

func comprobarBloque(posicion int, bloque string, notaMinima string) error {
	return errors.Join(
		errorSi(!numeroDecimal.MatchString(notaMinima), &NotaMinimaIlegibleError{Bloque: posicion, Texto: notaMinima}),
		errorSi(!lineaDestino.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "destino"}),
		errorSi(!tablaCursos.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "cursos"}),
		errorSi(!lineaIdiomas.MatchString(bloque), &LineaNoEncontradaError{Bloque: posicion, Linea: "idiomas"}),
	)
}

func buscar(linea *regexp.Regexp, texto string) map[string]string {
	return camposPorNombre(linea, linea.FindStringSubmatch(texto))
}

func buscarTodas(linea *regexp.Regexp, texto string) []map[string]string {
	todas := []map[string]string{}
	for _, coincidencia := range linea.FindAllStringSubmatch(texto, -1) {
		todas = append(todas, camposPorNombre(linea, coincidencia))
	}
	return todas
}

func camposPorNombre(linea *regexp.Regexp, coincidencia []string) map[string]string {
	campos := map[string]string{}
	for posicion, valor := range coincidencia {
		campos[linea.SubexpNames()[posicion]] = valor
	}
	return campos
}

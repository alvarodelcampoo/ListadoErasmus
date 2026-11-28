package erasmus

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"
)

type PlazaAdjudicada struct {
	numeroOrden int
	notaTotal   float64
}

type Adjudicaciones []map[string]PlazaAdjudicada

var (
	filaAdjudicacion  = regexp.MustCompile(`(?m)^(?P<estudiante>\S.*?) +\d+\.\d{3} +\d\.\d{3} +(?P<notaTotal>\d+\.\d{3}) .*?015-(?P<numeroOrden>\d{3})`)
	destinoAdjudicado = regexp.MustCompile(`015-\d{3}`)
)

type FilaIlegibleError struct {
	Ruta  string
	Linea int
}

func (e *FilaIlegibleError) Error() string {
	return fmt.Sprintf("la línea %d de %s tiene un destino adjudicado pero no se ha podido leer", e.Linea, e.Ruta)
}

func LeerAdjudicaciones(rutas ...string) (Adjudicaciones, error) {
	adjudicaciones := make(Adjudicaciones, len(rutas))
	errores := make([]error, len(rutas))
	for i, ruta := range rutas {
		adjudicaciones[i], errores[i] = leerAdjudicacion(ruta)
	}
	return adjudicaciones, errors.Join(errores...)
}

func leerAdjudicacion(ruta string) (map[string]PlazaAdjudicada, error) {
	contenido, err := os.ReadFile(ruta)
	plazas := map[string]PlazaAdjudicada{}
	for _, fila := range buscarTodas(filaAdjudicacion, string(contenido)) {
		plazas[fila["estudiante"]] = PlazaAdjudicada{numeroOrden: entero(fila["numeroOrden"]), notaTotal: decimal(fila["notaTotal"])}
	}
	return plazas, errors.Join(err, comprobarFilas(ruta, string(contenido)))
}

func comprobarFilas(ruta string, contenido string) error {
	return errorDeCada(strings.Split(contenido, "\n"), func(posicion int, linea string) error {
		return errorSi(destinoAdjudicado.MatchString(linea) && !filaAdjudicacion.MatchString(linea), &FilaIlegibleError{Ruta: ruta, Linea: posicion + 1})
	})
}

func (adjudicaciones Adjudicaciones) ultimaPlazaDeCadaEstudiante() map[string]PlazaAdjudicada {
	plazas := map[string]PlazaAdjudicada{}
	for _, adjudicacion := range adjudicaciones {
		maps.Copy(plazas, adjudicacion)
	}
	return plazas
}

func (adjudicaciones Adjudicaciones) notasPorDestino() map[int][]float64 {
	notas := map[int][]float64{}
	for _, plaza := range adjudicaciones.ultimaPlazaDeCadaEstudiante() {
		notas[plaza.numeroOrden] = append(notas[plaza.numeroOrden], plaza.notaTotal)
	}
	return notas
}

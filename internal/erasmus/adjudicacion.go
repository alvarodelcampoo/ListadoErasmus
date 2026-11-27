package erasmus

import (
	"errors"
	"maps"
	"os"
	"regexp"
)

type PlazaAdjudicada struct {
	numeroOrden int
	notaTotal   float64
}

type Adjudicaciones []map[string]PlazaAdjudicada

var filaAdjudicacion = regexp.MustCompile(`(?m)^(?P<estudiante>\S.*?) +\d+\.\d{3} +\d\.\d{3} +(?P<notaTotal>\d+\.\d{3}) .*?015-(?P<numeroOrden>\d{3})`)

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
	return plazas, err
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

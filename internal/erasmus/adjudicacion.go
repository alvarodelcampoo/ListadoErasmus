package erasmus

import (
	"errors"
	"io"
	"maps"
	"regexp"
)

type PlazaAdjudicada struct {
	numeroOrden int
	notaTotal   float64
}

type Adjudicaciones []map[string]PlazaAdjudicada

var filaAdjudicacion = regexp.MustCompile(`(?m)^(\S.*?) +\d+\.\d{3} +\d\.\d{3} +(\d+\.\d{3}) .*?015-(\d{3})`)

func LeerAdjudicaciones(textos ...io.Reader) (Adjudicaciones, error) {
	adjudicaciones := make(Adjudicaciones, len(textos))
	errores := make([]error, len(textos))
	for i, texto := range textos {
		adjudicaciones[i], errores[i] = leerAdjudicacion(texto)
	}
	return adjudicaciones, errors.Join(errores...)
}

func leerAdjudicacion(texto io.Reader) (map[string]PlazaAdjudicada, error) {
	contenido, err := io.ReadAll(texto)
	plazas := map[string]PlazaAdjudicada{}
	for _, fila := range filaAdjudicacion.FindAllStringSubmatch(string(contenido), -1) {
		plazas[fila[1]] = PlazaAdjudicada{numeroOrden: entero(fila[3]), notaTotal: decimal(fila[2])}
	}
	return plazas, err
}

func ultimaAdjudicacion(adjudicaciones ...map[string]PlazaAdjudicada) map[string]PlazaAdjudicada {
	plazas := map[string]PlazaAdjudicada{}
	for _, adjudicacion := range adjudicaciones {
		maps.Copy(plazas, adjudicacion)
	}
	return plazas
}

func notasPorDestino(plazas map[string]PlazaAdjudicada) map[int][]float64 {
	notas := map[int][]float64{}
	for _, plaza := range plazas {
		notas[plaza.numeroOrden] = append(notas[plaza.numeroOrden], plaza.notaTotal)
	}
	return notas
}

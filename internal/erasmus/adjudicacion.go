package erasmus

import (
	"io"
	"regexp"
)

type PlazaAdjudicada struct {
	numeroOrden int
	notaTotal   float64
}

var filaAdjudicacion = regexp.MustCompile(`(?m)^(\S.*?) +\d+\.\d{3} +\d\.\d{3} +(\d+\.\d{3}) .*?015-(\d{3})`)

func LeerAdjudicacion(texto io.Reader) (map[string]PlazaAdjudicada, error) {
	contenido, err := io.ReadAll(texto)
	plazas := map[string]PlazaAdjudicada{}
	for _, fila := range filaAdjudicacion.FindAllStringSubmatch(string(contenido), -1) {
		plazas[fila[1]] = PlazaAdjudicada{numeroOrden: entero(fila[3]), notaTotal: decimal(fila[2])}
	}
	return plazas, err
}

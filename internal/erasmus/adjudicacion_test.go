package erasmus

import (
	"fmt"
	"os"
	"testing"
)

func leerAdjudicacionReal(t *testing.T, numero int) map[string]PlazaAdjudicada {
	fichero, _ := os.Open(fmt.Sprintf("../../datos/adjudicacion%d.txt", numero))
	defer fichero.Close()
	plazas, err := LeerAdjudicacion(fichero)
	comprobar(t, err, nil)
	return plazas
}

func TestLeerAdjudicacion(t *testing.T) {
	plazas := leerAdjudicacionReal(t, 3)
	comprobar(t, plazas["ESTUDIANTE 166"], PlazaAdjudicada{numeroOrden: 61, notaTotal: 8.221})
	comprobar(t, plazas["ESTUDIANTE 021"], PlazaAdjudicada{numeroOrden: 97, notaTotal: 5.992})
}

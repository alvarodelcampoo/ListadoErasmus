package erasmus

import (
	"fmt"
	"os"
	"slices"
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

func TestLeerTodasLasAdjudicaciones(t *testing.T) {
	comprobar(t, len(leerAdjudicacionReal(t, 1)), 156)
	comprobar(t, len(leerAdjudicacionReal(t, 2)), 69)
	comprobar(t, len(leerAdjudicacionReal(t, 3)), 9)
}

func TestUltimaAdjudicacion(t *testing.T) {
	plazas := ultimaAdjudicacion(leerAdjudicacionReal(t, 1), leerAdjudicacionReal(t, 2), leerAdjudicacionReal(t, 3))
	comprobar(t, len(plazas), 172)
	comprobar(t, plazas["ESTUDIANTE 020"], PlazaAdjudicada{numeroOrden: 49, notaTotal: 8.13})
}

func TestNotasPorDestino(t *testing.T) {
	notas := notasPorDestino(ultimaAdjudicacion(leerAdjudicacionReal(t, 1), leerAdjudicacionReal(t, 2), leerAdjudicacionReal(t, 3)))
	comprobar(t, len(notas), 84)
	comprobar(t, slices.Sorted(slices.Values(notas[49])), []float64{8.13, 8.978, 9.019})
}

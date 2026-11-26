package erasmus

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func leerAdjudicacionReal(t *testing.T, numero int) map[string]PlazaAdjudicada {
	fichero, _ := os.Open(fmt.Sprintf("../../datos/adjudicacion%d.txt", numero))
	defer fichero.Close()
	plazas, err := LeerAdjudicacion(fichero)
	assert.Equal(t, nil, err)
	return plazas
}

func TestLeerAdjudicacion(t *testing.T) {
	plazas := leerAdjudicacionReal(t, 3)
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 61, notaTotal: 8.221}, plazas["ESTUDIANTE 166"])
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 97, notaTotal: 5.992}, plazas["ESTUDIANTE 021"])
}

func TestLeerTodasLasAdjudicaciones(t *testing.T) {
	assert.Equal(t, 156, len(leerAdjudicacionReal(t, 1)))
	assert.Equal(t, 69, len(leerAdjudicacionReal(t, 2)))
	assert.Equal(t, 9, len(leerAdjudicacionReal(t, 3)))
}

func TestUltimaAdjudicacion(t *testing.T) {
	plazas := ultimaAdjudicacion(leerAdjudicacionReal(t, 1), leerAdjudicacionReal(t, 2), leerAdjudicacionReal(t, 3))
	assert.Equal(t, 172, len(plazas))
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 49, notaTotal: 8.13}, plazas["ESTUDIANTE 020"])
}

func TestNotasPorDestino(t *testing.T) {
	notas := notasPorDestino(ultimaAdjudicacion(leerAdjudicacionReal(t, 1), leerAdjudicacionReal(t, 2), leerAdjudicacionReal(t, 3)))
	assert.Equal(t, 84, len(notas))
	assert.Equal(t, []float64{8.13, 8.978, 9.019}, slices.Sorted(slices.Values(notas[49])))
}

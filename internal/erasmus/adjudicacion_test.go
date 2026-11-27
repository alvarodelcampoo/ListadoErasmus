package erasmus

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func leerAdjudicacionesReales(t *testing.T) Adjudicaciones {
	rutas, _ := filepath.Glob("../../datos/adjudicacion*.txt")
	adjudicaciones, err := LeerAdjudicaciones(rutas...)
	assert.Equal(t, nil, err)
	return adjudicaciones
}

func TestLeerAdjudicacion(t *testing.T) {
	tercera := leerAdjudicacionesReales(t)[2]
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 61, notaTotal: 8.221}, tercera["ESTUDIANTE 166"])
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 97, notaTotal: 5.992}, tercera["ESTUDIANTE 021"])
}

func TestLeerTodasLasAdjudicaciones(t *testing.T) {
	adjudicaciones := leerAdjudicacionesReales(t)
	assert.Equal(t, 3, len(adjudicaciones))
	assert.Equal(t, 156, len(adjudicaciones[0]))
	assert.Equal(t, 69, len(adjudicaciones[1]))
	assert.Equal(t, 9, len(adjudicaciones[2]))
}

func TestUltimaPlazaDeCadaEstudiante(t *testing.T) {
	plazas := leerAdjudicacionesReales(t).ultimaPlazaDeCadaEstudiante()
	assert.Equal(t, 172, len(plazas))
	assert.Equal(t, PlazaAdjudicada{numeroOrden: 49, notaTotal: 8.13}, plazas["ESTUDIANTE 020"])
}

func TestNotasPorDestino(t *testing.T) {
	notas := leerAdjudicacionesReales(t).notasPorDestino()
	assert.Equal(t, 84, len(notas))
	assert.Equal(t, []float64{8.13, 8.978, 9.019}, slices.Sorted(slices.Values(notas[49])))
}

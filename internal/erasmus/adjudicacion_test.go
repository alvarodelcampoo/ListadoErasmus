package erasmus

import (
	"errors"
	"os"
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
	munich, _ := destinoDelListado(t, "TECHNISCHE UNIVERSITÄT MÜNCHEN", IngenieriaInformatica)
	stip, _ := destinoDelListado(t, "REPUBLIC OF MACEDONIA GOCE DELCEV STATE UNIVERSITY STIP", IngenieriaInformatica)
	assert.Equal(t, PlazaAdjudicada{numeroOrden: munich, notaTotal: 8.221}, tercera["ESTUDIANTE 166"])
	assert.Equal(t, PlazaAdjudicada{numeroOrden: stip, notaTotal: 5.992}, tercera["ESTUDIANTE 021"])
}

func TestLeerTodasLasAdjudicaciones(t *testing.T) {
	adjudicaciones := leerAdjudicacionesReales(t)
	adjudicacionesPublicadas := 3
	assert.Equal(t, adjudicacionesPublicadas, len(adjudicaciones))
	primera, segunda, tercera := adjudicaciones[0], adjudicaciones[1], adjudicaciones[2]
	estudiantesEnLaPrimera, estudiantesEnLaSegunda, estudiantesEnLaTercera := 156, 69, 9
	assert.Equal(t, estudiantesEnLaPrimera, len(primera))
	assert.Equal(t, estudiantesEnLaSegunda, len(segunda))
	assert.Equal(t, estudiantesEnLaTercera, len(tercera))
}

func TestUltimaPlazaDeCadaEstudiante(t *testing.T) {
	plazas := leerAdjudicacionesReales(t).ultimaPlazaDeCadaEstudiante()
	aachenInformatica, _ := destinoDelListado(t, institucionAachen, IngenieriaInformatica)
	estudiantesDistintos := 172
	assert.Equal(t, estudiantesDistintos, len(plazas))
	assert.Equal(t, PlazaAdjudicada{numeroOrden: aachenInformatica, notaTotal: 8.13}, plazas["ESTUDIANTE 020"])
}

func TestNotasPorDestino(t *testing.T) {
	notas := leerAdjudicacionesReales(t).notasPorDestino()
	aachenInformatica, _ := destinoDelListado(t, institucionAachen, IngenieriaInformatica)
	destinosConAdjudicados := 84
	assert.Equal(t, destinosConAdjudicados, len(notas))
	assert.Equal(t, []float64{8.13, 8.978, 9.019}, slices.Sorted(slices.Values(notas[aachenInformatica])))
}

func TestLeerFilaDeAdjudicacionIlegible(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "adjudicacion.txt")
	os.WriteFile(ruta, []byte("cabecera\nESTUDIANTE 001   8,000 2.000   10.000 B LOUVAIN01   015-041   1   216   9   1\n"), os.ModePerm)
	lineaDelEstudiante := 2
	_, err := LeerAdjudicaciones(ruta)
	assert.Equal(t, errors.Join(errors.Join(errors.Join(&FilaIlegibleError{Ruta: ruta, Linea: lineaDelEstudiante}))), err)
}

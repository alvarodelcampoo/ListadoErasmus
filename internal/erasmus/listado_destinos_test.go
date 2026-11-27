package erasmus

import (
	"errors"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

const institucionAachen = "RHEINISCH-WESTFÄLISCHE TECHNISCHE HOCHSCHULE AACHEN"

func destinoCon(destino Destino, cambio func(*Destino)) Destino {
	cambio(&destino)
	return destino
}

func soloAachenTeleco(t *testing.T) map[int]Destino {
	numeroAachenTeleco, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	return map[int]Destino{numeroAachenTeleco: aachenTeleco}
}

func TestNewListadoDestinos(t *testing.T) {
	numeroAachenTeleco, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	numeroAachenInformatica, aachenInformatica := destinoDelListado(t, institucionAachen, IngenieriaInformatica)
	sinPlazas := destinoCon(aachenTeleco, func(destino *Destino) { destino.plazas = 0 })
	casos := []struct {
		nombre         string
		destinos       map[int]Destino
		adjudicaciones Adjudicaciones
		errorEsperado  error
	}{
		{"destino real con su nº de orden", map[int]Destino{numeroAachenTeleco: aachenTeleco}, nil, nil},
		{"nº de orden mínimo permitido", map[int]Destino{1: aachenTeleco}, nil, nil},
		{"nº de orden no positivo", map[int]Destino{0: aachenTeleco}, nil, errors.Join(&NumeroOrdenNoPositivoError{NumeroOrden: 0})},
		{"destino con datos imposibles", map[int]Destino{numeroAachenTeleco: sinPlazas}, nil, errors.Join(errors.Join(&PlazasNoPositivasError{Plazas: 0}))},
		{"un destino imposible junto a otro válido", map[int]Destino{numeroAachenTeleco: sinPlazas, numeroAachenInformatica: aachenInformatica}, nil, errors.Join(errors.Join(&PlazasNoPositivasError{Plazas: 0}))},
		{"notas de un destino del listado", map[int]Destino{numeroAachenTeleco: aachenTeleco}, Adjudicaciones{{"ESTUDIANTE 1": {numeroOrden: numeroAachenTeleco, notaTotal: 9.5}, "ESTUDIANTE 2": {numeroOrden: numeroAachenTeleco, notaTotal: 8}}}, nil},
		{"notas de un destino que no está en el listado", map[int]Destino{numeroAachenTeleco: aachenTeleco}, Adjudicaciones{{"ESTUDIANTE 1": {numeroOrden: numeroAachenInformatica, notaTotal: 8}}}, errors.Join(&NotasSinDestinoError{NumeroOrden: numeroAachenInformatica})},
		{"nota imposible de un destino", map[int]Destino{numeroAachenTeleco: aachenTeleco}, Adjudicaciones{{"ESTUDIANTE 1": {numeroOrden: numeroAachenTeleco, notaTotal: 15}}}, errors.Join(errors.Join(errors.Join(&NotaConBonificacionFueraDeRangoError{Nota: 15, Posicion: 1})))},
		{"destino sin institución", map[int]Destino{numeroAachenTeleco: destinoCon(aachenTeleco, func(destino *Destino) { destino.institucion = " " })}, nil, errors.Join(errors.Join(ErrDestinoSinInstitucion))},
		{"duración no positiva", map[int]Destino{numeroAachenTeleco: destinoCon(aachenTeleco, func(destino *Destino) { destino.duracion = 0 })}, nil, errors.Join(errors.Join(&DuracionNoPositivaError{Duracion: 0}))},
		{"curso no positivo", map[int]Destino{numeroAachenTeleco: destinoCon(aachenTeleco, func(destino *Destino) { destino.cursos = map[Titulacion][]int{Telecomunicaciones: {0}} })}, nil, errors.Join(errors.Join(errors.Join(&CursoNoPositivoError{Curso: 0})))},
		{"nota mínima imposible", map[int]Destino{numeroAachenTeleco: destinoCon(aachenTeleco, func(destino *Destino) { destino.notaMinima = 11 })}, nil, errors.Join(errors.Join(&NotaMinimaFueraDeRangoError{Nota: 11}))},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, err := NewListadoDestinos(caso.destinos, caso.adjudicaciones)
			assert.Equal(t, caso.errorEsperado, err)
		})
	}
}

func TestListadoGuardaLasNotasDeCadaDestino(t *testing.T) {
	numeroAachenTeleco, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	numeroAachenInformatica, aachenInformatica := destinoDelListado(t, institucionAachen, IngenieriaInformatica)
	destinos := map[int]Destino{numeroAachenTeleco: aachenTeleco, numeroAachenInformatica: aachenInformatica}
	adjudicaciones := Adjudicaciones{{"ESTUDIANTE 1": {numeroOrden: numeroAachenTeleco, notaTotal: 9.5}, "ESTUDIANTE 2": {numeroOrden: numeroAachenTeleco, notaTotal: 8}}}
	listado, _ := NewListadoDestinos(destinos, adjudicaciones)
	assert.Equal(t, []float64{8, 9.5}, slices.Sorted(slices.Values(listado.destinos[numeroAachenTeleco].notasTotales)))
	assert.Equal(t, true, listado.destinos[numeroAachenInformatica].EsNuevo())
}

type casoDestinosViables struct {
	nombre           string
	notaMedia        float64
	niveles          map[string]NivelIdioma
	viablesEsperados map[int]Destino
	errorEsperado    error
}

func nivelesQueCumplen() map[string]NivelIdioma {
	return map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelB1}
}

func probarDestinosViables(t *testing.T, destinos map[int]Destino, casos []casoDestinosViables) {
	listado := ListadoDestinos{destinos: destinos}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			viables, err := listado.DestinosViables(Estudiante{NotaMedia: caso.notaMedia, Niveles: caso.niveles})
			assert.Equal(t, caso.viablesEsperados, viables)
			assert.Equal(t, caso.errorEsperado, err)
		})
	}
}

func TestNotaMediaDelEstudiante(t *testing.T) {
	soloAachen := soloAachenTeleco(t)
	probarDestinosViables(t, soloAachen, []casoDestinosViables{
		{"nota media válida", 8, nivelesQueCumplen(), soloAachen, nil},
		{"nota media máxima", 10, nivelesQueCumplen(), soloAachen, nil},
		{"nota media negativa", -1, nivelesQueCumplen(), nil, &NotaMediaFueraDeRangoError{Nota: -1}},
		{"nota media por encima de 10", 10.5, nivelesQueCumplen(), nil, &NotaMediaFueraDeRangoError{Nota: 10.5}},
	})
}

func TestRequisitoDeIdiomas(t *testing.T) {
	soloAachen := soloAachenTeleco(t)
	probarDestinosViables(t, soloAachen, []casoDestinosViables{
		{"nivel superior al exigido", 8, map[string]NivelIdioma{"ALEMAN": NivelC1, "INGLES": NivelB2}, soloAachen, nil},
		{"le falta uno de los idiomas exigidos", 8, map[string]NivelIdioma{"ALEMAN": NivelB1}, map[int]Destino{}, nil},
	})
}

func TestNotaMinimaDelDestino(t *testing.T) {
	soloAachen := soloAachenTeleco(t)
	probarDestinosViables(t, soloAachen, []casoDestinosViables{
		{"nota media igual a la nota mínima", 7, nivelesQueCumplen(), soloAachen, nil},
		{"nota media por debajo de la nota mínima", 6.9, nivelesQueCumplen(), map[int]Destino{}, nil},
	})
}

func TestNotaPonderada(t *testing.T) {
	_, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	casos := []struct {
		nombre        string
		destino       Destino
		notaMedia     float64
		niveles       map[string]NivelIdioma
		notaPonderada float64
	}{
		{"mismo nivel que el exigido", aachenTeleco, 7, nivelesQueCumplen(), 7},
		{"supera el nivel exigido en un idioma", aachenTeleco, 8, map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelC1}, 9.5},
		{"supera el nivel en los dos idiomas y cuenta la mayor", aachenTeleco, 8, map[string]NivelIdioma{"ALEMAN": NivelB2, "INGLES": NivelC1}, 9.5},
		{"un idioma que el destino no exige no bonifica", aachenTeleco, 8, map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelB1, "FRANCES": NivelC2}, 8},
		{"destino sin idiomas exigidos y cuenta su mejor idioma", destinoCon(aachenTeleco, func(destino *Destino) { destino.idiomas = nil }), 8, map[string]NivelIdioma{"ALEMAN": NivelB2, "INGLES": NivelC1}, 9.5},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			assert.Equal(t, caso.notaPonderada, caso.destino.notaPonderada(caso.notaMedia, caso.niveles))
		})
	}
}

func TestNotaDeCorte(t *testing.T) {
	_, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	conAdjudicados := destinoCon(aachenTeleco, func(destino *Destino) { destino.notasTotales = []float64{9.5, 8} })
	casos := []struct {
		nombre    string
		destino   Destino
		notaMedia float64
		niveles   map[string]NivelIdioma
		viable    bool
	}{
		{"nota ponderada igual a la nota de corte", conAdjudicados, 8, nivelesQueCumplen(), true},
		{"nota ponderada por debajo de la nota de corte", conAdjudicados, 7.9, nivelesQueCumplen(), false},
		{"solo llega a la nota de corte con la bonificación", conAdjudicados, 7, map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelC1}, true},
		{"destino nuevo sin adjudicados del año anterior", aachenTeleco, 7, nivelesQueCumplen(), true},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			assert.Equal(t, caso.viable, caso.destino.esViable(caso.notaMedia, caso.niveles))
		})
	}
}

func TestDestinoNuevo(t *testing.T) {
	numeroAachenTeleco, aachenTeleco := destinoDelListado(t, institucionAachen, Telecomunicaciones)
	casos := []struct {
		nombre  string
		destino Destino
		nuevo   bool
	}{
		{"sin adjudicados del año anterior es nuevo", aachenTeleco, true},
		{"con adjudicados del año anterior no es nuevo", destinoCon(aachenTeleco, func(destino *Destino) { destino.notasTotales = []float64{8} }), false},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			viables, _ := ListadoDestinos{destinos: map[int]Destino{numeroAachenTeleco: caso.destino}}.DestinosViables(Estudiante{NotaMedia: 8, Niveles: nivelesQueCumplen()})
			assert.Equal(t, caso.nuevo, viables[numeroAachenTeleco].EsNuevo())
		})
	}
}

func TestNotaPonderadaEnDestinosReales(t *testing.T) {
	_, politecnicoDiBari := destinoDelListado(t, "POLITECNICO DI BARI", Telecomunicaciones)
	_, uclouvain := destinoDelListado(t, "Uclouvain", Telecomunicaciones)
	notaMedia := 7.0
	ingles := map[string]NivelIdioma{"INGLES": NivelB2}
	notaConBonificacionB2 := 8.0
	assert.Equal(t, notaConBonificacionB2, politecnicoDiBari.notaPonderada(notaMedia, ingles))
	assert.Equal(t, notaMedia, uclouvain.notaPonderada(notaMedia, ingles))
}

func TestDestinosViablesConDatosReales(t *testing.T) {
	destinos := leerListadoReal(t)
	adjudicaciones := leerAdjudicacionesReales(t)
	listado, err := NewListadoDestinos(destinos, adjudicaciones)
	assert.Equal(t, nil, err)
	estudiante := Estudiante{NotaMedia: 5, Niveles: map[string]NivelIdioma{"PORTUGUES": NivelB2}}
	viables, err := listado.DestinosViables(estudiante)
	assert.Equal(t, nil, err)
	politecnicoDiMilano, _ := destinoDelListado(t, "POLITECNICO DI MILANO", Telecomunicaciones)
	politecnicoDiBari, _ := destinoDelListado(t, "POLITECNICO DI BARI", IngenieriaInformatica)
	yildiz, _ := destinoDelListado(t, "YILDIZ TECHNICAL UNIVERSITY", IngenieriaInformatica)
	ankara, _ := destinoDelListado(t, "ANKARA UNIVERSITY", IngenieriaInformatica)
	esperados := []int{politecnicoDiMilano, politecnicoDiBari, yildiz, ankara}
	assert.Equal(t, slices.Sorted(slices.Values(esperados)), slices.Sorted(maps.Keys(viables)))
	assert.Equal(t, false, viables[politecnicoDiMilano].EsNuevo())
	assert.Equal(t, true, viables[politecnicoDiBari].EsNuevo())
	assert.Equal(t, false, viables[yildiz].EsNuevo())
	assert.Equal(t, true, viables[ankara].EsNuevo())
}

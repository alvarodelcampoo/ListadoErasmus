package erasmus

import (
	"errors"
	"fmt"
	"testing"
)

func destinoValido() Destino {
	return Destino{
		institucion: "RHEINISCH-WESTFÄLISCHE TECHNISCHE HOCHSCHULE AACHEN",
		plazas:      2,
		duracion:    CursoCompleto,
		cursos:      []int{3, 4},
		notaMinima:  7,
		idiomas:     map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelB1},
	}
}

func destinoCon(cambio func(*Destino)) Destino {
	destino := destinoValido()
	cambio(&destino)
	return destino
}

func comprobar(t *testing.T, obtenido, esperado any) {
	t.Helper()
	if fmt.Sprint(obtenido) != fmt.Sprint(esperado) {
		t.Errorf("se esperaba %v, se obtuvo %v", esperado, obtenido)
	}
}

func TestNewListadoDestinos(t *testing.T) {
	sinPlazas := destinoCon(func(destino *Destino) { destino.plazas = 0 })
	casos := []struct {
		nombre           string
		destinos         map[int]Destino
		erroresEsperados []error
	}{
		{"destino real con su nº de orden", map[int]Destino{6: destinoValido()}, nil},
		{"nº de orden mínimo permitido", map[int]Destino{1: destinoValido()}, nil},
		{"nº de orden no positivo", map[int]Destino{0: destinoValido()}, []error{&NumeroOrdenNoPositivoError{NumeroOrden: 0}}},
		{"destino con datos imposibles", map[int]Destino{6: sinPlazas}, []error{&PlazasNoPositivasError{Plazas: 0}}},
		{"un destino imposible entre varios válidos", map[int]Destino{1: destinoValido(), 2: destinoValido(), 6: sinPlazas}, []error{&PlazasNoPositivasError{Plazas: 0}}},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, err := NewListadoDestinos(caso.destinos)
			comprobar(t, err, errors.Join(caso.erroresEsperados...))
		})
	}
}

type casoDestinosViables struct {
	nombre           string
	notaMedia        float64
	niveles          map[string]NivelIdioma
	viablesEsperados map[int]Destino
	erroresEsperados []error
}

func nivelesQueCumplen() map[string]NivelIdioma {
	return map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelB1}
}

func probarDestinosViables(t *testing.T, casos []casoDestinosViables) {
	listado := ListadoDestinos{destinos: map[int]Destino{6: destinoValido()}}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			viables, err := listado.DestinosViables(caso.notaMedia, caso.niveles)
			comprobar(t, []any{viables, err}, []any{caso.viablesEsperados, errors.Join(caso.erroresEsperados...)})
		})
	}
}

func TestNotaMediaDelEstudiante(t *testing.T) {
	probarDestinosViables(t, []casoDestinosViables{
		{"nota media válida", 8, nivelesQueCumplen(), map[int]Destino{6: destinoValido()}, nil},
		{"nota media máxima", 10, nivelesQueCumplen(), map[int]Destino{6: destinoValido()}, nil},
		{"nota media negativa", -1, nivelesQueCumplen(), map[int]Destino{}, []error{&NotaMediaFueraDeRangoError{Nota: -1}}},
		{"nota media por encima de 10", 10.5, nivelesQueCumplen(), map[int]Destino{}, []error{&NotaMediaFueraDeRangoError{Nota: 10.5}}},
	})
}

func TestRequisitoDeIdiomas(t *testing.T) {
	probarDestinosViables(t, []casoDestinosViables{
		{"nivel superior al exigido", 8, map[string]NivelIdioma{"ALEMAN": NivelC1, "INGLES": NivelB2}, map[int]Destino{6: destinoValido()}, nil},
		{"le falta uno de los idiomas exigidos", 8, map[string]NivelIdioma{"ALEMAN": NivelB1}, map[int]Destino{}, nil},
	})
}

func TestNotaMinimaDelDestino(t *testing.T) {
	probarDestinosViables(t, []casoDestinosViables{
		{"nota media igual a la nota mínima", 7, nivelesQueCumplen(), map[int]Destino{6: destinoValido()}, nil},
		{"nota media por debajo de la nota mínima", 6.9, nivelesQueCumplen(), map[int]Destino{}, nil},
	})
}

func TestNotaPonderada(t *testing.T) {
	casos := []struct {
		nombre        string
		destino       Destino
		notaMedia     float64
		niveles       map[string]NivelIdioma
		notaPonderada float64
	}{
		{"mismo nivel que el exigido", destinoValido(), 7, nivelesQueCumplen(), 7},
		{"supera el nivel exigido en un idioma", destinoValido(), 8, map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelC1}, 9.5},
		{"supera el nivel en los dos idiomas y cuenta la mayor", destinoValido(), 8, map[string]NivelIdioma{"ALEMAN": NivelB2, "INGLES": NivelC1}, 9.5},
		{"un idioma que el destino no exige no bonifica", destinoValido(), 8, map[string]NivelIdioma{"ALEMAN": NivelB1, "INGLES": NivelB1, "FRANCES": NivelC2}, 8},
		{"destino sin idiomas exigidos y cuenta su mejor idioma", destinoCon(func(destino *Destino) { destino.idiomas = nil }), 8, map[string]NivelIdioma{"ALEMAN": NivelB2, "INGLES": NivelC1}, 9.5},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			comprobar(t, caso.destino.notaPonderada(caso.notaMedia, caso.niveles), caso.notaPonderada)
		})
	}
}

func TestNotaDeCorte(t *testing.T) {
	conAdjudicados := destinoCon(func(destino *Destino) { destino.notasTotales = []float64{9.5, 8} })
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
		{"destino nuevo sin adjudicados del año anterior", destinoValido(), 7, nivelesQueCumplen(), true},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			comprobar(t, caso.destino.esViable(caso.notaMedia, caso.niveles), caso.viable)
		})
	}
}

func TestDestinoNuevo(t *testing.T) {
	casos := []struct {
		nombre  string
		destino Destino
		nuevo   bool
	}{
		{"sin adjudicados del año anterior es nuevo", destinoValido(), true},
		{"con adjudicados del año anterior no es nuevo", destinoCon(func(destino *Destino) { destino.notasTotales = []float64{8} }), false},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			viables, _ := ListadoDestinos{destinos: map[int]Destino{6: caso.destino}}.DestinosViables(8, nivelesQueCumplen())
			comprobar(t, viables[6].EsNuevo(), caso.nuevo)
		})
	}
}

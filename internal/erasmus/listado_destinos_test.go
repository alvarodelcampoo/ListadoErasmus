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

func TestNewListadoDestinos(t *testing.T) {
	casos := []struct {
		nombre           string
		destinos         map[int]Destino
		erroresEsperados []error
	}{
		{"destino real con su nº de orden", map[int]Destino{6: destinoValido()}, nil},
		{"nº de orden mínimo permitido", map[int]Destino{1: destinoValido()}, nil},
		{"nº de orden no positivo", map[int]Destino{0: destinoValido()}, []error{&NumeroOrdenNoPositivoError{NumeroOrden: 0}}},
		{"destino con datos imposibles", map[int]Destino{6: destinoCon(func(destino *Destino) { destino.plazas = 0 })}, []error{&PlazasNoPositivasError{Plazas: 0}}},
		{"un destino imposible entre varios válidos", map[int]Destino{1: destinoValido(), 2: destinoValido(), 6: destinoCon(func(destino *Destino) { destino.plazas = 0 })}, []error{&PlazasNoPositivasError{Plazas: 0}}},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, err := NewListadoDestinos(caso.destinos)
			obtenido, esperado := fmt.Sprint(err), fmt.Sprint(errors.Join(caso.erroresEsperados...))
			if obtenido != esperado {
				t.Errorf("se esperaba %v, se obtuvo %v", esperado, obtenido)
			}
		})
	}
}

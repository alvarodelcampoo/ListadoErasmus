package erasmus

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func leerListadoReal(t *testing.T) map[int]Destino {
	fichero, _ := os.Open("../../datos/destinos.txt")
	defer fichero.Close()
	destinos, err := LeerDestinos(fichero)
	comprobar(t, err, nil)
	return destinos
}

func TestLeerDestinos(t *testing.T) {
	destinos := leerListadoReal(t)
	comprobar(t, destinos[2], Destino{
		institucion: "HES-SO HAUTE ÉCOLE SPÉCIALISÉE DE SUISSE OCCIDENTALE",
		plazas:      2,
		duracion:    9,
		cursos:      []int{3, 4},
		notaMinima:  7,
		idiomas:     map[string]NivelIdioma{"FRANCES": NivelB1},
	})
	comprobar(t, destinos[106], Destino{
		institucion: "UNIVERSIDADE DO PORTO",
		plazas:      2,
		duracion:    9,
		cursos:      []int{3, 4},
		idiomas:     map[string]NivelIdioma{"INGLES": NivelB1, "PORTUGUES": NivelB1},
	})
}

func TestLeerTodoElListado(t *testing.T) {
	destinos := leerListadoReal(t)
	_, err := NewListadoDestinos(destinos)
	comprobar(t, len(destinos), 133)
	comprobar(t, err, nil)
}

func TestLeerDestinoSinLineas(t *testing.T) {
	_, err := LeerDestinos(strings.NewReader("cabecera\n 015-texto que no es un destino\n"))
	comprobar(t, err, errors.Join(errors.Join(
		&LineaNoEncontradaError{Bloque: 1, Linea: "destino"},
		&LineaNoEncontradaError{Bloque: 1, Linea: "cursos"},
		&LineaNoEncontradaError{Bloque: 1, Linea: "idiomas"},
	)))
}

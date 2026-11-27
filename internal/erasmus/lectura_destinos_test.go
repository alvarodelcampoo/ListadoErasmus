package erasmus

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func leerListadoReal(t *testing.T) map[int]Destino {
	fichero, _ := os.Open("../../datos/destinos.txt")
	defer fichero.Close()
	destinos, err := LeerDestinos(fichero)
	assert.Equal(t, nil, err)
	return destinos
}

func TestLeerDestinos(t *testing.T) {
	destinos := leerListadoReal(t)
	assert.Equal(t, Destino{
		institucion: "HES-SO HAUTE ÉCOLE SPÉCIALISÉE DE SUISSE OCCIDENTALE",
		plazas:      2,
		duracion:    9,
		cursos:      map[Titulacion][]int{Telecomunicaciones: {3, 4}},
		notaMinima:  7,
		idiomas:     map[string]NivelIdioma{"FRANCES": NivelB1},
	}, destinos[2])
	assert.Equal(t, Destino{
		institucion: "UNIVERSIDADE DO PORTO",
		plazas:      2,
		duracion:    9,
		cursos:      map[Titulacion][]int{InformaticaYADE: {3, 4}, IngenieriaInformatica: {3, 4}, InformaticaYMatematicas: {3, 4}},
		idiomas:     map[string]NivelIdioma{"INGLES": NivelB1, "PORTUGUES": NivelB1},
	}, destinos[106])
}

func TestLeerTodoElListado(t *testing.T) {
	destinos := leerListadoReal(t)
	_, err := NewListadoDestinos(destinos, nil)
	assert.Equal(t, 133, len(destinos))
	assert.Equal(t, nil, err)
}

func TestLeerDestinoSinLineas(t *testing.T) {
	_, err := LeerDestinos(strings.NewReader("cabecera\n 015-texto que no es un destino\n"))
	assert.Equal(t, errors.Join(errors.Join(
		&LineaNoEncontradaError{Bloque: 1, Linea: "destino"},
		&LineaNoEncontradaError{Bloque: 1, Linea: "cursos"},
		&LineaNoEncontradaError{Bloque: 1, Linea: "idiomas"},
	)), err)
}

func TestLeerNotaMinimaQueNoEsNumero(t *testing.T) {
	for _, nota := range []string{"7,5x", "7,5 x"} {
		_, err := LeerDestinos(strings.NewReader("cabecera\n 015-001   B LOUVAIN01    Uclouvain    2   9   OBSERVACIONES   " + nota + "\nCursos  Titulación\n 3,4  2TC\nIdiomas\n FRANCES (B2)\n"))
		assert.Equal(t, errors.Join(errors.Join(&NotaMinimaIlegibleError{Bloque: 1, Texto: nota})), err)
	}
}

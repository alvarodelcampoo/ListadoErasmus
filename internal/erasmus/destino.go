package erasmus

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type Destino struct {
	notasTotales []float64
	idiomas      map[string]NivelIdioma
	duracion     Duracion
	institucion  string
	plazas       int
	cursos       []int
	notaMinima   float64
}

var ErrDestinoSinInstitucion = errors.New("el destino no tiene nombre de institución")

type PlazasNoPositivasError struct {
	Plazas int
}

func (e *PlazasNoPositivasError) Error() string {
	return fmt.Sprintf("las plazas son %d y tienen que ser más de 0", e.Plazas)
}

type CursoNoPositivoError struct {
	Curso int
}

func (e *CursoNoPositivoError) Error() string {
	return fmt.Sprintf("el curso es %d y tiene que ser mayor que 0", e.Curso)
}

func NuevoDestino(institucion string, plazas int, duracion Duracion, cursos []int, notaMinima float64, idiomas map[string]NivelIdioma, notasTotales []float64) (Destino, error) {
	if strings.TrimSpace(institucion) == "" {
		return Destino{}, ErrDestinoSinInstitucion
	}
	if plazas <= 0 {
		return Destino{}, &PlazasNoPositivasError{Plazas: plazas}
	}
	for _, curso := range cursos {
		if curso <= 0 {
			return Destino{}, &CursoNoPositivoError{Curso: curso}
		}
	}
	if !notaEnRango(notaMinima, notaExpedienteMaxima) {
		return Destino{}, &NotaMinimaFueraDeRangoError{Nota: notaMinima}
	}
	for i, nota := range notasTotales {
		if !notaEnRango(nota, notaConBonificacionMaxima) {
			return Destino{}, &NotaConBonificacionFueraDeRangoError{Nota: nota, Posicion: i + 1}
		}
	}
	return Destino{
		notasTotales: slices.Clone(notasTotales),
		idiomas:      maps.Clone(idiomas),
		duracion:     duracion,
		institucion:  institucion,
		plazas:       plazas,
		cursos:       slices.Clone(cursos),
		notaMinima:   notaMinima,
	}, nil
}

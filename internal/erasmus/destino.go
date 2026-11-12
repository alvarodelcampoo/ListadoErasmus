package erasmus

import (
	"errors"
	"fmt"
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

func errorSi(incumple bool, err error) error {
	if incumple {
		return err
	}
	return nil
}

func errorDeCada[T any](valores []T, comprobar func(int, T) error) error {
	errores := make([]error, len(valores))
	for i, valor := range valores {
		errores[i] = comprobar(i, valor)
	}
	return errors.Join(errores...)
}

func (d Destino) comprobar() error {
	return errors.Join(
		errorSi(strings.TrimSpace(d.institucion) == "", ErrDestinoSinInstitucion),
		errorSi(d.plazas <= 0, &PlazasNoPositivasError{Plazas: d.plazas}),
		errorSi(!notaEnRango(d.notaMinima, notaExpedienteMaxima), &NotaMinimaFueraDeRangoError{Nota: d.notaMinima}),
		errorDeCada(d.cursos, func(_ int, curso int) error {
			return errorSi(curso <= 0, &CursoNoPositivoError{Curso: curso})
		}),
		errorDeCada(d.notasTotales, func(i int, nota float64) error {
			return errorSi(!notaEnRango(nota, notaConBonificacionMaxima), &NotaConBonificacionFueraDeRangoError{Nota: nota, Posicion: i + 1})
		}),
	)
}

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

func (d Destino) cumpleIdiomas(niveles map[string]NivelIdioma) bool {
	return !slices.ContainsFunc(slices.Collect(maps.Keys(d.idiomas)), func(idioma string) bool {
		return niveles[idioma] < d.idiomas[idioma]
	})
}

func (d Destino) alcanzaNotaMinima(notaMedia float64) bool {
	return notaMedia >= d.notaMinima
}

func (d Destino) bonificaIdioma(idioma string, nivel NivelIdioma) bool {
	_, exigido := d.idiomas[idioma]
	return (exigido || len(d.idiomas) == 0) && nivel > d.idiomas[idioma]
}

func (d Destino) notaPonderada(notaMedia float64, niveles map[string]NivelIdioma) float64 {
	bonificables := maps.Clone(niveles)
	maps.DeleteFunc(bonificables, func(idioma string, nivel NivelIdioma) bool { return !d.bonificaIdioma(idioma, nivel) })
	return notaMedia + slices.Max(append(slices.Collect(maps.Values(bonificables)), 0)).bonificacion()
}

func (d Destino) EsNuevo() bool {
	return len(d.notasTotales) == 0
}

func (d Destino) superaNotaDeCorte(notaPonderada float64) bool {
	return d.EsNuevo() || slices.ContainsFunc(d.notasTotales, func(nota float64) bool { return nota <= notaPonderada })
}

func (d Destino) esViable(notaMedia float64, niveles map[string]NivelIdioma) bool {
	return d.cumpleIdiomas(niveles) && d.alcanzaNotaMinima(notaMedia) && d.superaNotaDeCorte(d.notaPonderada(notaMedia, niveles))
}

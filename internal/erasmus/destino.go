package erasmus

import (
	"errors"
	"slices"
	"strings"
)

type Destino struct {
	notasTotales []float64
	idiomas      map[string]NivelIdioma
	meses        int
	institucion  string
	plazas       int
	cursos       []int
	notaMinima   float64
}

var (
	ErrDestinoSinInstitucion = errors.New("el destino no tiene nombre de institución")
	ErrPlazasNoPositivas     = errors.New("las plazas tienen que ser mayores que 0")
	ErrMesesNoPositivos      = errors.New("los meses tienen que ser mayores que 0")
	ErrCursoNoPositivo       = errors.New("los cursos tienen que ser mayores que 0")
)

func NuevoDestino(institucion string, plazas, meses int, cursos []int, notaMinima float64, idiomas map[string]string, notasTotales []float64) (Destino, error) {
	if strings.TrimSpace(institucion) == "" {
		return Destino{}, ErrDestinoSinInstitucion
	}
	if plazas <= 0 {
		return Destino{}, ErrPlazasNoPositivas
	}
	if meses <= 0 {
		return Destino{}, ErrMesesNoPositivos
	}
	for _, curso := range cursos {
		if curso <= 0 {
			return Destino{}, ErrCursoNoPositivo
		}
	}
	if err := comprobarNota("nota mínima", notaMinima, 10); err != nil {
		return Destino{}, err
	}
	for _, nota := range notasTotales {
		if err := comprobarNota("nota total", nota, 12); err != nil {
			return Destino{}, err
		}
	}
	niveles, err := nuevosNiveles(idiomas)
	if err != nil {
		return Destino{}, err
	}
	return Destino{
		notasTotales: slices.Clone(notasTotales),
		idiomas:      niveles,
		meses:        meses,
		institucion:  institucion,
		plazas:       plazas,
		cursos:       slices.Clone(cursos),
		notaMinima:   notaMinima,
	}, nil
}

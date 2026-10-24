package erasmus

import (
	"errors"
	"slices"
	"strings"
)

type Destino struct {
	numeroOrden  int
	notasTotales []float64
	idiomas      map[string]NivelIdioma
	meses        int
	institucion  string
	plazas       int
	cursos       []int
	notaMinima   float64
}

var (
	ErrNumeroOrdenNoPositivo = errors.New("el nº de orden tiene que ser mayor que 0")
	ErrDestinoSinInstitucion = errors.New("el destino no tiene nombre de institución")
	ErrPlazasNoPositivas     = errors.New("las plazas tienen que ser mayores que 0")
	ErrMesesNoPositivos      = errors.New("los meses tienen que ser mayores que 0")
	ErrCursoNoPositivo       = errors.New("los cursos tienen que ser mayores que 0")
)

func NuevoDestino(numeroOrden int, institucion string, plazas, meses int, cursos []int, notaMinima float64, idiomas map[string]string, notasTotales []float64) (Destino, error) {
	if numeroOrden <= 0 {
		return Destino{}, ErrNumeroOrdenNoPositivo
	}
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
	if !notaEnRango(notaMinima, notaMaxima) {
		return Destino{}, ErrNotaMinimaFueraDeRango
	}
	for _, nota := range notasTotales {
		if !notaEnRango(nota, notaTotalMaxima) {
			return Destino{}, ErrNotaTotalFueraDeRango
		}
	}
	niveles, err := nuevosNiveles(idiomas)
	if err != nil {
		return Destino{}, err
	}
	return Destino{
		numeroOrden:  numeroOrden,
		notasTotales: slices.Clone(notasTotales),
		idiomas:      niveles,
		meses:        meses,
		institucion:  institucion,
		plazas:       plazas,
		cursos:       slices.Clone(cursos),
		notaMinima:   notaMinima,
	}, nil
}

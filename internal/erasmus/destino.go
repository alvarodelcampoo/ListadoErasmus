package erasmus

import (
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

type ErrorDestinoSinInstitucion struct{}

func (ErrorDestinoSinInstitucion) Error() string {
	return "el destino no tiene nombre de institución"
}

type ErrorCantidadNoPositiva struct {
	Cantidad string
	Valor    int
}

func (e ErrorCantidadNoPositiva) Error() string {
	return fmt.Sprintf("%s es %d y tiene que ser mayor que 0", e.Cantidad, e.Valor)
}

func NuevoDestino(institucion string, plazas int, duracion Duracion, cursos []int, notaMinima float64, idiomas map[string]NivelIdioma, notasTotales []float64) (Destino, error) {
	if strings.TrimSpace(institucion) == "" {
		return Destino{}, ErrorDestinoSinInstitucion{}
	}
	if plazas <= 0 {
		return Destino{}, ErrorCantidadNoPositiva{Cantidad: "el nº de plazas", Valor: plazas}
	}
	for _, curso := range cursos {
		if curso <= 0 {
			return Destino{}, ErrorCantidadNoPositiva{Cantidad: "el curso", Valor: curso}
		}
	}
	if err := comprobarNota("mínima", notaMinima, notaExpedienteMaxima); err != nil {
		return Destino{}, err
	}
	for i, nota := range notasTotales {
		if err := comprobarNota(fmt.Sprintf("total n.º %d", i+1), nota, notaConBonificacionMaxima); err != nil {
			return Destino{}, err
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

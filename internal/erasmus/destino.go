package erasmus

import (
	"fmt"
	"slices"
	"strings"
)

type Destino struct {
	notasTotales []float64
	idiomas      map[string]nivel
	meses        int
	institucion  string
	plazas       int
	cursos       []int
	notaMinima   float64
}

type ErrorDestinoSinInstitucion struct{}

func (e ErrorDestinoSinInstitucion) Error() string {
	return "el destino no tiene nombre de institución"
}

type ErrorCantidadNoPositiva struct {
	Dato  string
	Valor int
}

func (e ErrorCantidadNoPositiva) Error() string {
	return fmt.Sprintf("%s %d no válido: tiene que ser mayor que 0", e.Dato, e.Valor)
}

func NuevoDestino(institucion string, plazas, meses int, cursos []int, notaMinima float64, idiomas map[string]string, notasTotales []float64) (Destino, error) {
	if strings.TrimSpace(institucion) == "" {
		return Destino{}, ErrorDestinoSinInstitucion{}
	}
	if plazas <= 0 {
		return Destino{}, ErrorCantidadNoPositiva{Dato: "plazas", Valor: plazas}
	}
	if meses <= 0 {
		return Destino{}, ErrorCantidadNoPositiva{Dato: "meses", Valor: meses}
	}
	for _, curso := range cursos {
		if curso <= 0 {
			return Destino{}, ErrorCantidadNoPositiva{Dato: "curso", Valor: curso}
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

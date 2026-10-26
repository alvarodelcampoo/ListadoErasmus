package erasmus

import (
	"fmt"
	"maps"
)

type Estudiante struct {
	notaMedia float64
	niveles   map[string]NivelIdioma
}

func NuevoEstudiante(notaMedia float64, niveles map[string]NivelIdioma) (Estudiante, error) {
	if !notaEnRango(notaMedia, notaMaxima) {
		return Estudiante{}, fmt.Errorf("%w: %v", ErrNotaMediaFueraDeRango, notaMedia)
	}
	return Estudiante{notaMedia: notaMedia, niveles: maps.Clone(niveles)}, nil
}

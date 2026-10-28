package erasmus

import "maps"

type Estudiante struct {
	notaMedia float64
	niveles   map[string]NivelIdioma
}

func NewEstudiante(notaMedia float64, niveles map[string]NivelIdioma) (Estudiante, error) {
	if !notaEnRango(notaMedia, notaExpedienteMaxima) {
		return Estudiante{}, &NotaMediaFueraDeRangoError{Nota: notaMedia}
	}
	return Estudiante{notaMedia: notaMedia, niveles: maps.Clone(niveles)}, nil
}

package erasmus

import "maps"

type Estudiante struct {
	notaMedia float64
	niveles   map[string]NivelIdioma
}

func NuevoEstudiante(notaMedia float64, niveles map[string]NivelIdioma) (Estudiante, error) {
	if err := comprobarNota("media", notaMedia, notaExpedienteMaxima); err != nil {
		return Estudiante{}, err
	}
	return Estudiante{notaMedia: notaMedia, niveles: maps.Clone(niveles)}, nil
}

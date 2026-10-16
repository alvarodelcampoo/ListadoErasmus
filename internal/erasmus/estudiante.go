package erasmus

type Estudiante struct {
	notaMedia float64
	niveles   map[string]nivel
}

func NuevoEstudiante(notaMedia float64, niveles map[string]string) (Estudiante, error) {
	if err := comprobarNota("nota media", notaMedia, 10); err != nil {
		return Estudiante{}, err
	}
	nivelesIdioma, err := nuevosNiveles(niveles)
	if err != nil {
		return Estudiante{}, err
	}
	return Estudiante{notaMedia: notaMedia, niveles: nivelesIdioma}, nil
}

package erasmus

type Estudiante struct {
	notaMedia float64
	niveles   map[string]NivelIdioma
}

func NuevoEstudiante(notaMedia float64, niveles map[string]string) (Estudiante, error) {
	if !notaEnRango(notaMedia, notaMaxima) {
		return Estudiante{}, ErrNotaMediaFueraDeRango
	}
	nivelesIdioma, err := nuevosNiveles(niveles)
	if err != nil {
		return Estudiante{}, err
	}
	return Estudiante{notaMedia: notaMedia, niveles: nivelesIdioma}, nil
}

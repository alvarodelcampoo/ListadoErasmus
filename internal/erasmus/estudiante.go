package erasmus

type Estudiante struct {
	notaMedia float64
	niveles   map[string]nivel
}

func NuevoEstudiante(notaMedia float64, niveles map[string]string) (Estudiante, error) {
	if err := comprobarNota("nota media", notaMedia, 10); err != nil {
		return Estudiante{}, err
	}
	nivelesIdioma := make(map[string]nivel, len(niveles))
	for idioma, texto := range niveles {
		nivel, err := nuevoNivel(texto)
		if err != nil {
			return Estudiante{}, err
		}
		nivelesIdioma[idioma] = nivel
	}
	return Estudiante{notaMedia: notaMedia, niveles: nivelesIdioma}, nil
}

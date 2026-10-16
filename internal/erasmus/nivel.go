package erasmus

import "fmt"

type nivel int

const (
	nivelB1 nivel = iota + 1
	nivelB2
	nivelC1
	nivelC2
)

type ErrorNivelIdioma struct {
	Nivel string
}

func (e ErrorNivelIdioma) Error() string {
	return fmt.Sprintf("nivel de idioma %q no válido: solo se admiten B1, B2, C1 y C2", e.Nivel)
}

func nuevoNivel(texto string) (nivel, error) {
	switch texto {
	case "B1":
		return nivelB1, nil
	case "B2":
		return nivelB2, nil
	case "C1":
		return nivelC1, nil
	case "C2":
		return nivelC2, nil
	}
	return 0, ErrorNivelIdioma{Nivel: texto}
}

package erasmus

import "fmt"

type Nivel int

const (
	B1 Nivel = iota + 1
	B2
	C1
	C2
)

type ErrorNivelIdioma struct {
	Nivel string
}

func (e ErrorNivelIdioma) Error() string {
	return fmt.Sprintf("nivel de idioma %q no válido: solo se admiten B1, B2, C1 y C2", e.Nivel)
}

func NuevoNivel(texto string) (Nivel, error) {
	switch texto {
	case "B1":
		return B1, nil
	case "B2":
		return B2, nil
	case "C1":
		return C1, nil
	case "C2":
		return C2, nil
	}
	return 0, ErrorNivelIdioma{Nivel: texto}
}

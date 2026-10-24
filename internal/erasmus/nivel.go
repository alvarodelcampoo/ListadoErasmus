package erasmus

import (
	"errors"
	"fmt"
	"strings"
)

type NivelIdioma int

const (
	NivelB1 NivelIdioma = iota + 1
	NivelB2
	NivelC1
	NivelC2
)

type ErrorNivelIdioma struct {
	Nivel string
}

func (e ErrorNivelIdioma) Error() string {
	return fmt.Sprintf("nivel de idioma %q no válido: solo se admiten B1, B2, C1 y C2", e.Nivel)
}

var ErrIdiomaSinNombre = errors.New("hay un idioma sin nombre")

func nuevoNivelIdioma(texto string) (NivelIdioma, error) {
	switch texto {
	case "B1":
		return NivelB1, nil
	case "B2":
		return NivelB2, nil
	case "C1":
		return NivelC1, nil
	case "C2":
		return NivelC2, nil
	}
	return 0, ErrorNivelIdioma{Nivel: texto}
}

func nuevosNiveles(textos map[string]string) (map[string]NivelIdioma, error) {
	niveles := make(map[string]NivelIdioma, len(textos))
	for idioma, texto := range textos {
		if strings.TrimSpace(idioma) == "" {
			return nil, ErrIdiomaSinNombre
		}
		nivelIdioma, err := nuevoNivelIdioma(texto)
		if err != nil {
			return nil, err
		}
		niveles[idioma] = nivelIdioma
	}
	return niveles, nil
}

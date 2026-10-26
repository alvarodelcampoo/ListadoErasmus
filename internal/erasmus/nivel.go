package erasmus

import (
	"errors"
	"strings"
)

type NivelIdioma int

const (
	NivelB1 NivelIdioma = iota + 1
	NivelB2
	NivelC1
	NivelC2
)

var (
	ErrNivelIdiomaNoValido = errors.New("el nivel de idioma no es B1, B2, C1 ni C2")
	ErrIdiomaSinNombre     = errors.New("hay un idioma sin nombre")
)

func nuevosNiveles(niveles map[string]NivelIdioma) (map[string]NivelIdioma, error) {
	copia := make(map[string]NivelIdioma, len(niveles))
	for idioma, nivelIdioma := range niveles {
		if strings.TrimSpace(idioma) == "" {
			return nil, ErrIdiomaSinNombre
		}
		if nivelIdioma < NivelB1 || nivelIdioma > NivelC2 {
			return nil, ErrNivelIdiomaNoValido
		}
		copia[idioma] = nivelIdioma
	}
	return copia, nil
}

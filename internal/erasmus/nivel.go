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
	return 0, ErrNivelIdiomaNoValido
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

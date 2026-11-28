package erasmus

type NivelIdioma int

const (
	NivelB1 NivelIdioma = iota + 1
	NivelB2
	NivelC1
	NivelC2
)

func (n NivelIdioma) bonificacion() float64 {
	return float64(n) * bonificacionIdiomaMaxima / float64(NivelC2)
}

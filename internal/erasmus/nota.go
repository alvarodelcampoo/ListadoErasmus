package erasmus

import "errors"

const (
	notaMaxima      = 10
	notaTotalMaxima = 12
)

var (
	ErrNotaMediaFueraDeRango  = errors.New("la nota media está fuera de rango")
	ErrNotaMinimaFueraDeRango = errors.New("la nota mínima está fuera de rango")
	ErrNotaTotalFueraDeRango  = errors.New("una nota total está fuera de rango")
)

func notaEnRango(nota, maximo float64) bool {
	return nota >= 0 && nota <= maximo
}

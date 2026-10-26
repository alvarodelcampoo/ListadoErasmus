package erasmus

import "errors"

const (
	notaExpedienteMaxima      = 10
	bonificacionIdiomaMaxima  = 2
	notaConBonificacionMaxima = notaExpedienteMaxima + bonificacionIdiomaMaxima
)

var (
	ErrNotaMediaFueraDeRango  = errors.New("la nota media está fuera de rango")
	ErrNotaMinimaFueraDeRango = errors.New("la nota mínima está fuera de rango")
	ErrNotaTotalFueraDeRango  = errors.New("una nota total está fuera de rango")
)

func notaEnRango(nota, maximo float64) bool {
	return nota >= 0 && nota <= maximo
}

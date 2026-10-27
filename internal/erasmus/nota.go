package erasmus

import "fmt"

const (
	notaExpedienteMaxima      = 10
	bonificacionIdiomaMaxima  = 2
	notaConBonificacionMaxima = notaExpedienteMaxima + bonificacionIdiomaMaxima
)

type NotaMediaFueraDeRangoError struct {
	Nota float64
}

func (e *NotaMediaFueraDeRangoError) Error() string {
	return fmt.Sprintf("la nota media es %v y tiene que estar entre 0 y %v", e.Nota, notaExpedienteMaxima)
}

type NotaMinimaFueraDeRangoError struct {
	Nota float64
}

func (e *NotaMinimaFueraDeRangoError) Error() string {
	return fmt.Sprintf("la nota mínima es %v y tiene que estar entre 0 y %v", e.Nota, notaExpedienteMaxima)
}

type NotaConBonificacionFueraDeRangoError struct {
	Nota     float64
	Posicion int
}

func (e *NotaConBonificacionFueraDeRangoError) Error() string {
	return fmt.Sprintf("la nota con bonificación de idioma n.º %d es %v y tiene que estar entre 0 y %v", e.Posicion, e.Nota, notaConBonificacionMaxima)
}

func notaEnRango(nota, maximo float64) bool {
	return nota >= 0 && nota <= maximo
}

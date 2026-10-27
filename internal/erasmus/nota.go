package erasmus

import "fmt"

const (
	notaExpedienteMaxima      = 10
	bonificacionIdiomaMaxima  = 2
	notaConBonificacionMaxima = notaExpedienteMaxima + bonificacionIdiomaMaxima
)

type ErrorNotaFueraDeRango struct {
	Nota   string
	Valor  float64
	Maximo float64
}

func (e ErrorNotaFueraDeRango) Error() string {
	return fmt.Sprintf("la nota %s es %v y tiene que estar entre 0 y %v", e.Nota, e.Valor, e.Maximo)
}

func comprobarNota(nota string, valor, maximo float64) error {
	if valor < 0 || valor > maximo {
		return ErrorNotaFueraDeRango{Nota: nota, Valor: valor, Maximo: maximo}
	}
	return nil
}

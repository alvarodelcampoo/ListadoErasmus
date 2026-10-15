package erasmus

import "fmt"

type ErrorNotaFueraDeRango struct {
	Dato   string
	Nota   float64
	Maximo float64
}

func (e ErrorNotaFueraDeRango) Error() string {
	return fmt.Sprintf("%s %v fuera de rango: tiene que estar entre 0 y %v", e.Dato, e.Nota, e.Maximo)
}

func comprobarNota(dato string, nota, maximo float64) error {
	if nota < 0 || nota > maximo {
		return ErrorNotaFueraDeRango{Dato: dato, Nota: nota, Maximo: maximo}
	}
	return nil
}

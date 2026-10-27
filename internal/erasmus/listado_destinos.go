package erasmus

import (
	"fmt"
	"maps"
)

type ListadoDestinos struct {
	destinos map[int]Destino
}

type NumeroOrdenNoPositivoError struct {
	NumeroOrden int
}

func (e *NumeroOrdenNoPositivoError) Error() string {
	return fmt.Sprintf("el nº de orden es %d y tiene que ser mayor que 0", e.NumeroOrden)
}

func NuevoListadoDestinos(destinos map[int]Destino) (ListadoDestinos, error) {
	for numeroOrden := range destinos {
		if numeroOrden <= 0 {
			return ListadoDestinos{}, &NumeroOrdenNoPositivoError{NumeroOrden: numeroOrden}
		}
	}
	return ListadoDestinos{destinos: maps.Clone(destinos)}, nil
}

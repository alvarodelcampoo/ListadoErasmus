package erasmus

import "maps"

type ListadoDestinos struct {
	destinos map[int]Destino
}

func NuevoListadoDestinos(destinos map[int]Destino) (ListadoDestinos, error) {
	for numeroOrden := range destinos {
		if numeroOrden <= 0 {
			return ListadoDestinos{}, ErrorCantidadNoPositiva{Cantidad: "el nº de orden", Valor: numeroOrden}
		}
	}
	return ListadoDestinos{destinos: maps.Clone(destinos)}, nil
}

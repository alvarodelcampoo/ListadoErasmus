package erasmus

import (
	"errors"
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

func NewListadoDestinos(destinos map[int]Destino) (ListadoDestinos, error) {
	errores := make([]error, 0, 2*len(destinos))
	for numeroOrden, destino := range destinos {
		errores = append(errores, errorSi(numeroOrden <= 0, &NumeroOrdenNoPositivoError{NumeroOrden: numeroOrden}), destino.comprobar())
	}
	return ListadoDestinos{destinos: maps.Clone(destinos)}, errors.Join(errores...)
}

func (listado ListadoDestinos) DestinosViables(notaMedia float64, niveles map[string]NivelIdioma) (map[int]Destino, error) {
	if !notaEnRango(notaMedia, notaExpedienteMaxima) {
		return nil, &NotaMediaFueraDeRangoError{Nota: notaMedia}
	}
	viables := maps.Clone(listado.destinos)
	maps.DeleteFunc(viables, func(_ int, destino Destino) bool { return !destino.esViable(notaMedia, niveles) })
	return viables, nil
}

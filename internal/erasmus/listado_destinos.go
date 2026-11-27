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

type NotasSinDestinoError struct {
	NumeroOrden int
}

func (e *NotasSinDestinoError) Error() string {
	return fmt.Sprintf("hay notas del nº de orden %d, que no está en el listado", e.NumeroOrden)
}

func NewListadoDestinos(destinos map[int]Destino, adjudicaciones ...map[string]PlazaAdjudicada) (ListadoDestinos, error) {
	notas := Adjudicaciones(adjudicaciones).notasPorDestino()
	listado := ListadoDestinos{destinos: map[int]Destino{}}
	errores := make([]error, 0, 2*len(destinos)+len(notas))
	for numeroOrden, destino := range destinos {
		destino.notasTotales = notas[numeroOrden]
		listado.destinos[numeroOrden] = destino
		errores = append(errores, errorSi(numeroOrden <= 0, &NumeroOrdenNoPositivoError{NumeroOrden: numeroOrden}), destino.comprobar())
	}
	for numeroOrden := range notas {
		_, enElListado := destinos[numeroOrden]
		errores = append(errores, errorSi(!enElListado, &NotasSinDestinoError{NumeroOrden: numeroOrden}))
	}
	return listado, errors.Join(errores...)
}

func (listado ListadoDestinos) DestinosViables(estudiante Estudiante) (map[int]Destino, error) {
	if !notaEnRango(estudiante.NotaMedia, notaExpedienteMaxima) {
		return nil, &NotaMediaFueraDeRangoError{Nota: estudiante.NotaMedia}
	}
	viables := maps.Clone(listado.destinos)
	maps.DeleteFunc(viables, func(_ int, destino Destino) bool { return !destino.esViable(estudiante.NotaMedia, estudiante.Niveles) })
	return viables, nil
}

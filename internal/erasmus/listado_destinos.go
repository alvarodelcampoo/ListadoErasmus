package erasmus

import (
	"errors"
	"fmt"
	"maps"
	"slices"
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

func NewListadoDestinos(destinos map[int]Destino, adjudicaciones Adjudicaciones) (ListadoDestinos, error) {
	notas := adjudicaciones.notasPorDestino()
	destinosConNotas := repartirNotas(destinos, notas)
	return ListadoDestinos{destinos: destinosConNotas}, errors.Join(slices.Concat(erroresDeDestinos(destinosConNotas), notasSinDestino(notas, destinos))...)
}

func repartirNotas(destinos map[int]Destino, notas map[int][]float64) map[int]Destino {
	destinosConNotas := map[int]Destino{}
	for numeroOrden, destino := range destinos {
		destino.notasTotales = notas[numeroOrden]
		destinosConNotas[numeroOrden] = destino
	}
	return destinosConNotas
}

func erroresDeDestinos(destinos map[int]Destino) []error {
	errores := []error{}
	for numeroOrden, destino := range destinos {
		errores = append(errores, errorSi(numeroOrden <= 0, &NumeroOrdenNoPositivoError{NumeroOrden: numeroOrden}), destino.comprobar())
	}
	return errores
}

func notasSinDestino(notas map[int][]float64, destinos map[int]Destino) []error {
	errores := []error{}
	for numeroOrden := range notas {
		_, enElListado := destinos[numeroOrden]
		errores = append(errores, errorSi(!enElListado, &NotasSinDestinoError{NumeroOrden: numeroOrden}))
	}
	return errores
}

func (listado ListadoDestinos) DestinosViables(estudiante Estudiante) (map[int]Destino, error) {
	if !notaEnRango(estudiante.NotaMedia, notaExpedienteMaxima) {
		return nil, &NotaMediaFueraDeRangoError{Nota: estudiante.NotaMedia}
	}
	viables := maps.Clone(listado.destinos)
	maps.DeleteFunc(viables, func(_ int, destino Destino) bool { return !destino.esViable(estudiante.NotaMedia, estudiante.Niveles) })
	return viables, nil
}

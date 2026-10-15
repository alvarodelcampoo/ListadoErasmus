package erasmus

import (
	"fmt"
	"regexp"
)

type Listado struct {
	destinos map[string]Destino
}

type ErrorNumeroOrdenMalFormado struct {
	NumeroOrden string
}

func (e ErrorNumeroOrdenMalFormado) Error() string {
	return fmt.Sprintf("nº de orden %q mal formado: tiene que ser 015-NNN", e.NumeroOrden)
}

type ErrorNumeroOrdenRepetido struct {
	NumeroOrden string
}

func (e ErrorNumeroOrdenRepetido) Error() string {
	return fmt.Sprintf("nº de orden %q repetido en el listado", e.NumeroOrden)
}

var formatoNumeroOrden = regexp.MustCompile(`^015-\d{3}$`)

func (l *Listado) Registrar(numeroOrden string, destino Destino) error {
	if !formatoNumeroOrden.MatchString(numeroOrden) {
		return ErrorNumeroOrdenMalFormado{NumeroOrden: numeroOrden}
	}
	if _, existe := l.destinos[numeroOrden]; existe {
		return ErrorNumeroOrdenRepetido{NumeroOrden: numeroOrden}
	}
	if l.destinos == nil {
		l.destinos = make(map[string]Destino)
	}
	l.destinos[numeroOrden] = destino
	return nil
}

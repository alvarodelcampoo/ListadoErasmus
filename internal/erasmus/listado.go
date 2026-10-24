package erasmus

import "errors"

type ListadoDestinos struct {
	destinosPorOrden map[int]Destino
}

var ErrNumeroOrdenRepetido = errors.New("el nº de orden ya está en el listado")

func (l *ListadoDestinos) Registrar(destino Destino) error {
	if destino.institucion == "" {
		return ErrDestinoSinInstitucion
	}
	if _, existe := l.destinosPorOrden[destino.numeroOrden]; existe {
		return ErrNumeroOrdenRepetido
	}
	if l.destinosPorOrden == nil {
		l.destinosPorOrden = make(map[int]Destino)
	}
	l.destinosPorOrden[destino.numeroOrden] = destino
	return nil
}

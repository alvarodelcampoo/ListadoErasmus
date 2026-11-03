# ListadoErasmus

## Descripción del problema


Un estudiante de la Universidad de Granada, en este caso de ingeniería informática, quiere irse de Erasmus, pero no sabe cuáles son sus destinos más probables. Tiene que ir uno a uno mirando si aquel sitio que le ha gustado no tuvo una nota de corte muy superior a la suya. Esta nota viene definida por la nota media del expediente en ese momento de la solicitud y se añade alguna bonificación en función del nivel de idioma:

	-B1: +0.5 puntos
	-B2:   +1 puntos
	-C1: +1.5 puntos
	-C2:   +2 puntos

Esta bonificación sólo se aplica si tienes más nivel de idioma del que se pide, si es obligatorio para acceder al destino un B2 de inglés y tienes ese título no se te añade ninguna bonificación, sin embargo si tienes un C1 de inglés se te suma un 1.5 puntos adicionales.
	
Hay un detalle a tener en cuenta, algunos de estos destinos exigen un nivel mínimo de idioma, por ejemplo si en un destino se pide un B2 de inglés y tienes una nota superior a la nota de corte pero tienes un B1 de inglés, no puedes optar a este destino.

Cada estudiante tiene hasta 10 opciones de destinos distintos para solicitar, el problema es que si algún estudiante pone varios destinos en la que la nota de corte es demasiado superior a la del estudiante, seguramente serán opciones que el estudiante ha perdido y le van quedando menos que quizás si sean más asequibles a su nota.

Por todas estas razones, necesito alguna forma de poder ver de forma clara y rápida cuáles son las opciones de destinos más realistas para mí, según mi nota y mi nivel de idioma sin tener que revisar cada una de las opciones disponibles.

Las fuentes de datos públicas para poder realizar esto son:

- Listado de Destinos (015 ETSI Informática y Telecomunicación): 
https://internacional.ugr.es/estudiantes/movilidad-saliente/grado-estudio/movilidad-internacional/destinos-erasmus

- Listas Seleccionados:
3ªadjudicación (015 ETSI Informática y Telecomunicación): 
https://internacional.ugr.es/estudiantes/movilidad-saliente/grado-estudio/movilidad-internacional/definitiva-3-adjudicacion-erasmus-seleccionados

2ªadjudicación (015 ETSI Informática y Telecomunicación): 
https://internacional.ugr.es/estudiantes/movilidad-saliente/grado-estudio/movilidad-internacional/definitiva-2-adjudicacion-erasmus-seleccionados

1ª adjudicación (015 ETSI Informática y Telecomunicación): 
https://internacional.ugr.es/estudiantes/movilidad-saliente/grado-estudio/movilidad-internacional/definitiva-1-adjudicacion-erasmus-seleccionados

Es necesario tener acceso a las listas de las 3 convocatorias, ya que las personas que aceptan el destino qué se les ha asignado no vuelven a aparecer en las listas posteriores, ya que no hay posibilidad de acceder a esas plazas.


## Tarjeta de rol

[Tarjeta de rol](images/foto.jpeg)

## Configuración del repositorio

[Configuración del repositorio](docs/configuracion.md)


## Historias de usuario

[Historias de usuario](docs/historias-de-usuario.md)


## Milestones

[Milestones](docs/milestones.md)


## User Journey

[User Journey](docs/user-journey.md)


## Perfiles de Usuario

[Perfiles de Usuario](docs/perfiles-usuario.md)

## Instalación de Go

Descarga Go desde [go.dev/dl](https://go.dev/dl/) e instálalo en tu usuario, sin permisos de administrador:

    curl -LO https://go.dev/dl/go<versión>.linux-amd64.tar.gz
    mkdir -p ~/.local && tar -C ~/.local -xzf go<versión>.linux-amd64.tar.gz
    echo 'export PATH="$HOME/.local/go/bin:$PATH"' >> ~/.bashrc

La versión que usa el proyecto está en el `go.mod`. Si tienes una más antigua, Go descarga la necesaria automáticamente.

La justificación de esta elección está en [docs/automatizacion.md](docs/automatizacion.md).

## Tareas

Las tareas del proyecto están en el `Makefile` y se ejecutan con `make` (si no lo tienes: `sudo apt install make`):

- `make installdeps`: instala las dependencias con Go modules.
- `make check`: comprueba solo la sintaxis de los fuentes con `gofmt -e`, sin compilarlos.
- `make test`: ejecuta los tests.

La justificación de las herramientas elegidas está en [docs/automatizacion.md](docs/automatizacion.md).

## Entidad: ListadoDestinos

`ListadoDestinos` ([internal/erasmus/listado_destinos.go](internal/erasmus/listado_destinos.go)) es el listado de destinos Erasmus de la ETSIIT tal y como lo publica la UGR. Guarda cada destino identificado por su número de orden en el listado.

Cada destino (`Destino`) tiene la institución, el número de plazas, la duración (un cuatrimestre o el curso completo), los cursos desde los que se puede pedir, los idiomas que exigen un nivel mínimo, la nota mínima si la pide y las notas totales de quienes lo consiguieron en las adjudicaciones.

En el servicio web, el listado será la base para responder a un estudiante qué destinos son viables para él: los que cumplen sus idiomas y cuya nota de corte no supera su nota con la bonificación por idioma.

Para comprobar la sintaxis de la entidad: `make check`.
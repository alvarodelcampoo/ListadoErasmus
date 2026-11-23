# Herramientas para los tests

## Criterios de búsqueda

Se han buscado herramientas para Go que funcionen con `go test`, que es la orden que ya lanza la tarea `test` del Makefile. Para cada nivel, biblioteca de aserciones y framework de tests, se ha buscado por separado, partiendo de lo que trae el propio lenguaje en su documentación oficial ([testing](https://pkg.go.dev/testing)) y añadiendo las alternativas más conocidas del apartado de testing de [Awesome Go](https://github.com/avelino/awesome-go#testing), una lista revisada de herramientas de Go.

## Biblioteca de aserciones

Con el paquete de tests de Go, `testing`, ya se puede testear: el resultado se comprueba con un `if` y el fallo se marca con `t.Errorf`. Lo que hay que decidir es si se compara así o con una biblioteca de aserciones externa.

### Criterios

1. Tiene que poder comparar lo que devuelve el código del proyecto: un Destino completo, con sus campos privados, listas y mapas, y los errores que devuelve NewListadoDestinos juntos con errors.Join.
2. Que el mensaje de fallo diga qué se esperaba y qué se obtuvo.

### Opciones

- **Sin biblioteca de aserciones, con el paquete estándar [testing](https://pkg.go.dev/testing):** se comprueba el resultado con un `if` y se marca el fallo con `t.Errorf`, indicando qué se esperaba y qué se obtuvo. Para comparar valores y errores la biblioteca estándar ya trae `reflect.DeepEqual`, `errors.Is` y `errors.As`.
- **[testify](https://github.com/stretchr/testify):** la biblioteca de aserciones más usada en Go. Funciones como `assert.Equal` o `assert.ErrorAs` dan mensajes de fallo más detallados y ahorran escribir los `if`, pero añaden varias dependencias externas al proyecto.
- **[is](https://github.com/matryer/is):** biblioteca de aserciones mínima y sin dependencias, con `is.Equal`, `is.NoErr` e `is.True`. Tiene muy pocas funciones y, para que el mensaje de fallo explique qué se comprobaba, usa el comentario escrito en la línea de la aserción.

### Decisión

No usamos ninguna biblioteca de aserciones: los resultados se comprueban con un `if` y `t.Errorf` del paquete estándar `testing`. Viene con el lenguaje, así que no añade dependencias ni tiempo de instalación en el contenedor ni en la integración continua, y es la forma idiomática de testear en Go. testify aporta sobre todo comodidad en los mensajes, porque todo lo que comprueba se puede comprobar igual con la biblioteca estándar, y no compensa añadir sus dependencias por eso. is solo ahorra el `if` y obliga a escribir comentarios en los tests para que los fallos se entienda, con lo que se perdería el código limpio.

## Framework de tests

Las comprobaciones con `if` y `t.Errorf` no agrupan nada por sí mismas, así que hace falta un framework que organice los tests. No hace falta planificar de antemano cuántos tests se ejecutan ni funciones antes o después de cada test, porque los tests deben ser independientes entre sí.

### Criterios

1. Tiene que poder agrupar las comprobaciones en subtests con un nombre descriptivo, para saber qué caso ha fallado.
2. Tiene que poder crear los objetos que comparten varios tests.
3. Que añada las menos dependencias posibles, por el mismo motivo que la biblioteca de aserciones.
4. Que funcione con la forma elegida de comprobar los resultados, `if` y `t.Errorf` del paquete estándar `testing`.

### Opciones

- **Paquete estándar [testing](https://pkg.go.dev/testing):** cada test es una función `TestXxx(t *testing.T)`. Los casos se agrupan en subtests con `t.Run("nombre", ...)`, los objetos comunes se crean con funciones auxiliares o con `TestMain`, y si hay que limpiar algo está `t.Cleanup`.
- **[Ginkgo](https://github.com/onsi/ginkgo):** framework de estilo BDD (Behavior-Driven Development), en el que los tests se escriben como frases anidadas con bloques `Describe`, `Context` e `It`, y funciones como `BeforeEach` para preparar los objetos antes de cada caso. Está pensado para usarse con la biblioteca de aserciones Gomega y tiene su propia orden, `ginkgo`, para ejecutarlos.
- **[GoConvey](https://github.com/smartystreets/goconvey):** framework de estilo BDD (Behavior-Driven Development) con bloques `Convey` anidados y sus propias aserciones (`So`), con una interfaz web que relanza los tests al guardar.

### Decisión

Usamos el paquete estándar `testing`. Cumple los cuatro criterios: agrupa en subtests con `t.Run`, los objetos comunes se crean con funciones auxiliares, no añade dependencias y es el mismo paquete con el que comprobamos los resultados. Ginkgo y GoConvey traen su propia biblioteca de aserciones, lo que no cumple el cuarto criterio, y además añaden dependencias.

## Ejecutor de tests

El ejecutor es la orden que busca los tests del repositorio, los lanza e informa de cuáles pasan y cuáles fallan.

### Criterios

1. Tiene que encontrar y ejecutar todos los tests desde la raíz del repositorio con una sola orden, para lanzarla desde la tarea `test` del Makefile.
2. No puede escribir en el directorio de los fuentes, o tiene que poder configurarse para escribir en otro sitio, porque en el contenedor del siguiente objetivo ese directorio no tendrá permiso de escritura.
3. Que funcione con el framework elegido, el paquete estándar `testing`.

### Opciones

- **[go test](https://pkg.go.dev/cmd/go#hdr-Test_packages):** la orden que trae Go. Con `go test ./...` busca las funciones `TestXxx` de todos los paquetes, las compila y las ejecuta. Lo que compila lo guarda en la caché de Go, que está fuera del repositorio y se cambia de sitio con la variable `GOCACHE`.
- **[gotestsum](https://github.com/gotestyourself/gotestsum):** lanza `go test` por debajo y cambia el formato de la salida, con resúmenes y ficheros JUnit para la integración continua. Hay que instalarlo aparte.
- **[ginkgo](https://github.com/onsi/ginkgo):** la orden del framework Ginkgo, que también lanza tests de `testing`, pero está pensada para las suites de Ginkgo y hay que instalarla aparte.

### Decisión

Usamos `go test ./...`, lanzado con `make test`. Cumple los tres criterios: ejecuta todos los tests desde la raíz, solo escribe en la caché de Go, fuera de los fuentes y configurable con `GOCACHE`, y es la orden del propio paquete `testing`. gotestsum y ginkgo acaban llamanado a `go test` y solo cambian la salida, así que no aportan nada que haga falta y habría que instalarlos.


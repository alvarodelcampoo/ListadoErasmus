# Herramientas para los tests

## Criterios de búsqueda

Las opciones se han buscado en la documentación oficial de Go ([testing](https://pkg.go.dev/testing)) y en el apartado de testing de [Awesome Go](https://github.com/avelino/awesome-go#testing).

## Biblioteca de aserciones

El paquete de tests de Go, `testing`, no trae aserciones, así que hay que decidir si los tests comprueban los resultados sin biblioteca de aserciones, con un `if` y `t.Errorf`, o con una biblioteca de aserciones externa.

## Criterios

1. Tiene que poder comprobar el tipo de error devuelto, porque también hay que probar los casos en los que falla.
2. Que añada las menos dependencias posibles, porque se instalarán en cada construcción del contenedor y en cada ejecución de la integración continua.
3. Que el mensaje de fallo diga qué se esperaba y qué se obtuvo.

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

## Ejecutor de tests (test runner)

Los tests se ejecutan con `go test ./...`, la orden que trae Go para buscar, desde la raíz del repositorio, las funciones `TestXxx` de todos los paquetes, ejecutarlas e informar de cuáles pasan y cuáles fallan. Se lanza desde la tarea `test` del Makefile con `make test`, así que no hay que instalar nada más.


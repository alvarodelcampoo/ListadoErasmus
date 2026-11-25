# Herramientas para los tests

## Criterios de búsqueda

Se han buscado herramientas para Go. Para cada nivel, biblioteca de aserciones, framework de tests y test runner, se ha buscado por separado, partiendo de lo que trae el propio lenguaje en su documentación oficial ([testing](https://pkg.go.dev/testing)) y añadiendo las alternativas más conocidas del apartado de testing de [Awesome Go](https://github.com/avelino/awesome-go#testing), una lista revisada de herramientas de Go.

## Biblioteca de aserciones

La biblioteca de aserciones es la que compara el resultado obtenido con el esperado y marca el test como fallido si no coinciden.

### Criterios

1. Lo fácil que es entender qué comprueba el test, siguiendo la S de [FIRST](https://medium.com/@tasdikrahman/f-i-r-s-t-principles-of-testing-1a497acda8d6): las líneas que ocupa la misma comprobación y las funciones propias de la herramienta que hay que conocer para leerla.
2. Lo que añade al tiempo de cada ejecución de `make test`, siguiendo la F de FIRST, que se paga cada vez que se lanzan los tests.

### Opciones

Cada opción se ha probado comparando el destino 2 leído de `datos/destinos.txt` con el esperado. Para el primer criterio se cuentan las líneas de comprobación y las funciones de la herramienta que hay que conocer para leerla. Para el segundo, el tiempo se ha medido con `time`, lanzando 30 veces seguidas `go test -count=1` con esa misma comparación y dividiendo entre 30.

![Tiempo de 30 ejecuciones de go test con cada opción](../images/Captura_aserciones.png)

- **Sin biblioteca de aserciones, con el paquete estándar [testing](https://pkg.go.dev/testing):** se compara con `reflect.DeepEqual` y el fallo se marca con `t.Errorf`, escribiendo el mensaje a mano. La comprobación ocupa 3 líneas y usa 2 funciones, `reflect.DeepEqual` y `t.Errorf`. Tarda 265 ms.
- **[go-cmp](https://github.com/google/go-cmp):** biblioteca para comparar valores. Hay que pasarle `cmp.AllowUnexported` con cada tipo que tenga campos privados, y el fallo se marca con `t.Errorf`. La comprobación ocupa 3 líneas y usa 3 funciones, `cmp.Diff`, `cmp.AllowUnexported` y `t.Errorf`. Tarda 289 ms.
- **[testify](https://github.com/stretchr/testify):** la comprobación es `assert.Equal(t, esperado, obtenido)`, que ocupa 1 línea y usa 1 función. Tarda 413 ms.
- **[is](https://github.com/matryer/is):** la comprobación es `is.New(t).Equal(obtenido, esperado)`, que ocupa 1 línea y usa 2 funciones, `is.New` y `Equal`. Tarda 255 ms.

### Decisión

Usamos testify. Es la más fácil de entender según el primer criterio: la comprobación ocupa una línea y solo hay que conocer `assert.Equal`, mientras que is necesita dos funciones, y testing y go-cmp, tres líneas. En el segundo criterio es la más lenta, 413 ms frente a los 255 ms de is, pero son unos 160 ms más por ejecución de los tests, una diferencia que no compensa tener que leer más para entender cada comprobación.

## Framework de tests

El framework de tests es el que organiza los tests: los agrupa con un nombre, para saber cuál ha fallado, y permite preparar lo que comparten varios de ellos.

### Criterios

1. Lo fácil que es entender el test, siguiendo la S de [FIRST](https://medium.com/@tasdikrahman/f-i-r-s-t-principles-of-testing-1a497acda8d6): las líneas que ocupa el mismo test y las funciones propias del framework que hay que conocer para leerlo.
2. Lo que añade al tiempo de cada ejecución de `make test`, siguiendo la F de FIRST, que se paga cada vez que se lanzan los tests.

### Opciones

- **Paquete estándar [testing](https://pkg.go.dev/testing):** cada test es una función `TestXxx(t *testing.T)`. Los casos se agrupan en subtests con `t.Run("nombre", ...)`, los objetos comunes se crean con funciones auxiliares o con `TestMain`, y si hay que limpiar algo está `t.Cleanup`.
- **[Ginkgo](https://github.com/onsi/ginkgo):** framework de estilo BDD (Behavior-Driven Development), en el que los tests se escriben como frases anidadas con bloques `Describe`, `Context` e `It`, y funciones como `BeforeEach` para preparar los objetos antes de cada caso. Está pensado para usarse con la biblioteca de aserciones Gomega y tiene su propia orden, `ginkgo`, para ejecutarlos.
- **[GoConvey](https://github.com/smartystreets/goconvey):** framework de estilo BDD (Behavior-Driven Development) con bloques `Convey` anidados y sus propias aserciones (`So`), con una interfaz web que relanza los tests al guardar.

### Decisión

Usamos el paquete estándar `testing`, que es el mismo paquete con el que comprobamos los resultados. Ginkgo y GoConvey no cumplen el criterio, porque traen su propia forma de comprobarlos: Ginkgo con Gomega y GoConvey con `So`.

## Test runner

El test runner es la orden que busca los tests del repositorio, los compila, los lanza e informa de cuáles pasan y cuáles fallan.

### Criterios

1. Tiene que ser él quien busque los tests desde la raíz del repositorio, los compile y los ejecute, y no una orden que llame a otra para alguna de esas tareas.
2. No puede escribir en el directorio de los fuentes, o tiene que poder configurarse para escribir en otro sitio, porque en el contenedor del siguiente objetivo ese directorio no tendrá permiso de escritura.

### Opciones

- **[go test](https://pkg.go.dev/cmd/go#hdr-Test_packages):** la orden que trae Go. Con `go test ./...` busca las funciones `TestXxx` de todos los paquetes, las compila y las ejecuta. Lo que compila lo guarda en la caché de Go, que está fuera del repositorio y se cambia de sitio con la variable `GOCACHE`.
- **[gotestsum](https://github.com/gotestyourself/gotestsum):** lanza `go test -json`, que es quien busca, compila y ejecuta los tests, y cambia el formato de la salida, con resúmenes y ficheros JUnit para la integración continua.
- **[ginkgo](https://github.com/onsi/ginkgo):** la orden del framework Ginkgo, pensada para sus suites. Busca los paquetes y ejecuta él mismo los tests, pero para compilarlos llama a `go test -c`, que deja el binario en el directorio del paquete salvo que se cambie con `-o`.

### Decisión

Usamos `go test ./...`, lanzado con `make test`. Busca, compila y ejecuta él mismo todos los tests desde la raíz, y solo escribe en la caché de Go, fuera de los fuentes y configurable con `GOCACHE`. gotestsum no cumple el primero, porque todo lo hace `go test -json`, y ginkgo tampoco, porque compila llamando a `go test -c`.


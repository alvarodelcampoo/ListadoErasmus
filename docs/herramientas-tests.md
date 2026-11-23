# Herramientas para los tests

## Criterios de búsqueda

Se han buscado herramientas para Go que funcionen con `go test`, que es la orden que ya lanza la tarea `test` del Makefile. Para cada nivel, biblioteca de aserciones y framework de tests, se ha buscado por separado, partiendo de lo que trae el propio lenguaje en su documentación oficial ([testing](https://pkg.go.dev/testing)) y añadiendo las alternativas más conocidas del apartado de testing de [Awesome Go](https://github.com/avelino/awesome-go#testing), una lista revisada de herramientas de Go.

## Biblioteca de aserciones

Con el paquete de tests de Go, `testing`, ya se puede testear: el resultado se comprueba con un `if` y el fallo se marca con `t.Errorf`. Lo que hay que decidir es si se compara así o con una biblioteca de aserciones externa.

### Criterios

1. Tiene que poder comparar un objeto del dominio entero con el esperado, aunque guarde sus datos en campos privados, listas y mapas.
2. Tiene que poder comprobar que errores se devuelven cuando una operación falla por varias razones a la vez.

### Opciones

Cada opción se ha probado con un test que compara un destino leído de `datos/destinos.txt` con el esperado y comprobando los errores de un listado con varios destinos mal formados.

- **Sin biblioteca de aserciones, con el paquete estándar [testing](https://pkg.go.dev/testing):** se comprueba el resultado con un `if` y se marca el fallo con `t.Errorf`. El objeto se puede comparar con `reflect.DeepEqual` o comparando lo que imprimen los dos valores con `fmt.Sprint`. Con los errores juntados por `errors.Join`, `reflect.DeepEqual` compara también cómo están agrupados por dentro, mientras que `fmt.Sprint` muestra todos los que han fallado, y `errors.As` encuentra uno concreto.
- **[go-cmp](https://github.com/google/go-cmp):** biblioteca para comparar valores, pensada como alternativa a `reflect.DeepEqual`. Al comparar el Destino con `cmp.Equal` da un panic por sus campos privados, salvo que se le pase `cmp.AllowUnexported` con cada tipo. Para buscar el error no tiene nada propio y hay que usar `errors.As`.
- **[testify](https://github.com/stretchr/testify):** la biblioteca de aserciones más usada en Go. Compara el Destino con `assert.Equal` y encuentra el error con `assert.ErrorAs`.
- **[is](https://github.com/matryer/is):** biblioteca de aserciones mínima, con `is.Equal`, `is.NoErr` e `is.True`. Compara el Destino con `is.Equal`, pero para buscar el error hay que usar `errors.As` dentro de `is.True`.

### Decisión

No usamos ninguna biblioteca de aserciones: los resultados se comprueban con un `if` y `t.Errorf` del paquete estándar `testing`,  comparando lo que imprimen los dos valores con `fmt.Sprint`. Así se cumplen los dos criterios de la misma forma: compara un objeto entero con sus campos privados y, con los errores juntados, muestra cuáles han fallado sin depender de cómo están agrupados. go-cmp no cumple el primero sin configurarlo para cada tipo con campos privados. testify e is comparan con `reflect.DeepEqual` por debajo y buscan errores con `errors.As`, así que no dan nada que no tenga ya la biblioteca estándar.

## Framework de tests

Las comprobaciones con `if` y `t.Errorf` no agrupan nada por sí mismas, así que hace falta un framework que organice los tests. No hace falta planificar de antemano cuántos tests se ejecutan ni funciones antes o después de cada test, porque los tests deben ser independientes entre sí.

### Criterios

1. Tiene que poder agrupar las comprobaciones en subtests con un nombre descriptivo, para saber qué caso ha fallado.
2. Tiene que poder crear los objetos que comparten varios tests.
3. Que funcione con la forma elegida de comprobar los resultados, `if` y `t.Errorf` del paquete estándar `testing`.

### Opciones

- **Paquete estándar [testing](https://pkg.go.dev/testing):** cada test es una función `TestXxx(t *testing.T)`. Los casos se agrupan en subtests con `t.Run("nombre", ...)`, los objetos comunes se crean con funciones auxiliares o con `TestMain`, y si hay que limpiar algo está `t.Cleanup`.
- **[Ginkgo](https://github.com/onsi/ginkgo):** framework de estilo BDD (Behavior-Driven Development), en el que los tests se escriben como frases anidadas con bloques `Describe`, `Context` e `It`, y funciones como `BeforeEach` para preparar los objetos antes de cada caso. Está pensado para usarse con la biblioteca de aserciones Gomega y tiene su propia orden, `ginkgo`, para ejecutarlos.
- **[GoConvey](https://github.com/smartystreets/goconvey):** framework de estilo BDD (Behavior-Driven Development) con bloques `Convey` anidados y sus propias aserciones (`So`), con una interfaz web que relanza los tests al guardar.

### Decisión

Usamos el paquete estándar `testing`. Cumple los tres criterios: agrupa en subtests con `t.Run`, los objetos comunes se crean con funciones auxiliares y es el mismo paquete con el que comprobamos los resultados. Ginkgo y GoConvey traen su propia biblioteca de aserciones, lo que no cumple el tercero.

## Test runner

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

Usamos `go test ./...`, lanzado con `make test`. Cumple los tres criterios: ejecuta todos los tests desde la raíz, solo escribe en la caché de Go, fuera de los fuentes y configurable con `GOCACHE`, y es la orden del propio paquete `testing`. gotestsum y ginkgo acaban llamando a `go test` y solo cambian la salida, así que no aportan nada que haga falta.


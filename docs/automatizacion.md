# Herramientas para trabajar con el proyecto

Todo lo que se elige aquí se usará también en el contenedor y en la integración continua. Para probar cada opción usamos un contenedor de podman con Ubuntu 24.04 y el usuario alumno, sin permisos de
administrador. De la instalación nos fijamos en el tiempo, porque se repetirá en cada construcción del contenedor y en cada ejecución de la integración continua. Parte de ese tiempo es de descarga y depende 
de la red, así que cuando una herramienta tarda porque tiene que compilarse lo indicamos, porque eso no depende de la red.

## Gestión de las dependencias (#26)

Go tiene un gestor de dependencias estándar, [Go modules](https://go.dev/ref/mod), que viene con el propio lenguje y usa el fichero `go.mod`. No hay alternativas que valorar, ya que las que aparecen al buscar o están abandonadas, como [dep](https://github.com/golang/dep) y [Glide](https://github.com/Masterminds/glide), o son sistemas de construcción para monorepos, no son gestores de dependencias, que también usan el `go.mod`, como Bazel. Por eso usamos Go modules.

## Gestor de tareas (#27)

Para trabajar con el código hacen falta tres tareas: instalar las dependencias, comprobar la sintaxis y pasar los tests. Ahora cada una es una orden distinta de `go` que hay que conocer.

### Criterios

- **Tiempo de instalación** La herramienta se instalará al preparar el entorno.
- **Tiempo de ejecución**  Las tareas se ejecutarán en cada cambio, así que cuanto menos tarden, mejor.
- **Puede generar un fichero a partir de otro, solo cuando este cambia.** No entra en este milestone, pero más adelante habrá que generar los datos a partir de los PDF de la UGR, y así no habría que cambiar
de herramienta más adelante.


### Opciones

- **Solo `go`:** usar directamente las órdenes de Go, sin fichero de tareas.
- **[Make](https://www.gnu.org/software/make/):** gestor de tareas genérico, con un `Makefile`.
- **[Task](https://taskfile.dev/):** gestor de tareas escrito en Go, con un `Taskfile.yml`.
- **[Just](https://github.com/casey/just):** gestor de tareas escrito en Rust, con un `justfile`.
- **[Mage](https://magefile.org/):** las tareas se escriben en Go, en un `magefile.go`.

Hay más opciones, pero se parecen a alguna de estas: goyek también define las tareas en Go, como Mage; un script de shell es como usar `go` directamente pero en un fichero; y Earthly o Bazel son sistemas de 
construcción completos, más de lo que hace falta aquí. Por eso probamos una de cada tipo.

### Pruebas

**Orden para comprobar la sintaxis.** Antes de comparar los gestores de tareas, hay que elegir qué orden ejecuta la tarea `check`, porque solo tiene que comprobar la sintaxis. Lo probamos con tres ficheros: 
uno válido pero mal formateado, uno con un error de sintaxis y uno con un error de tipos (una variable que no existe), que es un error de compilación pero no de sintaxis. `gofmt -e` solo falla con el error 
de sintaxis, e indica la línea y la columna.

![gofmt -e con los tres casos](../images/Captura_check_gofmt.png)

`go build` y `go vet` también fallan con el error de tipos, así que comprueban más que la sintaxis: `go build` compila el código y `go vet` además lo analiza. Ninguna de las tres órdenes crea ficheros en el 
directorio. Por eso la tarea `check` usará `gofmt -e`.

![go build y go vet con los tres casos](../images/Captura_check_comparacion.png)

**Solo `go`.** No hace falta instalar nada más que Go, y las tres tareas funcionan con el código en solo lectura sin crear ningún fichero: `go mod download` para las dependencias, `gofmt -e .` para la 
sintaxis y `go test ./...` para los tests. El problema es que son tres órdenes sueltas que cada persona tiene que conocer, y `gofmt` ni siquiera es una orden de `go`, sino otro programa. No hay ningún 
fichero de tareas que las reúna.

Para generar ficheros, Go tiene `go generate`, pero no comprueba si el fichero de origen ha cambiado: lo ejecutamos dos veces seguidas sin tocar `entrada.txt` y las dos veces vuelve a generar `salida.txt`.

![go generate regenera siempre](../images/Captura_sologo2.png)

**Make.** No hace falta instalarlo donde se van a ejecutar las tareas: viene en la imagen oficial de Go (GNU Make 4.4.1), que es la que usaríamos en el contenedor, y en la máquina Ubuntu 24.04 de la 
integración continua de GitHub, según su [lista oficial de software](https://github.com/actions/runner-images/blob/main/images/ubuntu/Ubuntu2404-Readme.md). En cambio, no viene ni en un Ubuntu mínimo ni en
Ubuntu Desktop, según la [lista de paquetes de Ubuntu 24.04 Desktop](https://releases.ubuntu.com/24.04/ubuntu-24.04.5.1-desktop-amd64.manifest). En esas máquinas hay que instalarlo con apt, que necesita 
permisos de administrador.

![Make en Ubuntu limpio, en la imagen de Go y en la máquina de desarrollo](../images/Captura_make1.png)

Las tres tareas quedan en un `Makefile` (`make installdeps`, `make check` y `make test`) y funcionan con el código en solo lectura, sin crear ningún fichero.

![Las tres tareas con Make en solo lectura](../images/Captura_make2.png)

Además, puede generar un fichero a partir de otro solo cuando hace falta: la segunda vez que se ejecuta, dice que `salida.txt` ya está al día y no hace nada, y en cuanto cambia `entrada.txt` lo vuelve a 
generar.

![Make solo regenera cuando cambia el fichero de origen](../images/Captura_make3.png)

**Task.** No viene instalado. Se puede instalar con `go install`, sin permisos de administrador, pero tardó casi 7 minutos en compilarse. También se puede dejar fijado en el `go.mod` del proyecto con 
`go get -tool`, pero entonces mete todas sus dependencias: el `go.mod` pasa de 3 líneas a 130 y aparece un `go.sum` de 252 líneas, cuando el proyecto no tiene ninguna dependencia.

Las tres tareas quedan en un `Taskfile.yml` y funcionan con el código en solo lectura.

![Las tres tareas con Task en solo lectura](../images/Captura_task2.png)

También puede generar un fichero solo cuando cambia el de origen: la segunda vez dice que la tarea está al día, y al cambiar `entrada.txt` la vuelve a ejecutar. Pero para saberlo guarda sus comprobaciones en 
una carpeta `.task` dentro del proyecto, que habría que añadir al `.gitignore` y que no podría escribir con el código en solo lectura.

![Task genera solo cuando cambia, pero crea .task](../images/Captura_task3.png)

**Just.** No viene instalado y, como está escrito en Rust, no se puede instalar con Go ni fijar en el `go.mod`. Su script oficial descarga el programa ya compilado sin permisos de administrador: tardó 12 s.
Aun así, es una herramienta más que habría que instalar en cada máquina, en el contenedor y en la integración continua.

Las tres tareas quedan en un `justfile` y funcionan con el código en solo lectura, sin crear ningún fichero.

![Las tres tareas con Just en solo lectura](../images/Captura_just2.png)

No puede generar un fichero solo cuando cambia el de origen: al ejecutar la misma tarea dos veces seguidas sin cambiar `entrada.txt`, vuelve a copiarlo las dos veces.

![Just ejecuta la tarea siempre](../images/Captura_just3.png)

**Mage.** No viene instalado. Se instala con `go install`, sin permisos de administrador: tardó 15 s.

Las tareas se escriben en Go, en un `magefile.go`, que para las mismas tres órdenes necesitó 18 líneas. Pero con el código en solo lectura no funciona: las tres tareas fallan porque, para ejecutarlas, Mage 
intenta escribir un fichero (`mage_output_file.go`) dentro del proyecto.

![Mage falla con el código en solo lectura](../images/Captura_mage2.png)

Puede generar un fichero solo cuando cambia el de origen, con su librería `target`: la segunda vez dice que está al día y, al cambiar `entrada.txt`, lo vuelve a generar. Pero para usar esa librería Mage 
tiene que añadirse como dependencia del proyecto, en el `go.mod` y con un `go.sum`.

![Mage genera solo cuando cambia, pero como dependencia](../images/Captura_mage3.png)

**Tiempo de ejecución.** Para medir cuánto añade cada gestor, ejecutamos `check` y `test` dos veces sobre el propio proyecto, en un contenedor nuevo para cada uno: la primera vez con la caché de Go vacía, 
como en la integración continua, y la segunda con la caché ya creada. Sin ningún gestor, `gofmt -e .` tarda 0,007 s, y `go test ./...` tarda 3,65 s la primera vez (casi todo es Go compilando lo necesario 
para los tests) y 0,09 s la segunda.

![Tiempos sin gestor de tareas](../images/Captura_ejecucion_sologo.png)

Make y Just casi no añaden nada: la primera vez, `check` tarda 0,03 s con Make y 0,013 s con Just, y `test` 3,65 s y 3,58 s.

![Tiempos con Make](../images/Captura_ejecucion_make.png)
![Tiempos con Just](../images/Captura_ejecucion_just.png)

Task añade unos 75 ms a cada orden, porque antes de ejecutar nada lee el `Taskfile.yml`: `check` tarda 0,1 s. En `test` esa diferencia se pierde dentro de lo que tarda Go (3,54 s).

![Tiempos con Task](../images/Captura_ejecucion_task.png)

Mage es el único con una diferencia clara: con la caché vacía, el primer `mage check` tarda 3,57 s, porque antes de nada compila el `magefile.go`. Sumando las dos tareas en frío, como en la integración 
continua, Mage tarda 4,7 s y Make 3,7 s.

![Tiempos con Mage](../images/Captura_ejecucion_mage.png)

### Decisión

Elegimos **Make**.

Es la única opción que cumple los cinco criterios. Usar solo `go` no da un fichero de tareas: son tres órdenes sueltas, una de ellas de otro programa (`gofmt`), y `go generate` regenera los ficheros siempre. 
Task funciona, pero instalarlo tardó casi 7 minutos porque lo compila entero, o, si se fija en el `go.mod`, le añade 130 líneas y un `go.sum` de 252 a un proyecto sin dependencias; además, para saber qué 
regenerar crea una carpeta `.task` dentro del proyecto. Just es otra herramienta que instalar aparte y no puede generar ficheros solo cuando cambian. Mage ni siquiera funciona con el código en solo lectura,
y para generar ficheros tiene que añadirse como dependencia.

Make está en la imagen oficial de Go, así que en el contenedor y en la integración continua no hay que instalar nada, funciona con el código en solo lectura y es el único que genera un fichero a partir de 
otro solo cuando cambia sin escribir nada más en el proyecto. Su coste es que no viene en Ubuntu Desktop ni en un Ubuntu mínimo: en una máquina de desarrollo sin Make hay que instalarlo con apt, con permisos 
de administrador. Task, Just y Mage tampoco vienen en ninguna de esas máquinas, y además habría que instalarlos en el contenedor y en la integración continua.

El tiempo de ejecución no cambia la decisión: Make es, junto con Just, el que menos añade a las órdenes de Go, y Mage, que ya fallaba con el código en solo lectura, es además el más lento cuando la caché 
está vacía.

No usamos el gestor de dependencias (`go`) como gestor de tareas porque no tiene un fichero de tareas, pero Make no lo sustituye: cada tarea del `Makefile` es una orden de `go` o `gofmt`.
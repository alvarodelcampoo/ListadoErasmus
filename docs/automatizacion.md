# Herramientas para trabajar con el proyecto

## Instalación y versión de Go (#14)

Para trabajar en el proyecto hay que instalar Go y poder cambiar de versión desde la línea de órdenes. Ahora mismo el repositorio no dice qué versión usar ni cómo instalarla, así que cada persona podría compilar con una distinta.

### Criterios

- **Sin permisos de administrador.** Se tiene que poder instalar en el espacio del usuario, como acordamos en el #14. Lo mismo hará falta en el contenedor y en la integración continua.
- **Fija la versión en el repositorio.** Así todos, el contenedor y la integración continua usan la misma versión.
- **Cambia de versión con una orden.** Para probar el código con otra versión sin reinstalar nada a mano.

### Opciones

- **apt:** el paquete de Go de Ubuntu.
- **Binario oficial:** descargar Go de [go.dev/dl](https://go.dev/dl/) y usar el mecanismo de Go para cambiar de versión desde el `go.mod` ([toolchain](https://go.dev/doc/toolchain)).
- **Versiones paralelas oficiales:** [`golang.org/dl`](https://go.dev/doc/manage-install), que instala otras versiones junto a la que ya hay.
- **[mise](https://mise.jdx.dev/lang/go.html):** gestor de versiones para varios lenguajes.
- **[goenv](https://github.com/go-nv/goenv):** gestor de versiones solo para Go.

Hay más opciones, pero hacen lo mismo que alguna de estas: asdf es parecido a mise; gvm, g y gobrew son parecidos a goenv; snap es otro paquete del sistema, como apt. Nix y los contenedores de desarrollo preparan un entorno completo, que es más de lo que pide el #14. Por eso probamos una de cada tipo.

### Pruebas

Para probar cada opción usamos un entorno limpio con un usuario sin permisos de administrador. Para ello usamos un contenedor de podman con Ubuntu 24.04 y el usuario `alumno`.

**apt.** Sin permisos de administrador no se puede instalar Go con apt, como se ve en la captura. Por eso esta opción queda descartada.

![apt sin permisos de administrador](../images/Captura_apt.png)

**Binario oficial.** Sí se puede descargar e instalar sin permisos de administrador, en `~/.local/go`. La descarga es de 63,8 MB y, una vez instalado, ocupa 270 MB. Hay un paso manual: añadir `~/.local/go/bin` al `PATH` en el `~/.bashrc`.

![Instalación del binario oficial](../images/Captura_instalacion_binario.png)

La versión queda fijada en el propio `go.mod` (`go mod init` escribe `go 1.26.10`) y no hace falta ningún otro fichero. Para cambiar de versión basta una orden, `go mod edit -go=1.27.3`: la siguiente vez que se ejecuta `go`, descarga la 1.27.3 y la usa. Cada versión que se descarga así ocupa 253 MB más.

![Fijar y cambiar la versión con el binario oficial](../images/Captura_prueba_binario.png)

**Versiones paralelas oficiales (`golang.org/dl`).** Para usarlas hace falta tener Go instalado antes, así que no sirven para instalar Go, solo para añadir otras versiones. Se instalan sin permisos de administrador. Instalar la orden `go1.27.3` tardó 9 s y la versión descargada ocupa 350 MB, más que los 253 MB de la misma versión con el binario oficial.

![Instalación de golang.org/dl](../images/Captura_golang1.png)

Las dos versiones conviven: `go` sigue siendo la 1.26.10 y `go1.27.3` es la 1.27.3. Para cambiar de versión hay que usar otra orden. El problema es que no fija la versión en el repositorio: en el directorio solo está el `go.mod`, y ese `go 1.26.10` lo escribió el `go` normal. Nada en el repositorio dice qué orden tiene que usar cada persona. Sirve para probar otra versión, pero no para fijar la del proyecto.

![Uso de golang.org/dl](../images/Captura_golang2.png)

**mise.** Es un gestor de versiones para varios lenguajes, no solo para Go. Se instala sin permisos de administrador: el script deja el programa en `~/.local/bin` y tardó 13 s. El programa ocupa 169 MB antes de instalar ninguna versión de Go.

![Instalación de mise](../images/Captura_mise1.png)

Permite fijar la versión de Go, pero en su propio fichero, `mise.toml`, así que la versión quedaría escrita en dos sitios: ahí y en el `go.mod`. Cambiar de versión es una orden (`mise use go@<versión>`), pero para usarla hay que lanzar las órdenes con `mise exec` o activar mise en la terminal. La 1.26.10 ocupa 270 MB y la 1.27.3, 283 MB.

![Fijar y cambiar la versión con mise](../images/Captura_mise2.png)

Por defecto no lee la versión del `go.mod`: sin `mise.toml`, `mise current go` avisa de que no hay ninguna versión fijada. Se puede activar con `mise settings add idiomatic_version_file_enable_tools go`, pero esa opción se guarda en `~/.config/mise/config.toml`, en la configuración de cada usuario y no en el repositorio. Así que o se añade un `mise.toml` al repositorio, o cada persona, el contenedor y la CI tienen que configurarlo por su cuenta.

![mise leyendo el go.mod](../images/Captura_mise3.png)
![Dónde guarda mise la opción](../images/Captura_mise4.png)

**goenv.** Es un gestor de versiones solo para Go. Se instala sin permisos de administrador con un `git clone` en `~/.goenv`, que tardó 3,7 s y ocupa 8,1 MB. A cambio, necesita `git` y tres líneas en el `~/.bashrc`: la variable `GOENV_ROOT`, el `PATH` y `eval "$(goenv init -)"`.

![Instalación de goenv](../images/Captura_goenv1.png)

Su lista de versiones de Go está al día: incluye la 1.26.10 y la 1.27.3.

![Versiones disponibles en goenv](../images/Captura_goenv2.png)

Permite fijar la versión, pero en su propio fichero, `.go-version`, así que, igual que con mise, la versión quedaría en dos sitios. Cambiar de versión es una orden (`goenv local <versión>`) y después `go` usa esa versión directamente, sin `exec`. Antes hay que instalarla con `goenv install`. La 1.26.10 ocupa 270 MB y la 1.27.3, 283 MB.

![Fijar y cambiar la versión con goenv](../images/Captura_goenv3.png)

Por defecto no lee la versión del `go.mod`: sin `.go-version`, `goenv version` devuelve la del sistema. Con la variable `GOENV_GOMOD_VERSION_ENABLE=1` sí la lee, pero es una variable de cada usuario, no del repositorio. Igual que con mise, o se añade un `.go-version` al repositorio, o cada persona, el contenedor y la CI tienen que configurarla.

![goenv leyendo el go.mod](../images/Captura_goenv4.png)
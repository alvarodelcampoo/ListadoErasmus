# Herramientas para los tests

## Criterios de búsqueda

Las opciones se han buscado en la documentación oficial de Go ([testing](https://pkg.go.dev/testing)) y en el apartado de testing de [Awesome Go](https://github.com/avelino/awesome-go#testing).

## Criterios de aceptación

Biblioteca de aserciones: 
1. Tiene que poder comprobar el tipo de error devuelto, porque también hay que probar los casos en los que falla.
2. Que añada las menos dependecias posibles, porque se instalarán en cada construcción del contenedor y en cada ejecución de la integración continua.
3. Que el mensaje de fallo diga qué se esperaba y qué se obtuvo.

## Biblioteca de aserciones

Go no trae aserciones en su paquete de tests, así que hay que decidir si los tests comprueban los resultados con la biblioteca estándar o una biblioteca externa.

## Opciones

- **Paquete estándar [testing](https://pkg.go.dev/testing):** se comprueba el resultado con un `if` y se marca el fallo con `t.Errorf`, indicando qué se esperaba y qué se obtuvo. Para comparar valores y errores la biblioteca estándar ya trae `reflect.DeepEqual`, `errors.Is` y `errors.As`.
- **[testify](https://github.com/stretchr/testify):** la biblioteca de aserciones más usada en Go. Funciones como `assert.Equal` o `assert.ErrorAs` dan mensajes de fallo más detallados y ahorran escribir los `if`, pero añaden varias dependencias externas al proyecto.
- **[is](https://github.com/matryer/is):** biblioteca de aserciones mínima y sin dependencias, con `is.Equal`, `is.NoErr` e `is.True`. Tiene muy pocas funciones y, para que el mensaje de fallo explique qué se comprobaba, usa el comentario escrito en la línea de la aserción.

### Decisión

Usamos el paquete estándar `testing` con `if` y `t.Errorf`. Viene con el lenguaje, así que no añade dependecias ni tiempo de instalación en el contenedor ni en la integración continua, y es la forma idiomática de testear en Go. testify aporta sobre todo comodidad en los mensajes, porque todo lo que comprueba se puede comprobar igual con la biblioteca estándar, y no compensa añadir sus dependencias por eso. is solo ahorra el `if` y obliga a escribir comentarios en los tests para que los fallos se entienda, con lo que se perdería el código limpio.
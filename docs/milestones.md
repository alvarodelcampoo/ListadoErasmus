[M0] Paquete de estructuras de datos Erasmus

El objetivo es  crear un paquete de código que contenga las estructuras de datos necesarias para almacenar en memoria la información de las diferentes convocatorias Erasmus y el perfil del estudiante. (Relacionado con [HU1]).

Se debe entregar la implementación empaquetada como una librería que representa la información del dominio (estudiantes, destinos y requisitos), sin añadir todavía nada de lógica de negocio ni cálculos.

Será viable si, mediante una revisión del código, se comprueba que las estructuras programadas reflejan fielmente los conceptos extraídos de las historias de usuario, y que los tipos de datos definidos permiten almacenar de forma exacta la información que aparece en la fuente original (sin omitir conceptos clave ni inventar atributos innecesarios).

[M1] Lógica de cálculo de bonificaciones

A partir de los datos disponibles gracias al M0, el objetivo es calcular la nota final del estudiante para cada destino del listado teniendo en cuenta las bonificaciones por idioma. (Relacionado con [HU2]). 

Se debe entregar el código con la lógica de cálculo. Un conjunto de funciones que reciba los parámetros del estudiante y los destinos para aplicar las bonificaciones correspondientes. Junto al código de negocio, se deben entregar las pruebas automatizadas preparadas para ejecutarse.

Será viable si al ejecutar las pruebas automatizadas, estas demuestran que al pasarle al sistema los datos de distintos perfiles de alumnos y listas de destinos, se devuelve para cada uno la nota final calculada correctamente y se valida con exactitud si es apto o no.

[M0] Paquete de estructuras de datos Erasmus

El objetivo es crear un paquete de código que contenga las estructuras de datos necesarias para almacenar en memoria la información de las diferentes convocatorias Erasmus y el perfil del estudiante. (Relacionado con [HU1]).

Se debe entregar la implementación empaquetada como una librería que representa la información del dominio (estudiantes, destinos y requisitos), sin añadir todavía nada de lógica de negocio ni cálculos. Junto al código, se debe entregar un documento de trazabilidad. Este documento es simplemente una lista explicativa que conecte el mundo real con el código, es decir, debe indicar de qué archivo oficial (de los que se menciona como fuentes de datos públicas que aparecen en el README) sale cada dato necesario y en qué variable del código se ha guardado.

Será válido si, al leer el documento explicativo, se puede comprobar que no se ha olvidado incluir ningún dato que sea necesario para resolver la HU y que cada variable que hay en el código viene directamente de los documentos originales (es decir, no hay atributos inventados por el programador que no existan en la fuente oficial).

[M1] Lógica de cálculo de bonificaciones

A partir de los datos disponibles gracias al M0, el objetivo es calcular la nota final del estudiante para cada destino del listado teniendo en cuenta las bonificaciones por idioma. (Relacionado con [HU2]). 

Se debe entregar el código con la lógica de cálculo. Un conjunto de funciones que reciba los parámetros del estudiante y los destinos para aplicar las bonificaciones correspondientes. Junto al código de negocio, se deben entregar las pruebas automatizadas preparadas para ejecutarse.

Será viable si al ejecutar las pruebas automatizadas, estas demuestran que al pasarle al sistema los datos de distintos perfiles de alumnos y listas de destinos, se devuelve para cada uno la nota final calculada correctamente y se valida con exactitud si es apto o no.

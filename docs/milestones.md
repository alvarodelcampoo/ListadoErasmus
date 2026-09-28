[M0] Modelo inicial de datos de destinos

El objetivo es representar de forma estructurada los conceptos mínimos necesarios para abordar el problema (destinos, notas de corte y requisitos de idioma) sin incluir todavía la lógica de negocio ni los cálculos. (Relacionado con [HU1]).

Se debe entregar la implementación en código del modelo del problema, es decir, la estructura mínima necesaria para representar un destino Erasmus y sus requisitos, sin añadir todavía nada de lógica de negocio.

Será viable si para un destino, se puede consultar su información y coincide con la de la fuente original.

[M1] Módulo de cálculo de bonificaciones

A partir de los datos disponibles gracias al M0, el objetivo es calcular la nota final del estudiante para cada destino del listado teniendo en cuenta las bonificaciones por idioma. (Relacionado con [HU2]). 

Se debe entregar el código con la lógica de cálculo. Un conjunto de funciones que reciba los parámetros del estudiante y los destinos para aplicar las bonificaciones correspondientes.

Será viable si al pasarle los datos de un alumno y la lista de destinos, se devuelve para cada uno la nota final calculada y si es apto o no.

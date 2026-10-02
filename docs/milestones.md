[M0] Modelo de dominio

El objetivo es aplicar una metodología de diseño guiado por dominio (DDD) para extraer los términos que aparecen en las HU y el problema. (Relacionado con [HU1] y [HU2]).

Se debe entregar la implementación empaquetada como una librería que representa los conceptos del dominio, sin añadir todavía nada de lógica de negocio ni cálculos.

Será válido si se puede comprobar que se ha aplicado el proceso metodológico (DDD). Esto se verificará con el código empaquetado, comprobando que toda línea del código va ligada a un issue sacado de las HU, a través de commits explicando el issue que resuelve ese fragmento de código.

[M1] Lógica de negocio

A partir del modelo de dominio construido en el M0, el objetivo es comenzar a resolver los problemas del usuario implementando las primeras reglas de la lógica de negocio (como las bonificaciones por idioma). (Relacionado con [HU3]). 

Se debe entregar el código con la lógica de cálculo empaquetado como un módulo funcional. Junto al código de negocio, se deben entregar las pruebas automatizadas preparadas para ejecutarse.

Será viable si se demuestra que ha seguido una metodología basada en pruebas. El desarrollador sabrá que su trabajo es correcto y a la hora de revisarlo se podrá validar de forma objetiva y determinista, si al ejecutar las pruebas automatizadas estas pasan correctamente sin errores, demostrando que las reglas de negocio cumplen con los escenarios esperados.

[M0] Paquete de dominio Erasmus

El objetivo es aplicar una metodología de diseño guiado por dominio (DDD) sobre las historias de usuario para extraer el lenguaje ubicuo y modelar las entidades y objetos valor que componen el problema. (Relacionado con [HU1] y [HU4]).

Se debe entregar la implementación empaquetada como una librería que representa los conceptos del dominio (estudiantes, destinos y requisitos), sin añadir todavía nada de lógica de negocio ni cálculos. Para evidenciar que se ha seguido la metodología correcta, el código se acompañará de un documento de diseño que recoja las decisiones tomadas, actuando como mapa de trazabilidad entre el código y el mundo real (indicando de qué fuente oficial o HU del README sale cada dato modelado).

Será válido si se puede comprobar que se ha aplicado el proceso metodológico (DDD). Esto se verificará leyendo el documento de diseño y comprobando que no se ha olvidado incluir ningún dato que sea necesario para resolver la HU y que cada atributo del código procede estrictamente de los documentos originales (sin variables inventadas por el programador que no existan en la fuente oficial).

[M1] Lógica de negocio

A partir del modelo de dominio construido en el M0, el objetivo es comenzar a resolver los problemas del usuario implementando las primeras reglas de la lógica de negocio (como las bonificaciones por idioma). (Relacionado con [HU2]). 

Se debe entregar el código con la lógica de cálculo empaquetado como un módulo funcional. Junto al código de negocio, se deben entregar las pruebas automatizadas preparadas para ejecutarse.

Será viable si se demuestra que ha seguido una metodología basada en pruebas. El desarrollador sabrá que su trabajo es correcto y a la hora de revisarlo se podrá validar de forma objetiva y determinista, si al ejecutar las pruebas automatizadas estas pasan correctamente sin errores, demostrando que las reglas de negocio cumplen con los escenarios esperados.

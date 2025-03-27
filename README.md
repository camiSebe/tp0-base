# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°8

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

### Solución

#### Análisis previo

Para que el servidor pueda aceptar conexiones y procesar mensajes en paralelo, hay que prestar atención a las secciones críticas del código y a las funciones que no son thread-safe.

A priori, tengo las siguientes secciones críticas:

1. La lista de clientes conectados: Cada vez que un nuevo cliente llega, este es aceptado y se lo agrega a la lista de clientes conectados. Por lo tanto, alguien puede llegar a conectarse mientras estoy consultando la cantidad de conectados (ya sea por consulta o porque estoy en medio de un proceso de borrado)

2. El contador de agencias completadas: La idea es la misma que en el anterior, una agencia puede que quiera agregarse como completada, al mismo tiempo que estoy consultado si ya todas terminaron para empezar el sorteo.

Por otra parte, ya la cátedra nos avisa que las funciones `def has_won(bet: Bet) -> bool:` , `def store_bets(bets: list[Bet]) -> None:` y `def load_bets() -> list[Bet]:` no son thread-safe. Esto quiere decir que debo de encargarme desde el lado del servidor de asegurarme que dos hilos no quieran utilizar estas funciones porque pueden quedar en un estado corrompido/inválido.

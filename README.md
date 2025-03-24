# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°6

Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_).
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del _batch_ fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB.

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

#### Solución

Se modificó el script de python que generaba los de docker-compose para que extraiga los datasets y los clientes puedan tener como volumen el archivo que les corresponde.

Cada cliente se encarga de:
    1) Leer su csv y cargarlo
    2) Partirlo en batches de maximo tamaño el `batch: maxAmount` y enviarle un mensaje al server con dicho tamaño
    3) Enviar ese batch (aca adentro se repite el protocolo del ej 5)
    4) Esperar la confirmación del server de que el batch se guardó correctamente

El server por su parte, recibe primero el size del batch, luego todos los bets de dentro del batch, y repite este proceso hasta que el cliente deja de enviar batches.

#### Estado de los tests

Por ahora este codigo no pasa las pruebas `test_bet_amount_A` y `test_bet_amount_B`. Deje un [comentario en el campus](https://campusgrado.fi.uba.ar/mod/forum/discuss.php?d=29858) preguntando qué es lo que sucede, porque no encontré la falla hasta el momento.

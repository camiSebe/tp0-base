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
    2) Partirlo en batches de maximo tamaño el `batch: maxAmount` y enviarle un mensaje al server con dicho tamaño (anteponiendo el codigo de envio de un batch (2))
    3) Enviar ese batch (aca adentro se repite el protocolo del ej 5)
    4) Esperar la confirmación del server de que el batch se guardó correctamente
    5) Al terminar con todos los batches, se envia un mensaje de fin de envío de batches (3)

El server por su parte, recibe primero el message code para saber si se está enviando un bet particular (1), un batch de bets(2) o el fin de batches (3).

- En caso de ser un bet se espera el orden mencionado en el ej 5.
- En caso de ser un batch, se esperaran todos los bets de ese batch.
- En caso de ser el mensaje de fin de batch se cerrará la comunicación.

#### MaxBatchAmount

El ejercicio pide que pongamos un max batch amount que no superare los 8 kB. Para estimar este valor lo que hice fue, usar para estimar el dato de ese ejercicio (NOMBRE="Santiago Lionel" APELLIDO="Lorca" DOCUMENTO="30904465" NACIMIENTO="1999-03-17" NUMERO="7574" AGENCIA="0"):

Cálculo del tamaño de una apuesta:

- Cada apuesta contiene los siguientes datos:
  - Bet_message_code (omitido en batch) → 0 bytes
  - size_first_name → 4 bytes
  - first_name → longitud variable
  - size_last_name → 4 bytes
  - last_name → longitud variable
  - size_document → 4 bytes
  - document → longitud variable
  - size_birthdate → 4 bytes
  - birthdate (YYYY-MM-DD, siempre 10 bytes)
  - size_agency → 4 bytes
  - agency → longitud variable

Entonces, los campos de tamaño fijo suman: 4 + 4 + 4 + 4 + 10 + 4 = 30 bytes

Asumiendo como valores promedio para los campos variables:

- Nombre: Santiago Lionel → 16 bytes
- Apellido: Lorca → 5 bytes
- Documento: 30904465 → 8 bytes
- Agencia: 1 → 1 bytes

Entones, los campos de tamaño variable suman: 16 + 5 + 8 + 1 = 30

Entonces, en total, una apuesta ocuparía: 30 + 30 = 60 bytes

Si configuramos batchMaxAmount = 100 apuestas por batch: 60 bytes * 100 = 6.0 kB < 8 kB

Ahora, Si batchMaxAmount = 150: 60 bytes * 150 = 9.0 kB > Supera los 8 kB

Asi que podemos tomar un valor alrededor de 100 y no habría problema, en principio. Hay que recordar que como tenemos datos variables, si llegaramos a tener un batch con nombre y apellidos muy largos, podríamos pasarnos. Por lo tanto, yo elegi un batch de 50 para tener bastante margen de error.

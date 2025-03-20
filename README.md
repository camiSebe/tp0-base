# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°4

Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

#### Solución

Se agregaron tanto para el servidor como el cliente handlers para SIGTERM para que cierren de manera correcta y cierren sus sockets correctamente. Un ejemplo de salida ahora quedaría:

```
docker compose -f docker-compose-dev.yaml logs -f
client1  | 2025-03-20 21:50:34 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: INFO
client1  | 2025-03-20 21:50:34 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
client1  | 2025-03-20 21:50:39 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2025-03-20 21:50:34 INFO     action: accept_connections | result: in_progress
server   | 2025-03-20 21:50:34 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:34 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2025-03-20 21:50:34 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:34 INFO     action: accept_connections | result: in_progress
server   | 2025-03-20 21:50:39 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:39 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2025-03-20 21:50:39 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:39 INFO     action: accept_connections | result: in_progress
server   | 2025-03-20 21:50:44 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:44 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2025-03-20 21:50:44 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2025-03-20 21:50:44 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:44 INFO     action: accept_connections | result: in_progress
server   | 2025-03-20 21:50:49 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:49 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2025-03-20 21:50:49 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-20 21:50:49 INFO     action: accept_connections | result: in_progress
client1  | 2025-03-20 21:50:49 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
client1  | 2025-03-20 21:50:50 INFO     action: Signal | result: success
client1  | 2025-03-20 21:50:50 INFO     action: stop_client | result: success | client_id: 1
client1 exited with code 0
client1 exited with code 1
server   | 2025-03-20 21:50:50 INFO     action: stop_server | result: success
server exited with code 0
```

Notese que ahora tenemos los mensajes `action: stop_client` y `action: stop_server` que cierran de manera ordenada los recursos. Además, cuando se detecta un cierre, el cliente muestra el log `action: Signal | result: success`. Cuando el cliente salga con codigo 1 (`client1 exited with code 1`) quiere decir que se detecto el signal y se cortó la ejecución.

#### Apuntes tomados durante el ejercicio

1) El flag -t, --timeout especifica un shutdown timeout en segundos. Uso: `docker compose down -t <cant_segundos>`. (Debe estar en el archivo de config)

# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°4

Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

#### Solución

Se agregaron tanto para el servidor como el cliente handlers para SIGTERM y SIGINT para que cierren de manera correcta y cierren sus sockets correctamente.

Un ejemplo de ejecución sería primero, levantar el docker `make docker-compose-up` y correr `make docker-compose-logs` en una consola. Luego levantar otra consola y correr `docker kill --signal=SIGTERM client1` para enviarle la señal de SIGTERM al cliente (obs, en este caso se usa con 1 porque solo tenemos 1 cliente, pero si se modificara el docker-compose se podría matar cualquier cliente cambiando el numero) o bien corriendo `docker kill --signal=SIGTERM server` para matar al servidor.

En este caso la salida se vería de la siguiente forma:

```
docker compose -f docker-compose-dev.yaml logs -f
client1  | 2025-03-22 17:54:30 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: INFO
client1  | 2025-03-22 17:54:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
client1  | 2025-03-22 17:54:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
client1  | 2025-03-22 17:54:40 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2025-03-22 17:54:30 INFO     action: accept_connections | result: in_progress
server   | 2025-03-22 17:54:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2025-03-22 17:54:30 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:30 INFO     action: accept_connections | result: in_progress
server   | 2025-03-22 17:54:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2025-03-22 17:54:35 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:35 INFO     action: accept_connections | result: in_progress
server   | 2025-03-22 17:54:40 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:40 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
server   | 2025-03-22 17:54:40 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:40 INFO     action: accept_connections | result: in_progress

client1  | 2025-03-22 17:54:42 INFO     action: Signal | result: success
client1  | 2025-03-22 17:54:42 INFO     action: stop_client | result: success | client_id: 1

server   | 2025-03-22 17:54:45 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:45 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
client1  | 2025-03-22 17:54:45 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2025-03-22 17:54:45 INFO     action: close_connection | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:45 INFO     action: accept_connections | result: in_progress
server   | 2025-03-22 17:54:51 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2025-03-22 17:54:51 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
server   | 2025-03-22 17:54:51 INFO     action: close_connection | result: success | ip: 172.25.125.3
client1  | 2025-03-22 17:54:51 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2025-03-22 17:54:51 INFO     action: accept_connections | result: in_progress
client1  | 2025-03-22 17:54:56 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
server   | 2025-03-22 17:55:21 INFO     action: stop_server | result: success
server exited with code 0
```

Notese que ahora tenemos los mensajes `action: stop_client` y `action: stop_server` que cierran de manera ordenada los recursos. Además, cuando se detecta un cierre, el cliente muestra el log `action: Signal | result: success`.

#### Apuntes tomados durante el ejercicio / Lecciones aprendidas

1) El flag -t, --timeout especifica un shutdown timeout en segundos. Uso: `docker compose down -t <cant_segundos>`. (Debe estar en el archivo de config)
2) No hay que usar exit cuando se detecte un SIGTERM/SIGINT, el programa debería de continuar con su ejecución normal, y terminar ordenadamente sin necesidad de ello. (Obs: Uno de los fix que tuve que hacer fue sacarlo, porque si no me aparecian 2 mensajes de salida del cliente, uno que decia `client1 exited with code 0` seguido de otro que decia `client1 exited with code 1`)

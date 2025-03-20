# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°3

Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `

#### Solución

Se agrega el script de bash `validar-echo-server.sh` que verifica el funcionamiento del servidor. Para correrlo simplemente correr:
`./validar-echo-server.sh`

Los posibles resultados son:

- `action: test_echo_server | result: success`   -> Indicará que el servidor funciona de manera correcta.
- `action: test_echo_server | result: fail`      -> Indicará que el servidor NO funciona de manera correcta.

#### Apuntes tomados durante el ejercicio / Lecciones aprendidas

1) Para que detecte correctamente la network, hay que anteponer el `tp0_` (name) al `testing_net` (nombre de la red del Docker Compose). Cualquier otra cosa dara como resultado: `docker: Error response from daemon: network testing_net_2 not found.`.

2) Usé BusyBox porque es una imagen ligera y ya incluye netcat.

3) El operador `<<<` (que redirige la entrada estándar de un programa en bash) no funciona en este caso porque estaría haciendo la redirección en mi shell antes de correr docker run, y como nc se ejecuta dentro del conteiner, no se recibe la redirección correctamente.

4) No se puede usar el operador `[[ ... ]]` porque si no fallan las pruebas (tira error de health check).

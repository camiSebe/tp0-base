# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°2

Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).

### Solución

Se agrega el uso de volumes al generador tanto para los clientes como el servidor en el archivo `mi-generador.py`.
Para el caso del server tenemos: `./server/config.ini:/config.ini`, y para el cliente: `./client/config.yaml:/config.yaml`.

También se comentó la linea 19 del archivo de Dockerfile del cliente porque ya no es necesaria.

Obs: Como algunas pruebas necesitaban que el logger este en modo INFO y otras en modo DEBUG, también se quitó esa configuración del `mi-generador.py` para evitar conflictos. [Para esto último se tomó como referencia la consulta hecha en el campus: https://campusgrado.fi.uba.ar/mod/forum/discuss.php?d=29405]

### Apuntes tomados durante el ejercicio

![como-funciona-docker-volumes](image.png)

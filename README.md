# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1

Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc.

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Solución

Se agrega el script de bash `generar-compose.sh` cuya utilización es de la siguiente forma:

`./generar-compose.sh <nombre_del_archivo_de_salida> <cantidad_de_clientes>`

Es decir, se esperan 2 argumentos luego del llamado al script donde el primero sea el nombre del archivo de salida y el segundo la cantidad de clientes.
No se realizan validaciones, por lo que se pide que la cantidad de clientes sea un numero entero mayor o igual a 1. También que el nombre del archivo de salida no incluyaq caracteres especiales ni espacios (se pueden utilizar giones).

Este script llama a un script de python llamado `mi-generador.py` que creara el archivo de configuración para el docker compose con los argumentos del script de bash anteriormente mencionado.

### Apuntes tomados durante el ejercicio

![docker-compose-yaml-apuntes](image.png)

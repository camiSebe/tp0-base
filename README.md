# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°5

#### Variables de entorno

Se agrego un archivo `.env` con valores default:
`NOMBRE="Santiago Lionel"`
`APELLIDO="Lorca"`
`DOCUMENTO="30904465"`
`NACIMIENTO="1999-03-17"`
`NUMERO="7574"`
`AGENCIA="0"`

Además, se modificó tambien el generador de script de docker-compose para que pueda detectar el archivo .env y tambien entender las variables de entorno en caso de que se pasen.

#### Protocolo cliente-servidor

##### Cliente a Servidor

###### Envío de mensajes al Servidor

Se van a mandar 2 mensajes por cada dato que se tenga que enviar, el primero sera el size del dato encodeado como un Uint8, y el siguiente el dato en si encodeado como enteros de bytes. Elegí un Uint8 porque, a priori, no parece haber ningun campo que pueda superar los 255 bytes que es el límite del Uint8.

Una representación "gráfica" sería:

| Field Length | Field Data - First Name       |
|-------------|--------------------------------|
| (1 byte)    | "Santiago Lionel" (15 bytes)  |
| `0F`        | `53 61 6E 74 69 61 67 6F 20 4C 69 6F 6E 65 6C` |

| Field Length | Field Data - Last Name        |
|-------------|--------------------------------|
| (1 byte)    | "Lorca" (5 bytes)             |
| `05`        | `4C 6F 72 63 61` |

| Field Length | Field Data - Document         |
|-------------|--------------------------------|
| (1 byte)    | "30904465" (8 bytes)          |
| `08`        | `33 30 39 30 34 34 36 35` |

| Field Length | Field Data - Birthdate        |
|-------------|--------------------------------|
| (1 byte)    | "1999-03-17" (10 bytes)       |
| `0A`        | `31 39 39 39 2D 30 33 2D 31 37` |

| Field Length | Field Data - Number           |
|-------------|--------------------------------|
| (1 byte)    | "7574" (4 bytes)              |
| `04`        | `37 35 37 34` |

| Field Length | Field Data - Agency           |
|-------------|--------------------------------|
| (1 byte)    | "1" (1 byte)                  |
| `01`        | `31` |

De esta forma me aseguro de que el servidor primero reciba el size del dato que tiene que leer, y luego lo reciba. De esta forma puedo chequear si efectivamente me llegaron todos los datos que envie.

###### Recepción de mensajes del Servidor

##### Servidor a Cliente

Se recibe solamente 1 byte con la confirmación del servidor. En caso de recibir un `<0x00>` significará que la apuesto se guardó correctamente, y `<0x01>` significará que hubo algún fallo.

###### Recepción de mensajes de los Clientes

Se recibe siempre primero el largo del campo, y luego la información del campo. El orden esperado es el descripto en la parte de envío del servidor (Nombre - Apellido - DNI - Nacimiento - Numero - Agencia)
Obs: Se tomó la decisión de que si no se informa la agencia o el número apostado, se pondrá el campo con 0. En caso de no informar la fecha de nacimiento, el campo dira None

###### Envío de mensajes a los Clientes

Se enviara solamente 1 byte indicando si se pudo guardar correctamente la apuesta. Es decir, el mensaje sera `<0x00>` en caso de éxito, y `<0x01>` en caso de falla.

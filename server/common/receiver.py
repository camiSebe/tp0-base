import socket
import logging

SIZE_OF_UINT8 = 1

def receive_data(client_socket: socket.socket, num_bytes: int) -> bytes:
    """
    Receives exactly num_bytes bytes from the client.
    """
    data = b""
    while len(data) < num_bytes:
        chunk = client_socket.recv(num_bytes - len(data))
        if not chunk:
            logging.error("Connection closed by client")
            raise ConnectionError("Connection closed by client")
        data += chunk
    return data

def receive_len(client_socket: socket.socket) -> int:
    """
    Receives an unsigned 8-bit integer from the client.
    The integer is sent in network byte order.
    """
    length_bytes = receive_data(client_socket, SIZE_OF_UINT8)
    return int.from_bytes(length_bytes, byteorder="big")

def receive_string(client_socket: socket.socket) -> str:
    """
    Receives a string from the client preceded by its length.
    The string is encoded in UTF-8.
    """
    length = receive_len(client_socket)
    logging.debug(f"action: receive_len | result: success | length: {length}")
    data = receive_data(client_socket, length)
    logging.debug(f"action: receive_data | result: success | data: {data}")
    return data.decode("utf-8")

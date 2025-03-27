import socket
import logging

LENGTH_CONFIRMATION_MESSAGE = 1

def send(sock: socket.socket, data: bytes) -> None:
    """
    Sends the data to the client
    Handles socket errors if the connection is closed or broken.
    """
    try:
        logging.info(f"action: send | result: in_progress | data: {data}")
        sock.sendall(data)
        logging.info(f"action: send | result: in_progress")
    except (BrokenPipeError, ConnectionResetError):
        logging.error("Connection closed by client")
    except socket.error as e:
        logging.error(f"Socket error: {e}")
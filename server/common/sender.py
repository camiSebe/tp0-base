import socket
import logging

LENGTH_CONFIRMATION_MESSAGE = 1

def send(sock: socket.socket, result: int) -> None:
    """
    Sends the confirmation message to the client as a byte.
    Handles socket errors if the connection is closed or broken.
    """
    try:
        logging.debug(f"action: send | result: in_progress")
        sock.sendall(result.to_bytes(LENGTH_CONFIRMATION_MESSAGE, byteorder="big"))
        logging.debug(f"action: send | result: success")
    except (BrokenPipeError, ConnectionResetError):
        logging.error("Connection closed by client")
    except socket.error as e:
        logging.error("Socket error")
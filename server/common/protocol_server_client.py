import socket
import logging

from common.bet import Bet
from common.receiver import receive_string

SIZE_OF_UINT32 = 4

class ProtocolServer:
    def __init__(self, client_socket: socket.socket):
        self._client_socket = client_socket

    def receive_bet(self) -> Bet:
        """
        Receives the data of the bet from the client and stores it in a Bet Object
        Data is received in the following order:
        - Name
        - Last Name
        - DNI
        - Birthdate
        - Number
        """
        nombre = receive_string(self._client_socket)
        apellido = receive_string(self._client_socket)
        dni = receive_string(self._client_socket)
        nacimiento = receive_string(self._client_socket)
        numero = receive_string(self._client_socket)

        bet = Bet(nombre, apellido, dni, nacimiento, numero)
        # bet.print_bet()
        return bet
    
    def send_message(self, message: str):
        """
        Send a message to the client

        Function sends a message to the client
        """
        # TODO: Modify the send to avoid short-writes
        self._client_socket.send("{}\n".format(message).encode('utf-8'))
    
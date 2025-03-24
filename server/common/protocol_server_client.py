import socket
import logging

from common.utils import Bet
from common.receiver import receive_string

SIZE_OF_UINT32 = 4

class ProtocolServer:
    def __init__(self, client_socket: socket.socket):
        self._client_socket = client_socket

    def receive_bet(self) -> Bet:
        """
        Receives the data of the bet from the client and stores it in a Bet Object
        Data is received in the following order:
        - First Name
        - Last Name
        - Document
        - Birthdate
        - Number
        - Agency
        """
        first_name = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {first_name}")
        last_name = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {last_name}")
        document = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {document}")
        birthdate = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {birthdate}")
        number = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {number}")
        agency = receive_string(self._client_socket)
        logging.info(f"action: receive_message | result: success | message: {agency}")

        bet = Bet(agency, first_name, last_name, document, birthdate, number)
        return bet
    
    def send_message(self, message: str):
        """
        Send a message to the client

        Function sends a message to the client
        """
        # TODO: Modify the send to avoid short-writes
        self._client_socket.send("{}\n".format(message).encode('utf-8'))
    
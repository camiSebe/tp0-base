import socket
import logging

from common.utils import Bet
from common.receiver import receive_string
from common.sender import send

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
    
    def send_confirmation(self, bet_result):
        """
        Send a confirmation message to the client including the document and the number of the bet following the format:
        <0x00> if the bet was made successfully
        <0x01>" if the bet had a failure
        """
        send(self._client_socket, bet_result)
        logging.info(f"action: send_confirmation_message | result: success | message: {bet_result}")
    
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

        if not first_name or not last_name or not document or not birthdate or not number or not agency:
            return None
        
        bet = Bet(agency, first_name, last_name, document, birthdate, number)
        return bet
    
    def send_confirmation(self, bet):
        """
        Send a confirmation message to the client including the document and the number of the bet following the format:
        "<document> <number> <0x00>" if the bet was made successfully
        "<document> <number> <0x01>" if the bet had a failure
        """
        # TODO: Modify the send to avoid short-writes
        if bet:
            message = f"{bet.document} {bet.number} {chr(0)}"
        else:
            message = f"{chr(0)} {chr(0)} {chr(1)}"
        self._client_socket.send("{}\n".format(message).encode('utf-8'))
    
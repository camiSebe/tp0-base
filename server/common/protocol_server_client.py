import socket
import logging

from common.utils import Bet
from common.receiver import receive_string, receive_data
from common.sender import send

SIZE_OF_UINT8 = 1
SIZE_OF_UINT32 = 4

class ProtocolServer:
    def __init__(self, client_socket: socket.socket):
        self._client_socket = client_socket

    def receive_message_code(self) -> int:
        """
        Receives the message code from the client
        """
        message_code = receive_data(self._client_socket, SIZE_OF_UINT8)
        deserialize_message_code = int.from_bytes(message_code, byteorder="big")
        if deserialize_message_code == "":
            logging.error("action: receive_message_code | result: fail | error: Message code did not arrive correctly")
            return None
        logging.debug(f"action: receive_message_code | result: success | message_code: {deserialize_message_code}")
        return deserialize_message_code

    def receive_batch_size(self) -> int:
        """
        Receives the batch size from the client
        """
        batch_size = receive_data(self._client_socket, SIZE_OF_UINT32)
        deserialize_batch_size = int.from_bytes(batch_size, byteorder="big")
        if deserialize_batch_size == "":
            logging.error("action: receive_batch_size | result: fail | error: Batch size did not arrive correctly")
            return None
        logging.debug(f"action: receive_batch_size | result: success | batch_size: {deserialize_batch_size}")
        return deserialize_batch_size

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
        last_name = receive_string(self._client_socket)
        document = receive_string(self._client_socket)
        birthdate = receive_string(self._client_socket)
        number = receive_string(self._client_socket)
        agency = receive_string(self._client_socket)
        
        bet = Bet(agency, first_name, last_name, document, birthdate, number)
        logging.debug(f"action: receive_bet | result: success | First Name: {first_name} | Last Name: {last_name} | Document: {document} | Birthdate: {birthdate} | Number: {number} | Agency: {agency}")
        return bet
    
    def send_confirmation(self, bet_result):
        """
        Send a confirmation message to the client including the document and the number of the bet following the format:
        <0x00> if the bet was made successfully
        <0x01>" if the bet had a failure
        """
        send(self._client_socket, bet_result)
        logging.debug(f"action: send_confirmation_message | result: success | message: {bet_result}")
    
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
        self._agency = 0

    def set_agency(self, agency: int):
        self._agency = agency

    def _client_is_connected(self) -> bool:
        """
        Check if the client is connected
        """
        return self._client_socket is not None

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
    
    def send_confirmation(self, bet_result : int):
        """
        Send a confirmation message to the client including the document and the number of the bet following the format:
        <0x00> if the bet was made successfully
        <0x01>" if the bet had a failure
        """
        send(self._client_socket, bet_result.to_bytes(SIZE_OF_UINT8, byteorder="big"))
        logging.debug(f"action: send_confirmation_message | result: success | message: {bet_result}")
    

    def receive_agency_number(self) -> int:
        """
        Receives the agency number from the client
        """
        agency = receive_data(self._client_socket, SIZE_OF_UINT32)
        deserialize_agency = int.from_bytes(agency, byteorder="big")
        if deserialize_agency == "":
            logging.error("action: receive_agency_number | result: fail | error: Agency number did not arrive correctly")
            return None
        logging.info(f"action: receive_agency_number | result: success | agency: {deserialize_agency}")

        self.set_agency(deserialize_agency)
        return deserialize_agency

    def send_winners_to_agency(self, winners: list[Bet]):
        """
        Send the winners to the client
        """
        self.send_list_of_winners_size(len(winners))

        for winner in winners:
            logging.info(f"action: send_winners_to_agency | result: in_progress | documento: {winner.document} | agencia: {winner.agency}")
            send(self._client_socket, int(winner.document).to_bytes(SIZE_OF_UINT32, byteorder="big"))
            logging.info(f"action: send_winners_to_agency | result: success | documento: {winner.document} | agencia: {winner.agency}")

        
    def send_list_of_winners_size(self, size: int):
        """
        Send the size of the list of winners to the client
        """
        send(self._client_socket, size.to_bytes(SIZE_OF_UINT32, byteorder="big"))
        logging.info(f"action: send_list_of_winners_size | result: success | size: {size}")


    def receive_confirmation(self) -> int:
        """
        Receives the confirmation message from the client
        """
        confirmation = receive_data(self._client_socket, SIZE_OF_UINT8)
        deserialize_confirmation = int.from_bytes(confirmation, byteorder="big")
        if deserialize_confirmation == "":
            logging.error("action: receive_confirmation | result: fail | error: Confirmation message did not arrive correctly")
            return None
        logging.info(f"action: receive_confirmation | result: success | confirmation: {deserialize_confirmation}")
        return deserialize_confirmation
import socket
import logging

import signal

from common.protocol_server_client import ProtocolServer
from common.utils import has_won, load_bets, store_bets, Bet

BET_SUCCESS = 0
BET_FAILURE = 1

BATCH_SUCCESS = 0
BATCH_FAILURE = 1

TOTAL_AGENCYS = 5

NEW_BET_MESSAGE = 1
NEW_BATCH_MESSAGE = 2
END_OF_BATCHES = 3
GET_WINNERS = 4

SUCCESS_SENDING_WINNERS = 0
class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._was_closed = False
        self._agencys_completed = 0
        self._agencies_and_address = {}

        signal.signal(signal.SIGTERM, self.stop_server)
        signal.signal(signal.SIGINT, self.stop_server)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while not self._was_closed:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as e:
                break

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        self._log_info('accept_connections', 'in_progress')
        c, addr = self._server_socket.accept()
        self._log_info('accept_connections', 'success', addr[0])
        return c

    def __handle_client_connection(self, client_sock):
        """
        Receives a Bet from a specific client socket, stores it, sends a confirmation message back 
        and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while not self._was_closed:
                message_code = ProtocolServer(client_sock).receive_message_code()
                logging.debug(f"action receive_message_code | result: success | message_code: {message_code}")

                if message_code == NEW_BET_MESSAGE:
                    self.process_bet(client_sock)

                elif message_code == NEW_BATCH_MESSAGE:
                    self.process_batch_of_bets(client_sock)

                elif message_code == END_OF_BATCHES:
                    self._log_info('receive_message_code', 'success', msg='End of batches received')
                    self._agencys_completed += 1

                elif message_code == GET_WINNERS:
                    self._log_info('receive_message_code', 'success', msg='Get winners received')
                    self._add_to_request_list(client_sock)
                    if self._agencys_completed < TOTAL_AGENCYS:
                        break
                    else:
                        self._log_debug('waiting_for_all_end_of_batches', 'success')
                        self.process_winners(client_sock)
                    
                else:
                    self._log_error('receive_message_code', 'fail', error='Message code not recognized')

            if self._agencys_completed == TOTAL_AGENCYS:
                self._log_debug('waiting_for_all_end_of_batches', 'success')
                self.process_winners(client_sock)

        except OSError as e:
            self._log_error('receive_message', 'fail', error=e)
        finally:
            addr = client_sock.getpeername()
            self._log_info('close_connection', 'success', addr[0])
            client_sock.close()

    def stop_server(self, signum, frame):
        """
        Stop server

        Function that stops the server and closes the server socket
        """
        self._log_info('stop_server', 'success')
        self._server_socket.close()
        self._was_closed = True

    def process_batch_of_bets(self, client_sock):
        batch_size = ProtocolServer(client_sock).receive_batch_size()
        if batch_size==None:
            self._log_error('receive_batch_size', 'fail', error='Batch size did not arrive correctly')
            ProtocolServer(client_sock).send_confirmation(BATCH_FAILURE)
        else:
            addr = client_sock.getpeername()
            self._log_debug('receive_batch_size', 'success', addr[0], msg=f"batch_size: {batch_size}")
            bets = []
            bets_failed = 0
            bets_succeded = 0
            for _ in range(batch_size):
                new_bet = self.process_bet(client_sock)
                if new_bet==None:
                    bets_failed += 1
                else:
                    bets.append(new_bet)
                    bets_succeded += 1
            store_bets(bets)
            if bets_failed > 0:
                logging.debug(f"action: apuesta_recibida | result: fail | cantidad: {bets_failed}")
                ProtocolServer(client_sock).send_confirmation(BATCH_FAILURE)
            else:
                logging.debug(f"action: apuesta_recibida | result: success | cantidad: {bets_succeded}")
                ProtocolServer(client_sock).send_confirmation(BATCH_SUCCESS)


    def process_bet(self, client_sock) -> Bet:
        bet = ProtocolServer(client_sock).receive_bet()
        if bet==None: 
            self._log_error('receive_bet', 'fail', error='Bet did not arrive correctly')
            return None
        else:
            addr = client_sock.getpeername()
            self._log_debug('receive_bet', 'success', addr[0], msg=bet.log_message())
            return bet

    def process_winners(self, client_sock):
        self._log_info('sorteo', 'success')

        winners = self.get_winners()

        logging.info(f"action: ganadores_obtenidos | result: success | cantidad: {len(winners)}")
        for winner in winners:
            logging.info(f"action: ganador_obtenido | result: success | documento: {winner.document} | agencia: {winner.agency}")
        
        for agency, client_sock in self._agencies_and_address.items():
            if client_sock.fileno() == -1:
                # self._log_error('send_winners_to_agency', 'fail', error='Invalid socket descriptor for agency' + str(agency))
                continue
            winners_from_agency = [winner for winner in winners if winner.agency == agency]
            ProtocolServer(client_sock).send_winners_to_agency(winners_from_agency, agency)
            if ProtocolServer(client_sock).receive_confirmation() == SUCCESS_SENDING_WINNERS:
                self._log_debug('send_winners_to_agency', 'success', msg='Confirmation received')
            else:
                self._log_error('send_winners_to_agency', 'fail', error='Confirmation not received')

    def get_winners(self) -> list[Bet]:
        all_bets = load_bets()
        winners = []
        logging.debug(f"action: checking_winners | result: in_progress")
        for bet in all_bets:
            if has_won(bet):
                winners.append(bet)
        logging.debug(f"action: checking_winners | result: success | winners: {len(winners)}")
        return winners

    def _add_to_request_list(self, client_sock):
        agency = ProtocolServer(client_sock).receive_agency_number()
        self._agencies_and_address[agency] = client_sock

    ### Logging helper functions
    def _log_info(self, action, result, ip=None, msg=None):
        """
        Helper function for logging informational messages
        """
        log_message = f"action: {action} | result: {result}"
        if ip:
            log_message += f" | ip: {ip}"
        if msg:
            log_message += f" | msg: {msg}"
        logging.info(log_message)

        ### Logging helper functions
    def _log_debug(self, action, result, ip=None, msg=None, doc=None, number=None):
        """
        Helper function for logging informational messages
        """
        log_message = f"action: {action} | result: {result}"
        if ip:
            log_message += f" | ip: {ip}"
        if msg:
            log_message += f" | msg: {msg}"
        if doc:
            log_message += f" | dni: {doc}"
        if number:
            log_message += f" | numero: {number}"
        logging.debug(log_message)

    def _log_error(self, action, result, error=None):
        """
        Helper function for logging error messages
        """
        log_message = f"action: {action} | result: {result}"
        if error:
            log_message += f" | error: {error}"
        logging.error(log_message)
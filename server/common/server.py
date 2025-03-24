import socket
import logging

import signal

from common.protocol_server_client import ProtocolServer
from common.utils import store_bets, Bet

BET_SUCCESS = 0
BET_FAILURE = 1

BATCH_SUCCESS = 0
BATCH_FAILURE = 1

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._was_closed = False
        signal.signal(signal.SIGTERM, self.stop_server)

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
        self._log_debug('accept_connections', 'in_progress')
        c, addr = self._server_socket.accept()
        self._log_debug('accept_connections', 'success', addr[0])
        return c

    def __handle_client_connection(self, client_sock):
        """
        Receives a Bet from a specific client socket, stores it, sends a confirmation message back 
        and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            self.process_batch_of_bets(client_sock)
        except OSError as e:
            self._log_error('receive_message', 'fail', error=e)
        finally:
            addr = client_sock.getpeername()
            self._log_debug('close_connection', 'success', addr[0])
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
            for _ in range(batch_size):
                new_bet = self.process_bet(client_sock)
                if new_bet==None:
                    bets_failed += 1
                else:
                    bets.append(new_bet)
            store_bets(bets)
            if bets_failed > 0:
                logging.info(f"action: apuesta_recibida | result: fail | cantidad: {bets_failed}")
                ProtocolServer(client_sock).send_confirmation(BATCH_FAILURE)
            else:
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
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
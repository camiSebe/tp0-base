import socket
import logging

import signal
import threading

from concurrent.futures import ThreadPoolExecutor

from common.protocol_server_client import ProtocolServer
from common.utils import has_won, load_bets, store_bets, Bet

BET_SUCCESS = 0
BET_FAILURE = 1

BATCH_SUCCESS = 0
BATCH_FAILURE = 1

NEW_BET_MESSAGE = 1
NEW_BATCH_MESSAGE = 2
# END_OF_BATCHES = 3
GET_WINNERS = 4

SUCCESS_SENDING_WINNERS = 0

class Server:
    def __init__(self, port, listen_backlog, clients_count):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._was_closed = False
        self._clients_conected = []
        self._agencys_completed = 0
        self._total_agencies = clients_count

        self._clients_conected_lock = threading.Lock()
        self._agencys_completed_lock = threading.Lock()
        self._storage_lock = threading.Lock()

        signal.signal(signal.SIGTERM, self.stop_server)
        signal.signal(signal.SIGINT, self.stop_server)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        with ThreadPoolExecutor() as executor:
            while not self._was_closed:
                try:
                    client_sock = self.__accept_new_connection()
                    executor.submit(self.__handle_client_connection, client_sock)
                    # self.__handle_client_connection(client_sock)
                except OSError as e:
                    if self._was_closed:
                        break
                    self._log_error('accept_connections', 'fail', error=e)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        self._log_info('accept_connections', 'in_progress')
        c, addr = self._server_socket.accept()

        with self._clients_conected_lock:
            self._log_info('CLIENTS CONNECTED LOCK in __accept_new_connection', 'in_progress', msg='Adding new client')
            self._clients_conected.append(c)
        self._log_info('CLIENTS CONNECTED LOCK in __accept_new_connection', 'success')
        
        self._log_info('accept_connections', 'success', addr[0])
        return c

    def __handle_client_connection(self, client_sock: socket.socket):
        """
        Receives a Bet from a specific client socket, stores it, sends a confirmation message back 
        and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol_server = ProtocolServer(client_sock)

            while protocol_server._client_is_connected():
                message_code = protocol_server.receive_message_code()
                logging.debug(f"action receive_message_code | result: success | message_code: {message_code}")

                if message_code == NEW_BET_MESSAGE:
                    self.process_bet(protocol_server)

                elif message_code == NEW_BATCH_MESSAGE:
                    self.process_batch_of_bets(protocol_server)

                elif message_code == GET_WINNERS:
                    with self._agencys_completed_lock:
                        self._log_info('AGENCIES COMPLETED LOCK in message_code == GET_WINNERS', 'in_progress', msg='Adding 1 to agencies_completed')
                        self._agencys_completed += 1
                    self._log_info('AGENCIES COMPLETED LOCK in message_code == GET_WINNERS', 'success')
                    # self._log_info('receive_message_code', 'success', msg='Get winners received')
                    
                    self._log_info('waiting_for_all_end_of_batches', 'in_progress', msg=f"agencies_completed: {self._agencys_completed} | total_agencies: {self._total_agencies}")

                    with self._agencys_completed_lock:
                        self._log_info('AGENCIES COMPLETED LOCK in message_code == GET_WINNERS', 'in_progress', msg='Checking if all agencies have completed')
                        if self._agencys_completed == self._total_agencies:
                            self._log_debug('waiting_for_all_end_of_batches', 'success')
                            self.process_winners()
                        self._log_info('AGENCIES COMPLETED LOCK in message_code == GET_WINNERS', 'success')
                    break
                    
                else:
                    self._log_error('receive_message_code', 'fail', error='Message code not recognized')

        except OSError as e:
            self._log_error('receive_message', 'fail', error=e)

        except Exception as e:
            self._log_error('receive_message', 'fail', error=e)

    def stop_server(self, signum, frame):
        """
        Stop server

        Function that stops the server and closes the server socket
        """
        self._log_info('stop_server', 'success')
        self._server_socket.close()
        self._was_closed = True

        with self._clients_conected_lock:
            self._log_info('CLIENTS CONNECTED LOCK in stop_server', 'in_progress', msg='Closing all clients connections')
            self._close_all_clients_connection()
            self._log_info('close_all_clients', 'success')
        self._log_info('CLIENTS CONNECTED LOCK in stop_server', 'success')

    def process_batch_of_bets(self, protocol_server: ProtocolServer):
        batch_size = protocol_server.receive_batch_size()
        if batch_size==None:
            self._log_error('receive_batch_size', 'fail', error='Batch size did not arrive correctly')
            protocol_server.send_confirmation(BATCH_FAILURE)
        else:
            self._log_debug('receive_batch_size', 'success',msg=f"batch_size: {batch_size}")
            bets = []
            bets_failed = 0
            bets_succeded = 0
            for _ in range(batch_size):
                new_bet = self.process_bet(protocol_server)
                if new_bet==None:
                    bets_failed += 1
                else:
                    bets.append(new_bet)
                    bets_succeded += 1
            
            # store_bets(bets)
            with self._storage_lock:
                self._log_debug('STORAGE LOCK in process_batch_of_bets', 'in_progress', msg=f"Storing {len(bets)} bets")
                store_bets(bets)
            self._log_debug('STORAGE LOCK in process_batch_of_bets', 'success')

            if bets_failed > 0:
                # logging.debug(f"action: apuesta_recibida | result: fail | cantidad: {bets_failed}")
                protocol_server.send_confirmation(BATCH_FAILURE)
            else:
                logging.debug(f"action: apuesta_recibida | result: success | cantidad: {bets_succeded}")
                protocol_server.send_confirmation(BATCH_SUCCESS)


    def process_bet(self, protocol_server: ProtocolServer) -> Bet:
        bet = protocol_server.receive_bet()
        if bet==None: 
            self._log_error('receive_bet', 'fail', error='Bet did not arrive correctly')
            return None
        else:
            # self._log_debug('receive_bet', 'success', msg=bet.log_message())
            return bet

    def process_winners(self):
        """
        Process all the winners in the current bets and send them to the corresponding agencies
        """
        try:
            self._log_info('sorteo', 'success')

            winners = self.get_winners()

            logging.info(f"action: ganadores_obtenidos | result: success | cantidad: {len(winners)}")
            for winner in winners:
                logging.info(f"action: ganador_obtenido | result: success | documento: {winner.document} | agencia: {winner.agency}")

            for client in self._clients_conected:
                protocol_server = ProtocolServer(client)
                this_agency = protocol_server.receive_agency_number()
                winners_from_this_agency = [winner for winner in winners if winner.agency == this_agency]
                protocol_server.send_winners_to_agency(winners_from_this_agency)
        
        finally:
            with self._clients_conected_lock:
                self._log_info('CLIENTS CONNECTED LOCK in process_winners', 'in_progress', msg='Closing all clients connections')
                self._close_all_clients_connection()
            self._log_info('CLIENTS CONNECTED LOCK in process_winners', 'success')


    def get_winners(self) -> list[Bet]:
        """
        Returns a list of all the winners in the current bets
        """
        # all_bets = load_bets()
        with self._storage_lock:
            self._log_info('STORAGE LOCK in get_winners for load_bets', 'in_progress', msg='Loading all bets')
            all_bets = load_bets()
        self._log_info('STORAGE LOCK in get_winners for load_bets', 'success')

        winners = []
        logging.debug(f"action: checking_winners | result: in_progress")
        for bet in all_bets:
            if has_won(bet):
                winners.append(bet)
        logging.debug(f"action: checking_winners | result: success | winners: {len(winners)}")
        return winners

    def _close_all_clients_connection(self):
        """
        Closes all client connections in the server
        """
        self._log_info('close_all_clients', 'in_progress')
        for client in self._clients_conected:
            client.close()
        self._clients_conected.clear()
        self._log_info('close_all_clients', 'success')
 
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
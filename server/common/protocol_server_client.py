import socket

class ProtocolServer:
    def __init__(self, client_socket: socket.socket):
        self._client_socket = client_socket

    def receive_message(self) -> str:
        """
        Receive a message from the client

        Function blocks until a message from the client is received
        """
        # TODO: Modify the receive to avoid short-reads
        return self._client_socket.recv(1024).rstrip().decode('utf-8')
    
    def send_message(self, message: str):
        """
        Send a message to the client

        Function sends a message to the client
        """
        # TODO: Modify the send to avoid short-writes
        self._client_socket.send("{}\n".format(message).encode('utf-8'))
    
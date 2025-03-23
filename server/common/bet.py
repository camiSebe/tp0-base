
class Bet:
    def __init__(self, nombre, apellido, dni, nacimiento, numero):
        self.nombre = nombre
        self.apellido = apellido
        self.dni = dni
        self.nacimiento = nacimiento
        self.numero = numero

    def print_bet(self):
        print(f"Nombre: {self.nombre}")
        print(f"Apellido: {self.apellido}")
        print(f"DNI: {self.dni}")
        print(f"Nacimiento: {self.nacimiento}")
        print(f"Numero: {self.numero}")
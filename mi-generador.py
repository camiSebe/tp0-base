import sys
import yaml

def generar_docker_compose(archivo_salida, cant_clientes):
    docker_compose_data = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1",
                    "LOGGING_LEVEL=DEBUG",
                ],
                "networks": ["testing_net"],
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [
                        {"subnet": "172.25.125.0/24"}
                    ]
                }
            }
        }
    }

    for i in range(1, cant_clientes+1):
        docker_compose_data["services"][f"client{i}"] = {
            "container_name": f"client{i}",
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
                "CLI_LOG_LEVEL=DEBUG"
            ],
            "networks": ["testing_net"],
            "depends_on": ["server"]
        }
    
    with open(archivo_salida, "w") as file:
        yaml.dump(docker_compose_data, file, sort_keys=False)



if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Uso: python3 mi-generador.py <archivo_salida> <cant_clientes>")
        sys.exit(1)

    archivo_salida = sys.argv[1]
    cant_clientes = int(sys.argv[2])

    generar_docker_compose(archivo_salida, cant_clientes)
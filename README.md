# 2024-M08-P1

## Objetivo

Ferramenta de linha de comando (CLI) para monitorar a temperatura em uma cadeia de supermercados.


### Construtor

Visando uma mutabilidade fácil, foi criado a coleta de dados por meio de um json, onde se passa o atributo e o valor maximo ou minimo que o sensor pode coletar.

```json
	[{
    "id": "lj01f01",
    "tipo": "freezer",
    "temperatura": -18,
    "timestamp": "01/03/2024 14:30:21"
  },
  {
    "id": "lj02g03",
    "tipo": "geladeira",
    "temperatura": 5,
    "timestamp": "01/03/2024 14:30:40"
  },
  {
    "id": "lj01f02",
    "tipo": "freezer",
    "temperatura": -26,
    "timestamp": "01/03/2024 14:31:01"
  },
  {
    "id": "lj03g01",
    "tipo": "geladeira",
    "temperatura": 12,
    "timestamp": "01/03/2024 14:30"
  }
]
```
Onde no código, é feito a leitura do json, verificação dos dados e transformação dos dados para o formato de envio.

## Como Rodar

### Script

1. Após rodar o mosquito em sua maquina local, com o comando:

```bash

    mosquitto -c mosquitto.conf

``` 
2. Abra um novo terminal e acesse o diretório do projeto, com o seguinte comando:

```bash

    cd /src

```
3. Agora, rode o seguinte comando:
    
```bash
    
    chmod +x start.sh
    
```
4. Rode o seguinte comando para instalar as dependências do projeto e executar o simulador:

```bash

    ./start.sh

```

5. O simulador irá rodar e enviará mensagens para o tópico `sensor` a cada 1 segundos.

### Testes

```bash

    mosquitto -c mosquitto.conf

``` 
2. Abra um novo terminal e acesse o diretório do projeto, com o seguinte comando:

```bash

    cd /src

```
3. Agora, rode o seguinte comando:
    
```bash
    
    chmod +x test.sh
    
```
4. Rode o seguinte comando para instalar as dependências do projeto e executar o simulador:

```bash

    ./test.sh
```

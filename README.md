# 🧭 OnFly API

API REST para gerenciamento da aplicação **OnFly**, desenvolvida em **Go** com o framework **Fiber**, utilizando arquitetura limpa, Swagger para documentação, e recursos como PostgreSQL, Redis e RabbitMQ — todos orquestrados com Docker Compose.

---

## 🚀 Como rodar a aplicação com Docker

Certifique-se de ter o **Docker** e o **Docker Compose** instalados em sua máquina.

1. Clone este repositório:
   ```bash
   git clone https://github.com/Natannms/onfly-teste
   cd onfly-api

2. Suba os containers:

```bash
 docker compose up -d
```
3. Acesse a aplicação no navegador:

🌐 API base: http://localhost:3000

📄 Swagger (documentação): http://localhost:3000/swagger/index.html


🧱 Serviços disponíveis
| Serviço        | Porta Local | Descrição                       |
| -------------- | ----------- | ------------------------------- |
| API (Go/Fiber) | `3000`      | Servidor principal da aplicação |
| PostgreSQL     | `5432`      | Banco de dados relacional       |
| Redis          | `6379`      | Cache e filas                   |
| RedisInsight   | `5540`      | Interface gráfica para o Redis  |
| RabbitMQ       | `15672`     | Painel web de gerenciamento     |
| RabbitMQ AMQP  | `5672`      | Protocolo de fila (interno)     |

🧪 Variáveis de Ambiente
As variáveis são automaticamente lidas do .env ou configuradas diretamente no docker-compose.yml (modo desenvolvimento).
Para ambientes de produção, configure adequadamente:

DATABASE_URL

RABBITMQ_URL

QUEUE_REDIS_HOST, PORT, PASSWORD, etc.


📁 Estrutura dos Containers
app: Código da API (Go + Fiber) com hot reload via Air

postgres: Banco de dados relacional

redis: Cache

redisinsight: Interface visual para Redis

rabbitmq: Fila de mensagens com dashboard


🧹 Encerrando os serviços
Para parar os serviços e remover os containers:
```bash
docker compose down
```


✨ Contribuições
Pull requests são bem-vindos! Sinta-se à vontade para abrir uma issue ou sugerir melhorias.
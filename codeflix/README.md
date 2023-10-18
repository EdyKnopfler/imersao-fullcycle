# Aplicação CodeFlix

Todo o gerenciamento do projeto Go é feito dentro do contêiner `app` (ver arquivo `docker-compose.yaml`). A pasta do projeto com os fontes é mapeada em um volume.

## Domain

Coração da aplicação, sem pensar em banco de dados ou infraestrutura de qualquer tipo.

**Entidade é lugar de regra de negócio.** Os nomes de métodos seguem os nomes dados pelo negócio às operações.

O domínio declara interfaces mas não implementa. Ex.: repositórios.

## gRPC

Streaming bidirecional usando HTTP/2.

**RPC**: _Remote Procedure Call_

Recomendado no backend, entre microsserviços. Entre o browser e o servidor ainda não está maduro.

Comunicação baseada em um contrato de dados entre as partes (arquivos `.proto`).

### HTTP/2

Usa a mesma conexão multiplexada (envia e recebe). Dados são binários.

### Protocol Buffers (ProtoBuf)

Serialização de dados estruturados em formato binário (em oposição ao formato texto de JSON, XML, etc.).

### API Unary

Request + Response.

### Server Streaming

Servidor continua enviando informação ao longo do tempo. Não é preciso esperar 100% da resposta estar pronta.

### Client Streaming

O inverso do anterior :)

### Bidirectional Streaming

:D




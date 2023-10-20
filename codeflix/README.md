# Aplicação CodeFlix

Todo o gerenciamento do projeto Go é feito dentro do contêiner `app` (ver arquivo `docker-compose.yaml`). A pasta do projeto com os fontes é mapeada em um volume.

*Obs.:** a imersão ensina o uso da ferramenta [cobra](https://github.com/spf13/cobra) para gerenciar múltiplos comandos via flags no mesmo binário; aqui preferi separar em vários binários para tornar o código mais enxuto para consulta.

## Domain

Coração da aplicação, sem pensar em banco de dados ou infraestrutura de qualquer tipo.

**Entidade é lugar de regra de negócio.** Os nomes de métodos seguem os nomes dados pelo negócio às operações.

O domínio declara interfaces mas não implementa. Ex.: repositórios.

## Application

Camada de infraestrutura: serviço de gRPC, etc.

## gRPC

Streaming bidirecional usando HTTP/2.

**RPC**: _Remote Procedure Call_

Recomendado no backend, entre microsserviços. Entre o browser e o servidor ainda não está maduro.

Comunicação baseada em um contrato de dados entre as partes (arquivos `.proto`).

### HTTP/2

Usa a mesma conexão multiplexada (envia e recebe). Dados são binários.

### API Unary

Request + Response.

### Server Streaming

Servidor continua enviando informação ao longo do tempo. Não é preciso esperar 100% da resposta estar pronta.

### Client Streaming

O inverso do anterior :)

### Bidirectional Streaming

:D


### Protocol Buffers (ProtoBuf)

Serialização de dados estruturados em formato binário (em oposição ao formato texto de JSON, XML, etc.). Ver pasta `application/grpc/protofiles`.

Compilando os `.proto` ():

```bash
protoc \
  --go_out=application/grpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=application/grpc/pb \
  --go-grpc_opt=paths=source_relative \
  --proto_path=application/grpc/protofiles \
  application/grpc/protofiles/*.proto
```

* `protoc` instalado na imagem pelo `Dockerfile`
* saída em `application/grpc/pb` conforme descrito nos arquivos

Subindo o servidor:

```bash
docker compose up [-d] app db
docker compose exec app go run cmd/server/main.go
```

Testando com o Evans (já instalado na imagem):

```bash
docker compose exec -it app evans -r repl

> package derso.com.imersao.fullcycle.codeflix
> service PixService
> call RegisterPixKey
```

## gORM

Trabalhar com ORM traz ônus e bônus...

Um ônus do gORM é a exigência das colunas de ID nas chaves estrangeiras:

```go
type Account struct {
	BankID string `gorm:"column:bank_id;type:uuid;not null" valid:"-"`
}
```

## Kafka

Subindo a aplicação e o broker Kafka, criando os tópicos e verificando:

```bash
  docker compose up [-d] db app kafka
  docker compose run --rm kafka-topics-generator
  docker compose exec kafka kafka-topics --list --bootstrap-server=localhost:9092
```

Também já está configurado o centro de controle da Confluent (http://localhost:9021/):

```bash
docker compose up [-d] control-center  # ou `run --rm` se quiser ser econômico :)
```

Para os próximos testes criamos o tópico _teste_; no `docker-compose.yml` há amostras de comandos (também podemos usar o centro de controle para criar o tópico e consumir):

Emitindo mensagem:

```bash
docker compose exec app go run cmd/kafka/main.go
```

Consumindo

```bash
docker compose exec kafka kafka-console-consumer --topic=teste --bootstrap-server=localhost:9092
```
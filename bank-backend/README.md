# Aplicação Back-end da Instituição Bancária

Simula um banco com suas contas, e dispara transações "Pix" fictícias através do serviço `codeflix`.

Utiliza o framework [NestJS](https://nestjs.com/).

Não estou tratando todos os casos e erros aqui, apenas codando o suficiente para pegar a ideia e a prática (descartando tudo que é repetitivo).

## Docker

Tudo dentro de contêineres!

* Removemos as pastas `dist` e `node_modules` e fazemos a instalação lá por dentro.

## Instalação das dependências

```bash
$ npm install
```

## Rodando

```bash
# development
$ npm run start

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## Testes

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

## TypeORM

```bash
npm install typeorm @nestjs/typeorm
```

## gRPC

* Copiamos os `*.proto`
* Configuramos como assets no `nest-cli.json`
* Dependências:

```bash
npm install --save @grpc/proto-loader @grpc/grpc-js @nestjs/microservices
```
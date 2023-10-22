# Imersão Full Cycle

https://imersao.fullcycle.com.br/aulas/

Estudo de caso ("Pix" fictício) envolvendo questões arquiteturais.

Como é uma referência de consulta, existem vários `TODOs` com lembretes para completar a ideia que está sendo passada.

Mais detalhes no `README.md` de cada subprojeto!

Tecnologias:
* Docker
* Golang
* gRPC
* Kafka
* Nest.js
* Next.js
* Kubernetes

## `codeflix`

Aplicação central do backend em linguagem Golang. Comunica-se com outros serviços usando gRPC.

## `bank-backend`

Serviço do banco individual em NestJS. Faz comunicação com o `codeflix` via gRPC.

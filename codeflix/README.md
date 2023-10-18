# Anotações

## Aplicação

Todo o gerenciamento do projeto Go é feito dentro do contêiner `app` (ver arquivo `docker-compose.yaml`). A pasta do projeto com os fontes é mapeada em um volume.

## Domain

Coração da aplicação, sem pensar em banco de dados ou infraestrutura de qualquer tipo.

**Entidade é lugar de regra de negócio.** Os nomes de métodos seguem os nomes dados pelo negócio às operações.

Haja código repetido...

O domínio declara interfaces mas não implementa. Ex.: repositórios.


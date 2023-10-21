# Debian
FROM node:20-slim

# procps = habilita o reloading
RUN apt update -y && \
    apt install procps -y && \
    npm install -g @nestjs/cli@10.0.0

# Usuário node da imagem
WORKDIR /home/node/app

USER node

CMD [ "tail", "-f", "/dev/null" ]
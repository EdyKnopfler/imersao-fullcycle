# Debian
FROM node:20-slim

# procps = habilita o reloading
RUN apt update -y && \
    apt install procps && \
    npm install -g @nestjs/cli@10.0.0

WORKDIR /home/node/app  # Usuário node da imagem

USER node

CMD [ "tail", "-f", "/dev/null" ]
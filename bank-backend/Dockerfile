# Debian
FROM node:20-slim

# procps = habilita o reloading
RUN apt update -y && \
    apt install procps -y && \
    npm install -g @nestjs/cli@10.0.0

WORKDIR /home/node/app

CMD [ "tail", "-f", "/dev/null" ]
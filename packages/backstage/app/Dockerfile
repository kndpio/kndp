FROM nginx:mainline
LABEL org.opencontainers.image.source https://github.com/kndpio/kndp
RUN apt-get update && apt-get -y install jq && rm -rf /var/lib/apt/lists/*

COPY packages/backstage/app/dist /usr/share/nginx/html
COPY packages/backstage/app/conf.template /etc/nginx/templates/default.conf.template
COPY packages/backstage/app/inject-config.sh /docker-entrypoint.d/40-inject-config.sh
RUN chmod 0755 /docker-entrypoint.d/40-inject-config.sh

ENV PORT 3000
# Stage 1: Build Stage
FROM node:18-alpine as build

WORKDIR /packages
COPY . /packages
RUN apk update && \
    apk add --no-cache autoconf automake libtool nasm build-base
RUN npm install
RUN npx nx run kndp:build-extensions
RUN npx nx run kndp:build-backend
RUN npx nx run kndp:build-list

# Stage 2: Production Stage
FROM node:18-alpine

WORKDIR /packages

COPY --from=build /packages/packages/kndp/charts/kndp/files /packages
EXPOSE 3000
CMD ["node", "/packages/backend.js"]

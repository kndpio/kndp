FROM node:18-alpine
COPY ./charts/kndp/files /packages
WORKDIR /packages/backend
RUN npm ci 
CMD ["node", "/packages/backend/backend.js"]

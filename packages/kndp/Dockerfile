FROM node:18-alpine
COPY ./charts/kndp/files /packages
WORKDIR /packages/backend
RUN npm install 
CMD ["node", "/packages/backend/backend.js"]

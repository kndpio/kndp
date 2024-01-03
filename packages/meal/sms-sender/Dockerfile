FROM node:lts-alpine
WORKDIR /home/node
COPY . .
RUN yarn install
CMD ["node", "./sms-send.js"]
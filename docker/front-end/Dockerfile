FROM node:13.12.0-alpine

WORKDIR /srv

COPY front_end/package.json ./
COPY front_end/package-lock.json ./
RUN npm install

COPY front_end ./
RUN npm run build

# start app
CMD ["npm", "start"]

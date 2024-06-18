FROM node:20 AS builder

WORKDIR /srv

COPY front_end/package.json ./
COPY front_end/package-lock.json ./
RUN npm ci

COPY front_end ./
RUN npm run build

ENTRYPOINT ["npm"]
CMD ["run", "start"]

# FROM nginx:stable

# COPY --from=builder /srv/build/ /usr/share/nginx/html

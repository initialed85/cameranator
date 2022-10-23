FROM node:13.12.0-alpine

WORKDIR /srv

COPY front_end/package.json ./
COPY front_end/package-lock.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm npm set cache /usr/src/app/.npm && npm install

COPY front_end ./
RUN --mount=type=cache,target=/usr/src/app/.npm npm set cache /usr/src/app/.npm && npm run build

# start app
CMD ["npm", "start"]

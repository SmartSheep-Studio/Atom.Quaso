# Run image with this command
# docker run --rm --name quaso --net host -v $(pwd)/config.toml:/app/config.toml  -v $(pwd)/resources:/resources quaso

# Building Frontend
FROM node:18-alpine as quaso-web
WORKDIR /source
COPY . .
WORKDIR /source/packages/quaso-web
RUN rm -rf dist node_modules
RUN --mount=type=cache,target=/source/packages/quaso-web/node_modules,id=quaso_web_modules_cache,sharing=locked \
    --mount=type=cache,target=/root/.npm,id=quaso_web_node_cache \
    yarn install
RUN --mount=type=cache,target=/source/packages/quaso-web/node_modules,id=quaso_web_modules_cache,sharing=locked \
    yarn run build-only
RUN mv /source/packages/quaso-web/dist /dist

# Building Backend
FROM golang:alpine as quaso-server

WORKDIR /source
COPY . .
COPY --from=quaso-web /dist /source/packages/quaso-web/dist
RUN mkdir /dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/server ./pkg/cmd/server/main.go

# Runtime
FROM golang:alpine

COPY --from=quaso-server /dist/server /app/server

EXPOSE 9443

CMD ["/app/server", "serve"]
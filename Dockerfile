# build stage
FROM golang:1.18.1-alpine AS build-env
RUN apk --no-cache add build-base git curl
ADD . /app
RUN cd /app/cmd/server && go build -o bin

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/cmd/server/bin /app/bin

COPY --from=build-env /app/cmd/server/cert /app/cert
COPY --from=build-env /app/cmd/server/oauth /app/oauth
COPY --from=build-env /app/cmd/server/config.toml /app/config.toml

COPY --from=build-env /app/web /app/web

ENV HOST_STATIC="web/static/"
ENV HOST_TEMPLATES="web/templates/"

ENTRYPOINT "./bin" -env $ARGS
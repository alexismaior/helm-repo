# syntax=docker/dockerfile:experimental
FROM siop/golang as stage1

# Compile Delve
RUN apk add --no-cache git
RUN go get github.com/go-delve/delve/cmd/dlv

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN --mount=type=ssh go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build /build.sh

FROM alpine:3.7 as debug

# Port 8080 belongs to our application, 40000 belongs to Delve
EXPOSE 40000

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

WORKDIR /

LABEL siop.descricao="Serviço do qualitativo"
LABEL traefik.http.routers.qualitativo-servico.rule="PathPrefix(`/modulo/qualitativo/api`)"

COPY --from=stage1 /app/app-debug /
COPY --from=stage1 /go/bin/dlv /
COPY --from=stage1 /app/internal/graphql/schema.graphql /internal/graphql/schema.graphql

# Run delve
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/app-debug"]

FROM alpine:3.7 as with-shell

LABEL siop.descricao="Serviço do qualitativo"
LABEL traefik.http.routers.qualitativo-servico.rule="PathPrefix(`/modulo/qualitativo/api`)"

COPY --from=stage1 /app/app /
COPY --from=stage1 /app/internal/graphql/schema.graphql /internal/graphql/schema.graphql

CMD [ "/app" ]

FROM scratch

LABEL siop.descricao="Serviço do qualitativo"
LABEL traefik.http.routers.qualitativo-servico.rule="PathPrefix(`/modulo/qualitativo/api`)"

COPY --from=stage1 /app/app /
COPY --from=stage1 /app/internal/graphql/schema.graphql /internal/graphql/schema.graphql

CMD [ "/app" ]

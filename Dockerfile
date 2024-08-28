# FROM golang:1.22-alpine3.18 AS build
# WORKDIR /go/src/proglog
# COPY . .
# RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog
# # Definir la versión del probe y descargarla usando wget
# RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
#     wget -qO /go/bin/grpc_health_probe \
#     https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#     chmod +x /go/bin/grpc_health_probe

# FROM scratch
# COPY --from=build /go/bin/proglog /bin/proglog
# # END: beginning
# # START_HIGHLIGHT
# COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
# # END_HIGHLIGHT
# # verify that the binary is in the path

# RUN ls -l /bin
# ENV PATH="/bin"
# ENTRYPOINT ["/bin/proglog"]

FROM golang:1.22-alpine3.18 AS build
WORKDIR /go/src/proglog
COPY . .
RUN apk add --no-cache bash wget
RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog

ENTRYPOINT ["/go/bin/proglog"]
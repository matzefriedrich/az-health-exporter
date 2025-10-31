FROM golang:1.25 AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags "-s -w" -o /out/az-health-exporter ./cmd/az-health-exporter


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/az-health-exporter /usr/local/bin/az-health-exporter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/az-health-exporter", "monitor", "--p", "8080"]

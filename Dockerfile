FROM golang:1.25 AS builder

ARG TARGETOS
ARG TARGETARCH

ARG APP_RELEASE
ARG APP_RELEASE_DATE
ARG APP_VERSION
ARG CI_COMMIT_SHORT_SHA

WORKDIR /src

COPY . .

RUN go mod download

ENV PACKAGE_NAME="github.com/matzefriedrich/az-health-exporter"
ENV LDFLAGS="-X ${PACKAGE_NAME}/internal.CommitSha=${CI_COMMIT_SHORT_SHA} -X ${PACKAGE_NAME}/internal.Version=${APP_VERSION} -X ${PACKAGE_NAME}/internal.ReleaseDate=${APP_RELEASE_DATE} -X ${PACKAGE_NAME}/internal.ReleaseName=${APP_RELEASE}"

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "${LDFLAGS} -s -w" -o /out/az-health-exporter ./cmd/az-health-exporter


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/az-health-exporter .

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/az-health-exporter", "monitor", "--p", "8080"]

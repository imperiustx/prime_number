FROM --platform=${BUILDPLATFORM} golang:1.15.3-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/server-api

FROM scratch AS bin-unix
FROM bin-unix AS bin-darwin

COPY --from=build /out/server /

FROM bin-unix AS bin-linux

FROM bin-${TARGETOS} AS bin


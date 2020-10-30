FROM golang:1.15.3-alpine AS build

WORKDIR /src/
COPY . /src/

RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server-api
RUN CGO_ENABLED=0 go build -o /bin/admin ./cmd/server-admin

FROM scratch
COPY --from=build /bin/server /bin/server
COPY --from=build /bin/admin /bin/admin

CMD [ "/bin/admin", "keygen", "private.pem" ]

COPY . .

ENTRYPOINT ["/bin/server"]
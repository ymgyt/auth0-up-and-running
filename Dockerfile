FROM golang:1.11.2-alpine3.8 as build

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/ymgyt/auth0-up-and-running

COPY . .

# depenency managementどうするか悩み中(dep, module, vendoring...)

RUN CGO_ENABLED=0 go build -o /auth0-uar


FROM alpine:3.8

WORKDIR /root

COPY --from=build /auth-uar .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/ymgyt/auth0-up-and-running/static ./static
COPY --from=build /go/src/github.com/ymgyt/auth0-up-and-running/templates ./templates

EXPOSE 80

ENTRYPOINT [ "./auth-uar" ]
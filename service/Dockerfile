################################################################
#  BUILD
################################################################
FROM gcr.io/okcredit-42/golang:1.13.7 AS builder
COPY . /src
WORKDIR /src/service
RUN make build OUT=/service



################################################################
#  MAIN
#################9##############################################
FROM alpine:3.8
RUN echo "https://dl-4.alpinelinux.org/alpine/v3.9/main" >/etc/apk/repositories
RUN echo "https://dl-4.alpinelinux.org/alpine/v3.9/community" >>/etc/apk/repositories
RUN apk add --no-cache --update curl ca-certificates && update-ca-certificates
COPY --from=builder /service /service
HEALTHCHECK CMD curl --fail http://0.0.0.0:8080/ || exit 1
CMD ["/service"]
EXPOSE 8080
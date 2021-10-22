FROM golang:1.15.8-alpine as builder
RUN apk update && apk upgrade && apk add --no-cache bash libc6-compat git openssh
WORKDIR /notes-api
COPY . .
RUN rm -f go.sum
RUN go build -o ./app ./main.go

FROM alpine:3.13.1
WORKDIR /app
COPY --from=builder /notes-api/app .
EXPOSE 8006
CMD ["./app"]

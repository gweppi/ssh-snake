FROM golang:alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o ssh-snake

FROM scratch
WORKDIR /app
COPY --from=builder /app/ssh-snake .
EXPOSE 23234
CMD [ "./ssh-snake" ]
# Start from the latest golang base image
FROM golang:1.22.3 as builder
WORKDIR /app
COPY . .

ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/gharsallahmoez/*
RUN git config --global url."https://${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# CGO has to be disabled for alpine
RUN CGO_ENABLED=0 GOOS=linux make build

######## Start a new stage from scratch #######
FROM alpine:3.17.1
RUN apk update && apk add --no-cache ca-certificates bash git

WORKDIR /app/
COPY --from=builder /app/bin/messages/messagessvc .
RUN chmod +x messagessvc

EXPOSE 8080
CMD [ "./messagessvc" ]
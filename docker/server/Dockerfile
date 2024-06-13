FROM --platform=linux/amd64 alpine
RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
COPY bin/server server
EXPOSE 3000
CMD ["./bin/server"]
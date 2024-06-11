FROM --platform=linux/amd64 alpine
RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
COPY server server
EXPOSE 3000
CMD ["./server"]
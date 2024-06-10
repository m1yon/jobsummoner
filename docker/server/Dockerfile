FROM --platform=linux/amd64 alpine
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY . .
EXPOSE 3000
CMD ["./server"]
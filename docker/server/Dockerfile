FROM --platform=linux/amd64 alpine

ENV LOCAL_DB false
ENV DATABASE_URL ""

RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
COPY . .
EXPOSE 3000

CMD ./bin/server -local-db=${LOCAL_DB} -dsn=${DATABASE_URL}
FROM --platform=linux/amd64 alpine

ENV LOCAL_DB false
ENV DATABASE_URL ""
ENV PROXY_HOSTNAME ""
ENV PROXY_PORT ""
ENV PROXY_USERNAME ""
ENV PROXY_PASSWORD ""


RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
COPY . .

CMD ./bin/scraper -local-db=${LOCAL_DB} -dsn=${DATABASE_URL} -proxy-hostname=${PROXY_HOSTNAME} -proxy-port=${PROXY_PORT} -proxy-username=${PROXY_USERNAME} -proxy-password=${PROXY_PASSWORD}
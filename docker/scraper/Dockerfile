FROM --platform=linux/amd64 alpine
RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app
COPY bin/scraper scraper
CMD ["./bin/scraper"]
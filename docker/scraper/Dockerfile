FROM --platform=linux/amd64 alpine
RUN apk update && apk add ca-certificates tzdata
WORKDIR /app
COPY scraper scraper
CMD ["./scraper"]
FROM --platform=linux/amd64 ubuntu:jammy

ARG apt_sources="http://archive.ubuntu.com"

RUN sed -i "s|http://archive.ubuntu.com|$apt_sources|g" /etc/apt/sources.list && \
  apt-get update > /dev/null && \
  apt-get install --no-install-recommends -y \
  # chromium dependencies
  libnss3 \
  libxss1 \
  libasound2 \
  libxtst6 \
  libgtk-3-0 \
  libgbm1 \
  ca-certificates \ 
  # fonts
  fonts-liberation fonts-noto-color-emoji fonts-noto-cjk \
  # timezone
  tzdata \
  # process reaper
  dumb-init \
  # headful mode support, for example: $ xvfb-run chromium-browser --remote-debugging-port=9222
  xvfb \
  # misc
  make \
  > /dev/null && \
  # cleanup
  rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

VOLUME /app/db

EXPOSE 3000

CMD scripts/run-docker.sh
FROM golang:alpine3.20 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . ./
RUN go mod download

# Install Make and Git
RUN apk add --no-cache make git

# Build binary
RUN make build

FROM ubuntu:latest

# install cron
RUN apt-get update && apt-get -y install cron

# copy the binary from the previous state
COPY --from=build /app/log-checker /bin/log-checker

# config file
COPY .log-checker.yaml /root/.log-checker.yaml

# cron file
COPY <<EOF /etc/cron.d/log-checker-cron
* * * * * /bin/log-checker run -b

EOF

# Give execution rights on the cron job
RUN chmod 0644 /etc/cron.d/log-checker-cron

# Apply cron job
RUN crontab /etc/cron.d/log-checker-cron

# Run the command on container startup
CMD printenv > /etc/environment && cron && tail --follow=name --retry /root/log/app.log 2>/dev/null


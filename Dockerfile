# Start from the latest alpine golang base image
FROM golang:alpine as builder
LABEL maintainer="VksoR2FpZGFyb3Mp"

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .
RUN go get
RUN go build -o donkey main.go

######## Start a new stage from scratch #######
FROM scratch

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/donkey /usr/local/bin/donkey

# Command to run the executable
# CMD ["/usr/local/bin/donkey"]
CMD ["donkey"]

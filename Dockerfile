FROM golang:alpine AS builder
LABEL maintainer="VksoR2FpZGFyb3Mp"

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .
RUN apk add --no-cache git gcc musl-dev && \
    go get && \
    go build -tags musl -o donkey main.go
    # using "musl" tag to use the compiled misl lib as described in [go-fitz docs](https://github.com/gen2brain/go-fitz?tab=readme-ov-file#build-tags)
    # should consider using [pdfcpu](https://github.com/pdfcpu/pdfcpu)

######## Start a new stage from scratch #######
FROM alpine

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/donkey /usr/local/bin/donkey

# Command to run the executable
# CMD ["/usr/local/bin/donkey"]
CMD ["donkey"]

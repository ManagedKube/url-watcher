####################
# builder
####################
FROM golang:1.12-stretch as builder
#FROM golang:1.11.0-stretch as builder

# Update certificates
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY ./endpoint-runner $GOPATH/src/endpoint-runner
WORKDIR $GOPATH/src/endpoint-runner

# build
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o /go/bin/endpoint-runner

# Create a "nobody" non-root user for the next image by crafting an /etc/passwd
# file that the next image can copy in. This is necessary since the next image
# is based on scratch, which doesn't have adduser, cat, echo, or even sh.
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

####################
# start from scratch
####################
FROM scratch

# Copy our static executable
COPY --from=builder /go/bin/endpoint-runner /endpoint-runner
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the /etc_passwd file we created in the builder stage into /etc/passwd in
# the target stage. This creates a new non-root user as a security best
# practice.
COPY --from=builder /etc_passwd /etc/passwd

# Run as the new non-root by default
USER nobody

ENTRYPOINT ["/endpoint-runner"]

FROM golang:1.19-alpine

# Install required packages
RUN apk add --no-cache curl

# Download and install kubectl binary
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# Download and install Helm binary
RUN curl -LO "https://get.helm.sh/helm-v3.7.1-linux-amd64.tar.gz" \
    && tar -xf helm-v3.7.1-linux-amd64.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && rm -rf helm-v3.7.1-linux-amd64.tar.gz linux-amd64
    
# Set the working directory
WORKDIR /app

# Copy the Go application to the container
COPY . .

# Build the Go application
#RUN go build -o helm-cli .
RUN go build -o helm-cli ./cmd/main.go

ENTRYPOINT ["./helm-cli"]

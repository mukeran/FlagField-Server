FROM ubuntu:19.10

ENV DEBIAN_FRONTEND=noninteractive

# Replace apt source with byr mirror
# RUN sed -i 's/archive.ubuntu.com/mirrors.byrio.org/g' /etc/apt/sources.list
# RUN sed -i 's/security.ubuntu.com/mirrors.byrio.org/g' /etc/apt/sources.list

# Install essentials
RUN apt-get update && apt-get install -yq --no-install-recommends \
    ca-certificates \
    apt-utils \
    build-essential \
    wget \
    git \
    ssh \
    mysql-client \
    redis

# Go installation
RUN wget -P /root/ 'https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz'
RUN tar -xvf /root/go1.13.4.linux-amd64.tar.gz -C /usr/local/
RUN ln -s /usr/local/go/bin/go /usr/bin/go
RUN ln -s /usr/local/go/bin/gofmt /usr/bin/gofmt

# Set go environment
# RUN go env -w GOPROXY=https://goproxy.cn,direct

# Copy source and make
COPY ./ /opt/FlagField/
RUN rm -rf /opt/FlagField/config.json
RUN cd /opt/FlagField && make

# Make docker-entrypoint.sh executable
RUN chmod +x /opt/FlagField/docker-entrypoint.sh

# Set entrypoint
ENTRYPOINT ["/opt/FlagField/docker-entrypoint.sh"]

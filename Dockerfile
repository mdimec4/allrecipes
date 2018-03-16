FROM ubuntu:17.10 
#FROM armv7/armhf-ubuntu:16.10

ADD . /tmp/repo/src/allrecipes
WORKDIR /opt/webapi

ENV GOPATH="/tmp/repo"
ENV PATH="/tmp/repo/bin:${PATH}"

RUN apt-get update && \
    apt-get install -y --no-install-recommends --reinstall ca-certificates && \
    apt-get install -y --no-install-recommends build-essential \
    golang \
    git \
    wget && \
    # install dumb init
    cd /tmp && \
    wget -P /tmp https://github.com/Yelp/dumb-init/archive/v1.2.0.tar.gz && \
    tar xzf /tmp/v1.2.0.tar.gz -C /tmp/ && \
    cd /tmp/dumb-init-1.2.0 && \
    make && \
    cp dumb-init /sbin && \
    chmod +x /sbin/dumb-init && \
    cd / && \
    rm -rf /tmp/dumb-init-1.2.0 /tmp/v1.2.0.tar.gz && \
    # build
    cd /tmp/repo/src/allrecipes && \
    go get github.com/golang/dep/cmd/dep && \ 
    dep ensure && \
    go build allrecipes/cmd/webapi && \
    mkdir -p /opt/webapi && \
    cp ./webapi /opt/webapi && \
    cd / && \
    rm -r /tmp/repo && \
    # cleanup
    apt-get remove -y build-essential \
    golang \
    git \
    wget && \
    apt-get autoremove -y && \
    apt-get clean


# Runs "/usr/bin/dumb-init -- /my/script --with --args"
ENTRYPOINT ["/sbin/dumb-init", "--"]
# or if you use --rewrite or other cli flags
# ENTRYPOINT ["dumb-init", "--rewrite", "2:3", "--"]
# CMD ["/my/script", "--with", "--args"]
CMD ["./webapi"]





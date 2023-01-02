FROM golang:1.19

RUN curl -Lso staticcheck_linux_amd64.tar.gz \
    "https://github.com/dominikh/go-tools/releases/download/v0.3.3/staticcheck_linux_amd64.tar.gz" && \
    tar -xvf staticcheck_linux_amd64.tar.gz && \
    mv staticcheck/staticcheck /usr/local/bin/staticcheck && \
    rm -rf staticcheck*

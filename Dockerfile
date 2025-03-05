
FROM golang:1.23.7-alpine AS packages_image
ENV PATH /usr/local/go/bin:$PATH
WORKDIR /go/src
COPY go.mod go.sum ./
RUN echo 'packages_image' && \
    go mod download

FROM golang:1.23.7-alpine AS build_image
LABEL temporary=''
ARG BRANCH
WORKDIR /go/
COPY --from=packages_image /go ./
RUN echo -e "build_image" && \
    rm -r ./src && \
    apk update && \
    apk add --no-cache git && \
    git clone -b ${BRANCH} https://github.com/av-belyakov/placeholder_misp.git  ./src/ && \
    go build -C ./src/cmd/ -o ../app

FROM alpine
LABEL author='Artemij Belyakov'
ARG VERSION
ARG USERNAME=dockeruser
ARG US_DIR=/opt/placeholder_misp
RUN addgroup --g 1500 groupcontainer
RUN adduser -u 1500 -G groupcontainer -D ${USERNAME} --home ${US_DIR}
USER ${USERNAME}
WORKDIR ${US_DIR}
RUN mkdir ./logs
COPY --from=build_image /go/src/app ./
COPY --from=build_image /go/src/README.md ./ 
COPY config/* ./config/

ENTRYPOINT [ "./app" ]

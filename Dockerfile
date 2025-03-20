FROM golang:1.23.7-alpine AS pmisp_packages_image
ENV PATH /usr/local/go/bin:$PATH
WORKDIR /go/src
COPY go.mod go.sum ./
RUN echo 'packages_image' && \
    go mod download

FROM golang:1.23.7-alpine AS pmisp_build_image
LABEL temporary=''
ARG BRANCH
ARG VERSION
WORKDIR /go/
COPY --from=pmisp_packages_image /go ./
RUN echo -e "pmisp_build_image" && \
    rm -r ./src && \
    apk update && \
    apk add --no-cache git && \
    git clone -b ${BRANCH} https://github.com/av-belyakov/placeholder_misp.git  ./src/${VERSION}/ && \
    go build -C ./src/${VERSION}/cmd/ -o ../app

FROM alpine
LABEL author='Artemij Belyakov'
#аргумент STATUS содержит режим запуска приложения prod или development
#если значение содержит запись development, то в таком режиме и будет
#работать приложение, во всех остальных случаях режим работы prod
ARG STATUS
ARG VERSION
ARG USERNAME=dockeruser
ARG US_DIR=/opt/placeholder_misp
ENV GO_PHMISP_MAIN=${STATUS}
RUN addgroup --g 1500 groupcontainer && \ 
    adduser -u 1500 -G groupcontainer -D ${USERNAME} --home ${US_DIR}
USER ${USERNAME}
WORKDIR ${US_DIR}
RUN mkdir ./logs
COPY --from=pmisp_build_image /go/src/${VERSION}/app ./
COPY --from=pmisp_build_image /go/src/${VERSION}/README.md ./ 
COPY --from=pmisp_build_image /go/src/${VERSION}/version ./
COPY --from=pmisp_build_image /go/src/${VERSION}/backupdb/sqlite3_backup.db ./backupdb/
COPY config/* ./config/

ENTRYPOINT [ "./app" ]

# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-golang-api"
LABEL REPO="https://github.com/deadpyxel/company-api"

ENV PROJPATH=/go/src/github.com/deadpyxel/company-api

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

COPY . /go/src/github.com/deadpyxel/company-api
WORKDIR /go/src/github.com/deadpyxel/company-api

RUN make build

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/deadpyxel/company-api"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/company-api/bin

WORKDIR /opt/company-api/bin

COPY --from=build-stage /go/src/github.com/deadpyxel/company-api/company-api /opt/company-api/bin/
RUN chmod +x /opt/company-api/bin/company-api

# Create appuser
RUN adduser -D -g '' company-api
USER company-api

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
EXPOSE 8000

CMD ["/opt/company-api/bin/company-api"]

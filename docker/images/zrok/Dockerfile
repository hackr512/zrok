# this builds docker.io/openziti/zrok
ARG ZITI_CLI_TAG="1.0.0"
ARG ZITI_CLI_IMAGE="docker.io/openziti/ziti-cli"

FROM ${ZITI_CLI_IMAGE}:${ZITI_CLI_TAG}

ARG ARTIFACTS_DIR=./dist
ARG DOCKER_BUILD_DIR=.
# e.g. linux
ARG TARGETOS
# e.g. arm64
ARG TARGETARCH

### Required OpenShift Labels 
LABEL name="openziti/zrok" \
      maintainer="support@zrok.io" \
      vendor="NetFoundry" \
      summary="Run the zrok CLI" \
      description="Run the zrok CLI"

USER root

### install packages: findutils provides xargs which is used by the zrok Helm chart's controller bootstrapping script to
#create the default account enable token
RUN   INSTALL_PKGS="findutils" && \
      microdnf -y update --setopt=install_weak_deps=0 --setopt=tsflags=nodocs && \
      microdnf -y install --setopt=install_weak_deps=0 --setopt=tsflags=nodocs ${INSTALL_PKGS}

### add licenses to this directory
RUN mkdir -p -m0755 /licenses
COPY ./LICENSE /licenses/apache.txt

RUN mkdir -p /usr/local/bin
COPY ${ARTIFACTS_DIR}/${TARGETARCH}/${TARGETOS}/zrok \
      ./nfpm/zrok-enable.bash \
      ./nfpm/zrok-share.bash \
      /usr/local/bin/
RUN chmod 0755 \
      /usr/local/bin/zrok \
      /usr/local/bin/zrok-enable.bash \
      /usr/local/bin/zrok-share.bash

USER ziggy
ENTRYPOINT [ "zrok" ]

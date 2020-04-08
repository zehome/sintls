FROM debian:bullseye
MAINTAINER Laurent Coustet <ed@zehome.com>

RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/norecommands && \
    echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/norecommands
RUN export DEBIAN_FRONTEND=noninteractive; \
    apt-get -qqy update && apt-get -y install \
    dh-golang \
    git-buildpackage \
    git-core \
    golang-1.14 \
    golang-go \
    ssh \
    sudo

RUN apt-get clean

RUN useradd -m -s /bin/bash -d /ci -G sudo ci && \
    echo "ci ALL=(root) NOPASSWD: ALL" > /etc/sudoers.d/ci && \
    chmod 440 /etc/sudoers.d/ci && \
    mkdir -p /ci/.pip /ci/.ssh
RUN mkdir -p /ci/.ssh && \
    echo 'StrictHostKeyChecking no' >> /ci/.ssh/config && \
    mkdir -p /ci/.pip

RUN chown ci:ci -R /ci

USER ci
WORKDIR /ci
CMD ["/bin/bash"]


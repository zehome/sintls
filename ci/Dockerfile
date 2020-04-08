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
CMD ["/bin/bash"]

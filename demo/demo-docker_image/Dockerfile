ARG BUILDROOT=/pg
ARG PGROOT=/usr/local/pgsql
ARG PGUSER=mtmuser
ARG PGUID=5000
ARG PGGROUP=mtmuser
ARG PGGID=5000

FROM ubuntu:22.04 as build
ARG BUILDROOT
ARG PGROOT
ARG PGUSER
ARG PGUID
ARG PGGROUP
ARG PGGID

ENV BUILDROOT=$BUILDROOT
ENV PGROOT=$PGROOT
ENV PGUSER=$PGUSER
ENV PGUID=$PGUID
ENV PGGROUP=$PGGROUP
ENV PGGID=$PGGID

RUN apt-get update && \
    apt-get -y install gcc make bison flex libreadline-dev zlib1g-dev perl golang git

RUN mkdir BUILDROOT
COPY postgres_cluster $BUILDROOT/postgres_cluster

RUN cd $BUILDROOT/postgres_cluster && \
    ./configure && \
    make clean && \
    make && \
    make install

RUN cd $BUILDROOT/postgres_cluster/contrib/mmts && \
    make clean && \
    make install

FROM ubuntu:22.04
ARG BUILDROOT
ARG PGROOT
ARG PGUSER
ARG PGUID
ARG PGGROUP
ARG PGGID

ENV BUILDROOT=$BUILDROOT
ENV PGROOT=$PGROOT
ENV PGUSER=$PGUSER
ENV PGUID=$PGUID
ENV PGGROUP=$PGGROUP
ENV PGGID=$PGGID

RUN apt-get update && \
    apt-get -y install libreadline-dev curl less vim netcat net-tools lsof gdb rr gdbserver iproute2 iputils-ping iptables

RUN groupadd -g $PGUID $PGGROUP
RUN useradd -m $PGUSER -u $PGUID -g $PGGID

USER $PGUSER

COPY --chown=$PGUID:$PGGID --from=build $PGROOT $PGROOT

COPY scripts /scripts

RUN touch "$PGROOT"/mtm-lockfile

CMD ['/bin/sh', '/scripts/entrypoint.sh']

FROM golang:1.15.2-buster AS build

MAINTAINER Pavel Smirnov

RUN mkdir /go/src/forum-api

COPY . /go/src/forum-api

WORKDIR /go/src/forum-api

RUN make build

FROM ubuntu:20.04 AS release

MAINTAINER Pavel Smirnov

RUN mkdir /opt/forum-api

# install Postgres
ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib
RUN apt install make

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt installed``
USER postgres

COPY --from=build /go/src/forum-api /opt/forum-api

WORKDIR /opt/forum-api




RUN /etc/init.d/postgresql start &&\
    psql -f configs/sql/init.sql &&\
    make dbsetupc &&\
    /etc/init.d/postgresql stop

RUN echo "listen_addresses='*'\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf



# Add VOLUMEs to allow backup of config, logs and databases
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Собранный сервер
COPY --from=build /go/src/forum-api/build/bin/forum /usr/bin/forum

EXPOSE 5432
EXPOSE 5000

CMD service postgresql start && forum
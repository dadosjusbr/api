FROM postgres

WORKDIR /database

COPY ./init_db.sql /docker-entrypoint-initdb.d/

USER root
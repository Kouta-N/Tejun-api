FROM mysql:8.0
ENV LANG ja_JP.UTF-8

COPY ./database/*.sql /docker-entrypoint-initdb.d/
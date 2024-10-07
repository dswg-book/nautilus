FROM golang as build-image
RUN apt update && apt-get install -yy build-essential make
COPY . /app
WORKDIR /app/
RUN make build-server

FROM ubuntu
LABEL org.opencontainers.image.source="https://github.com/dswg-book/nautilus"
COPY --from=build-image /app/build/nautilus-server /usr/bin/
EXPOSE 3030
ENTRYPOINT [ "nautilus-server" ]
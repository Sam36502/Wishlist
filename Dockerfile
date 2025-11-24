FROM alpine:latest

ADD data /data

EXPOSE 80

ENV WISHLIST_PORT 80
ENV WISHLIST_DB_FILENAME /data/wishlist.db
ENV WISHLIST_BASE_URL https://wishlist.pearcenet.ch

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

CMD ["/data/run_server_linux"]

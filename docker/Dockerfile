FROM alpine

COPY ./release/pinmonl-linux /usr/local/bin/pinmonl
WORKDIR /app
RUN chmod 755 /usr/local/bin/pinmonl \
  && chmod -R 755 /app

CMD ["pinmonl", "server"]

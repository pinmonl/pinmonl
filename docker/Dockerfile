FROM alpine

COPY ./release/pinmonl-linux /app/pinmonl
WORKDIR /app
RUN chmod -R 755 /app \
  && ln -sf /app/pinmonl /usr/local/bin/pinmonl

CMD ["pinmonl", "server"]

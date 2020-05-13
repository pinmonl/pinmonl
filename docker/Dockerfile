FROM alpine
ARG BUILD_OS
ARG BUILD_ARCH

COPY release/pinmonl-${BUILD_OS}-${BUILD_ARCH} /usr/local/bin/pinmonl
WORKDIR /pinmonl

EXPOSE 3399

ENTRYPOINT ["pinmonl"]
CMD ["server"]

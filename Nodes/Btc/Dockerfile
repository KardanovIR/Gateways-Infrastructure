FROM alpine
ARG VERSION=0.18.1
ENV FILENAME bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz
ENV DOWNLOAD_URL https://bitcoin.org/bin/bitcoin-core-${VERSION}/${FILENAME}

RUN apk update

RUN wget $DOWNLOAD_URL
RUN tar xzvf /bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz
RUN mkdir /root/.bitcoin
RUN mv /bitcoin-${VERSION}/bin/* /usr/local/bin/
RUN rm -rf /bitcoin-${VERSION}/
RUN rm -rf /bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz

EXPOSE 8332 8333 18332 18333 28332 28333
ADD ./bin/docker_entrypoint.sh /usr/local/bin/docker_entrypoint.sh
RUN chmod a+x /usr/local/bin/docker_entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker_entrypoint.sh"]
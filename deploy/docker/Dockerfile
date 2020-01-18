FROM cchantep/libvips:amazonlinux1

COPY nuggan /root/
COPY server.conf /root/

ENV LD_LIBRARY_PATH /opt/local/lib

EXPOSE 8080/tcp
ENTRYPOINT ["/root/nuggan", "-server", ":8080", "-server-config", "/root/server.conf"]

FROM coredns/coredns:latest

# Copy custom CoreDNS configuration file into the container
COPY Corefile /Corefile

# Expose standard DNS ports (UDP and TCP)
EXPOSE 53/udp
EXPOSE 53/tcp

# Run CoreDNS using the custom Corefile
ENTRYPOINT ["/coredns", "-conf", "/Corefile"]
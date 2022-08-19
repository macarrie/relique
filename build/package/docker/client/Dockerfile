FROM golang:1.18-bullseye as builder

#RUN apt-get install -y bash gcc g++ make git && git clone https://github.com/macarrie/relique --branch ${version} /src
RUN apt-get install -y bash gcc g++ make
COPY . /src

WORKDIR /src
RUN make build_client install INSTALL_ROOT=./output/install INSTALL_ARGS="--client --skip-user-creation"

#-----------------------------
FROM debian:stable-slim

RUN mkdir -p /etc/relique/ /var/lib/relique /var/log/relique

# Copy application files
RUN apt-get update -y && apt-get install -y rsync
WORKDIR /app
COPY --from=builder /src/output/install/usr/bin/relique-client /usr/bin/
COPY --from=builder /src/output/install/etc/relique /etc/relique
COPY --from=builder /src/output/install/etc/relique/client.toml.sample /etc/relique/client.toml
COPY --from=builder /src/output/install/var/lib/relique /var/lib/relique

RUN groupadd --gid 8433 relique && useradd -r -d /var/lib/relique -s /bin/sh -g relique -u 8433 relique
USER relique

# Run
EXPOSE 8433
CMD ["/usr/bin/relique-client", "start"]


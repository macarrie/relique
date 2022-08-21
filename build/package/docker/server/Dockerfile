FROM node:18-bullseye as builder

RUN apt-get update && apt-get install -y bash gcc g++ make golang
COPY . /src

WORKDIR /src
RUN make build_server install INSTALL_ROOT=./output/install INSTALL_ARGS="--server --skip-user-creation"

#-----------------------------
FROM debian:stable-slim

RUN mkdir -p /etc/relique/ /var/lib/relique /var/log/relique

# Copy application files
RUN apt-get update -y && apt-get install -y rsync ssh
COPY --from=builder /src/output/install/usr/bin/relique-server /usr/bin/
COPY --from=builder /src/output/install/etc/relique /etc/relique
COPY --from=builder /src/output/install/etc/relique/server.toml.sample /etc/relique/server.toml
COPY --from=builder /src/output/install/var/lib/relique /var/lib/relique

RUN groupadd --gid 833 relique && useradd -r -d /var/lib/relique -s /bin/sh -g relique -u 833 relique
RUN mkdir ~relique/.ssh && ssh-keygen -t rsa -f ~relique/.ssh/id_rsa -N ''
RUN chown -R relique:relique /var/lib/relique /var/log/relique /etc/relique
USER relique

# Run
EXPOSE 8433
CMD ["/usr/bin/relique-server", "start"]


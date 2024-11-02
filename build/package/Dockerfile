FROM golang:1.22 AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build



FROM rockylinux/rockylinux:9-minimal
RUN microdnf install -y openssh-clients rsync \
    && microdnf clean all \
    && rm -rf /var/cache/yum
RUN groupadd -g 8433 relique \
    && useradd -u 8433 -g 8433 relique 
COPY --from=base /app/output/relique /bin/relique
RUN relique config init && chown -R relique:relique /etc/relique /var/lib/relique

USER relique
CMD [ "/bin/bash" ]
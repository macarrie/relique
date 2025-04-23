FROM node:23-bookworm AS webui
WORKDIR /app
RUN apt install -y make
COPY . .
RUN make webui


FROM golang:1.24 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN apt install -y git
RUN go mod download
COPY . .
COPY --from=webui /app/webui/dist /app/internal/server/dist
RUN ls /app/internal/server
RUN make bin


FROM rockylinux/rockylinux:9-minimal
RUN microdnf install -y openssh-clients rsync \
    && microdnf clean all \
    && rm -rf /var/cache/yum
RUN groupadd -g 8433 relique \
    && useradd -u 8433 -g 8433 relique 
COPY --from=base /app/output/relique /bin/relique
RUN relique config init && chown -R relique:relique /etc/relique /var/lib/relique

USER relique
CMD [ "/bin/relique", "server", "start" ]
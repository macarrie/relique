FROM golang:1.22 AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build



FROM rockylinux/rockylinux:9-minimal
RUN groupadd -g 8433 relique 
RUN useradd -u 8433 -g 8433 relique 
COPY --from=base /app/output/relique /bin/relique
RUN relique config init

USER relique
CMD [ "/bin/bash" ]
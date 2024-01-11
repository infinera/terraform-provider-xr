FROM golang:1.18.0-alpine3.15 AS builder
LABEL stage=builder
RUN apk add --no-cache bash git gcc libc-dev make cmake
ADD . /workspace/terraform-provider-xr
RUN cd /workspace/terraform-provider-xr \
    && make build install && cd /tmp \
    && git clone -b v0.0.6 https://github.com/infinera/terraform-xr-network-setup.git

FROM alpine:3.16.2 AS final
RUN apk add --no-cache bash terraform git
WORKDIR /xr_terraform/network-setup/
COPY --from=builder /root/.terraform.d /root/.terraform.d
COPY --from=builder /tmp/terraform-xr-network-setup/use-git-Repo/*.tf /xr_terraform/network-setup/
RUN  terraform init

CMD ["bash"]

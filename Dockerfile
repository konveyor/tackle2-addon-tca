FROM registry.access.redhat.com/ubi8/go-toolset:latest as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM quay.io/konveyor/tackle-container-advisor:v2.0.1
WORKDIR /app
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]

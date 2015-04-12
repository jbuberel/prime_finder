# Dockerfile extending the generic Go image with application files for a
# single application.
FROM ubuntu:latest

ENV MYSQL=173.194.84.87
COPY server /server
ENTRYPOINT ["/server"]

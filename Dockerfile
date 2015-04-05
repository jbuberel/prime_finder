# Dockerfile extending the generic Go image with application files for a
# single application.
FROM gcr.io/google_appengine/go-compat

ADD . /app
RUN /bin/bash /app/_ah/build.sh

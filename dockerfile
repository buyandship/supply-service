FROM --platform=linux/amd64 golang:1.22

ARG SERVICE=supply-svr
ARG PORT=8573

WORKDIR /app

# Copy the required binary
COPY ./supply-svr /app/supply-svr

# Copy the config.yaml file
COPY ./config.test.yml /app/config.yaml

# Set the environment variable for the port
ENV PORT=$PORT

# Expose the port
EXPOSE $PORT

CMD ["./supply-svr"]
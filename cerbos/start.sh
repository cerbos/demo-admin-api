docker run -i -t -p 3592:3592 -p 3593:3593 \
  -v $(pwd)/config:/config \
  -v $(pwd)/policies:/policies \
  ghcr.io/cerbos/cerbos:latest \
  server \
  --config=/config/conf.yaml

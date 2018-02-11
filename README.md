[![Build Status](https://travis-ci.org/blockassets/litecoinpool_exporter.svg?branch=master)](https://travis-ci.org/blockassets/litecoinpool_exporter)

# Litecoinpool Exporter

[Prometheus.io](https://prometheus.io/) exporter for the [litecoinpool.org](https://litecoinpool.org) API.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

### Usage (defaults):

``
LITECOINPOOL_API_KEY=2l3kj2l3kj43lj ./litecoinpool_exporter -port 5551 -timeout 10s
``

If you don't want to specify the API key in the environment, you can do it in a file:

``
echo "l2j3l4k23j4l3j" > ./litecoinpool-api-key.txt
./litecoinpool_exporter -key-file litecoinpool-api-key.txt
``

Note: if you remove a worker, you need to restart the exporter.

### Setup

Install [dep](https://github.com/golang/dep) and the dependencies...

`make dep`

### Build binary for amd64

`make amd64`

### Install onto miner

The [releases tab](https://github.com/blockassets/litecoinpool_exporter/releases) has `master` binaries cross compiled for AMD64. These are built automatically on [Travis](https://travis-ci.org/blockassets/litecoinpool_exporter).

Download the latest release and copy the `litecoinpool_exporter` binary to `/usr/local/bin`

```
chmod ugo+x litecoinpool_exporter
scp litecoinpool_exporter root@SERVER_IP:/usr/local/bin
```

Create `/etc/systemd/system/litecoinpool_exporter.service`

```
ssh root@SERVER_IP "echo '
[Unit]
Description=litecoinpool_exporter
After=init.service

[Service]
Type=simple
ExecStart=/usr/local/bin/litecoinpool_exporter -key-file /usr/local/etc/litecoinpool-api-key.txt
Restart=always
RestartSec=4s
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/litecoinpool_exporter.service"
```

Enable the service:

```
ssh root@MINER_IP "systemctl enable litecoinpool_exporter; systemctl start litecoinpool_exporter"
```

### Test install

Open your browser to `http://SERVER_IP:5551/metrics`

### Prometheus configuration

`prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'lcp_exporter'
    static_configs:
      - targets: ['localhost:5551']
```

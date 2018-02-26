[![Build Status](https://travis-ci.org/blockassets/bwpool_exporter.svg?branch=master)](https://travis-ci.org/blockassets/bwpool_exporter)

# BW Pool Exporter

[Prometheus.io](https://prometheus.io/) exporter for the [BW.com](https://bw.com) Pool API.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

### Usage (defaults):

bwpool.json:

```
{
  "Username": "foo",
  "PublicKey": "your public key",
  "PrivateKey": "your private key"
}
```

```
./bwpool_exporter-linux-amd64 -config bwpool.json
```

Note: if you remove a worker, you need to restart the exporter.

### Setup

Install [dep](https://github.com/golang/dep) and the dependencies...

`make dep`

### Build binary for linux

`make linux`

### Install onto miner

The [releases tab](https://github.com/blockassets/bwpool_exporter/releases) has `master` binaries cross compiled for Linux AMD64 and Darwin. These are built automatically on [Travis](https://travis-ci.org/blockassets/bwpool_exporter).

Download the latest release and copy the `bwpool_exporter-linux-amd64` binary to `/usr/local/bin`

```
chmod ugo+x bwpool_exporter
scp bwpool_exporter root@SERVER_IP:/usr/local/bin
```

Create `/etc/systemd/system/bwpool_exporter.service`

```
ssh root@SERVER_IP "echo '
[Unit]
Description=bwpool_exporter
After=init.service

[Service]
Type=simple
ExecStart=/usr/local/bin/bwpool_exporter-linux-amd64 -key-file /usr/local/etc/litecoinpool-api-key.txt
Restart=always
RestartSec=4s
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/bwpool_exporter.service"
```

Enable the service:

```
ssh root@MINER_IP "systemctl enable bwpool_exporter; systemctl start bwpool_exporter"
```

### Test install

Open your browser to `http://SERVER_IP:5551/metrics`

### Prometheus configuration

`prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'bwpool_exporter'
    static_configs:
      - targets: ['localhost:5551']
```

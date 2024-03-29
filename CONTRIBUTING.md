## Contributing

### Requirements

- Go 1.22
- GNU Make
- Docker (optional)
- OPNsense Box with admin access

### Environment

This guide is for osx and Linux.

### Create API key and secret in OPNsense

`SYSTEM>ACCESS>USERS>[user]>API KEYS`

[OPNsense Documentation](https://docs.opnsense.org/development/how-tos/api.html#creating-keys)

### Run the exporter locally

```bash
OPS_ADDRESS="ops.example.com" OPS_API_KEY=your-api-key OPS_API_SECRET=your-api-secret make local-run
```

- test it

```bash
curl http://localhost:8080/metrics
```

### Before PR

- Make sure to sync the vendor if the dependencies have changed.

```bash
make sync-vendor
```

- Make sure to run the tests and linters.

```bash
make test
make lint
```

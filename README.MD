# gcp-disable-kms-key-versions

This will delete all BDD secrets in Google secrets manager for a given project using the patter of a secret name including `pr-`
## Run

First build the binary, requires go 1.15

Build
```bash
make build
```
Run
```bash
./build/gcp-disable-kms-key-versions
```

# garbage collect google secret manager secrets

This will delete secrets in Google secrets manager for a given project optionally matching a 
secret name containing a string passed in with a `--filter` flag.

## Run

First build the binary, requires go 1.15

Build
```bash
make build
```
Run e.g.
```bash
./build/gc-gsm-secrets --project jenkins-x-labs-bdd1 --filter pr-
```

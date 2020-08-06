How to run all tests
====================

1. Copy `.env.example` to `.env`
2. Fill `.env`
3. Run `dotenv go test ./...`



How to deploy
=============

```console
$ cd path/to/unity-meta-check

$ # Write the new version to deploy into version.go
$ edit ./version/version.go
$ git add ./version/version.go
$ git commit -m "Bump to $(./scripts/print-version)"
$ git push

$ # Deploy to GitHub releases and Docker registry and unity-meta-check-bins:
$ scripts/deploy
```

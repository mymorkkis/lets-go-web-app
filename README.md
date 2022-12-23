# Let's Go App

Repo for following the [Let's Go Book](https://lets-go.alexedwards.net/) by Alex Edwards.

### Generate a local TLS certificate

Find where the Go standard library is installed and run the `generate_cert` script.
E.G:

```
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

If Go was installed using something like homebrew then the path will be different.

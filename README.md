# go-pwdgen

<img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" width="64" height="64">

A small, **cryptographically secure** random password generator written in Go — usable as a
library, a CLI, or an HTTP service. Live demo deployed [here](http://pwdgen.poln.org).

## Features

- 🔒 **Secure by default** — passwords are drawn from [`crypto/rand`](https://pkg.go.dev/crypto/rand)
  (the system CSPRNG), not `math/rand`.
- 🎯 **Unbiased** — characters are selected with *rejection sampling*, so every symbol in the
  charset has an exactly uniform probability (no modulo bias).
- 🧩 **Three ways to use it** — Go package, command-line tool, and a Gin-based HTTP API.
- 🌐 **JSON or HTML** output from the web endpoint.

### Charset & length

Passwords are composed from the following 68-character set:

```
^AZERTYUIOPMLKJHGFDSQWXCVBN_#@?1234567890-.!azertyuiopmlkjhgfdsqwxcvbn
```

The requested length is clamped to the range **[10, 1048576]** (`MinLength`..`MaxLength`);
any value outside that range falls back to the minimum of `10`.

## Requirements

- Go **1.26** or later

## HTTP API

Start the server (listens on `:8080`, or on `$PORT` if set):

```sh
make run          # build + run, or:
go run ./app
```

### Query parameters

| Param  | Default | Description                                        |
|--------|---------|----------------------------------------------------|
| `l`    | `16`    | Password length (clamped to `[10, 1048576]`)       |
| `mime` | `html`  | Response format: `json` or `html`                  |

### Example — JSON

```sh
curl -Ls "pwdgen.poln.org/?l=64&mime=json" | jq
```

```json
{
  "length": 64,
  "password": "4wC3ZAJjw1C5X_p9AyAJDDNzu3BF8khCQlAu9mKIzNdGPVjVni_ftSEjer.3EJQS"
}
```

### Example — HTML

```sh
curl -Ls "pwdgen.poln.org/?l=64"
```

```html
<!DOCTYPE html>
<html>

<head>
    <title>PoLN | 🔑 PWDGEN</title>
    <style>
        table {
          font-family: arial, sans-serif;
          border-collapse: collapse;
          width: 50%;
        }

        td, th {
          border: 1px solid #dddddd;
          text-align: left;
          padding: 8px;
        }

        </style>
</head>

<body>
    <table>
        <tr>
            <th>Length</th>
            <th>Password</th>
        </tr>
        <tr>
            <td class="length"><code>64</code></td>
            <td><code>-3D8-7MAqq2IAciip7w2426iV18vWhgizaJ?cI?aCkDy#gnxgeAvJ7rGkveccI!I</code></td>
        </tr>
    </table>
</body>

</html>
```

## Command-line tool

```sh
go run ./cmd -l 32
# or install it:
go install github.com/fairhive-labs/go-pwdgen/cmd@latest
```

```text
Password length : 32
Code : 58mhKGFTyiom1.jTIGXO06fas_IiO_uM
```

A length below the minimum is reported and clamped:

```sh
$ go run ./cmd -l 5
provided length 5 is less than 10, changed to 10 !!!
Password length : 10
Code : HbvBkgKJqx
```

## Library

```go
import "github.com/fairhive-labs/go-pwdgen/pkg/generator"

pwd := generator.Generate(32) // 32-character cryptographically secure password
```

`Generate(size int) string` returns a password of the requested `size`, clamped to
`[generator.MinLength, generator.MaxLength]`.

## Development

| Command      | Description                          |
|--------------|--------------------------------------|
| `make build` | Build the web app into `./bin/app`   |
| `make run`   | Build and run the web app            |
| `make test`  | `go vet ./...` then `go test ./...`  |
| `make clean` | Remove the `./bin` directory         |

```sh
go test ./... -cover   # run the full suite with coverage
```

## License

Released under the [MIT License](LICENSE).

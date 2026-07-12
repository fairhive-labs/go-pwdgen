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

On a terminal it plays a short cyberpunk "forge" sequence before revealing the
password (all in green — the password included, so it reads as part of the
interface):

```text
╔══════════════════════════════════════╗
║             -- PWDGEN --             ║
║         SECURE FORGE ONLINE          ║
╚══════════════════════════════════════╝
> SEEDING CSPRNG................ [OK]
> REJECTION-SAMPLING ENTROPY.... [OK]
> FORGING 32-CHAR KEY........... [OK]
[████████████████] 100%
>>> PAYLOAD [32]:
-py4Uk?kf_@A9OiN8VP#_e5aTGkNicsf
// NO RIGHT PASSWORD, ONLY BETTER TOOLS
```

The task lines mirror what the generator actually does — seed the system CSPRNG,
draw bytes with rejection sampling, and forge an *N*-character key.

### Flags

| Flag     | Default | Description                                             |
|----------|---------|---------------------------------------------------------|
| `-l`     | `16`    | Password length (clamped to `[10, 1048576]`)            |
| `-plain` | `false` | Print only the bare password — no animation, for scripts |

### Scriptable by design

The animation and all decorative output go to **stderr**; the bare password is
written to **stdout**. The show also auto-disables when output isn't a terminal.
So piping stays clean — you capture only the credential:

```sh
pwdgen -l 32 | pbcopy        # copies just the password; the show still plays on your terminal
pwdgen -l 32 2>/dev/null     # bare password only
pwdgen -plain -l 32          # bare password only, animation forced off
```

A length below the minimum is reported (on stderr) and clamped to `10`:

```sh
$ pwdgen -l 5 2>&1 1>/dev/null
!! requested length 5 is below minimum, forced to 10
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

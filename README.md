# go-pwdgen
<img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" width="64" height="64">

Password generator written in Go... and deployed on [here](http://pwdgen.trendev.fr)

### E.g : Generate a 64 length random password 
> curl -Ls "pwdgen.trendev.fr/json?l=64" | jq

```
{
  "length": 64,
  "password": "4wC3ZAJjw1C5X_p9AyAJDDNzu3BF8khCQlAu9mKIzNdGPVjVni_ftSEjer.3EJQS"
}
```

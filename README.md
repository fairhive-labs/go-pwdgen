# go-pwdgen
<img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" width="64" height="64">

Password generator written in Go... and deployed [here](http://pwdgen.fairhive.io)

### E.g : Generate a 64 length random password ; MIME = JSON
> curl -Ls "pwdgen.fairhive.io/?l=64&mime=json" | jq

```
{
  "length": 64,
  "password": "4wC3ZAJjw1C5X_p9AyAJDDNzu3BF8khCQlAu9mKIzNdGPVjVni_ftSEjer.3EJQS"
}
```

### E.g : Generate a 64 length password ; MIME = HTML
> curl -Ls "pwdgen.fairhive.io/?l=64"

```
<!DOCTYPE html>
<html>

<head>
    <title>🔑 PWDGEN</title>
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
            <th>
                Length
            </th>
            <th>
                Password
            </th>
        </tr>
        <tr>
            <td><code>64</code></td>
            <td><code>-3D8-7MAqq2IAciip7w2426iV18vWhgizaJ?cI?aCkDy#gnxgeAvJ7rGkveccI!I</code></td>
        </tr>
    </table>
</body>

</html>
```
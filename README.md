# GoSpread - Google Spreadsheets utility
GoSpread wraps `google.golang.org/api/sheets/v4`.

## Features
- manage sheets
  - create, sort, move

## Installing
Use `go get` to install the library.
```sh
go get -u github.com/daizpan/gospread
```

Install `gospraed` command.
```sh
go get github.com/daizpan/gospread/cmd/gospread
``` 


## Usage
1. Enable sheets api to `https://console.developers.google.com/`.
2. Create service account, and create key and download to `creadential.json`
3. Share a sheets for service account email.

```sh
gospread sheets create sheet-name <spread-sheet-id>
```

your application.
```go
import "github.com/daizpan/gospread"

func main() {
	g, err := gospread.NewGoSpreadWithCredentialFile("credential.json")
	if err != nil {
		fmt.Fatal(err)
	}
	if err := g.CreateSheet("spread-sheet-id", "new-sheet-name"); err != nil {
		fmt.Printf("createTestSheet error=%s", err)
	}
}
```

## Author

[daizpan](https://github.com/daizpan)

## Licence

[MIT](https://github.com/daizpan/gospread/blob/main/LICENSE)

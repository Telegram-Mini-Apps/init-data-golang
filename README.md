# Init Data

Package which provides utilities to work with Telegram Web Apps init data.
To learn more, what init data is, visit its [documentation](https://github.com/Telegram-Web-Apps/documentation/blob/master/init-data.md).

## Installation

```bash
go get github.com/telegram-web-apps/init-data-golang
```

## Usage

### Validation

```go
package main

import (
	"fmt"
	"github.com/Telegram-Web-Apps/init-data-golang"
	"time"
)

func main() {
	// Init data in raw format.
	initData := "query_id=AAHdF6IQAAAAAN0XohDhrOrc&..."

	// Telegram Bot secret key.
	token := "627618978:amnnncjocxKJf"

	// Define how long since init data generation date it is valid.
	expIn := 24 * time.Hour

	// Will return error in case, init data is invalid. To see,
	// which error could be returned, see errors.go file.
	fmt.Println(initdata.Validate(initData, token, expIn))
}
```

In case, expiration time is `0`, function will skip expiration time check. It
is recommended to specify non-zero value as long as this check is considered
important to prevent old stolen init data usage. 

### Parsing

package main

```go
package main

import (
    "fmt"
    "github.com/Telegram-Web-Apps/init-data-golang"
)

func main() {
	// Init data in raw format.
	initData := "query_id=AAHdF6IQAAAAAN0XohDhrOrc&..."
	
	// Will return 2 values.
	// 1. Pointer to InitData in case, passed data has correct format.
	// 2. Error in case, something is wrong. 
	fmt.Println(initdata.Parse(initData))
}
```

It is important to note, that `Parse` function does not do any checks which
`Validate` does. So, this function just parses incoming data without `hash`
or expiration time validations.

### API

To see full package API documentation, visit [this link](https://pkg.go.dev/github.com/telegram-web-apps/init-data-golang).

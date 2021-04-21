# brickset
This is sample code for the beginnings of a Go Module for the Brickset.com API.

The [www.brickset.com](http://www.brickset.com) website provides an [API](https://brickset.com/api/v3.asmx) for working with their data. This Go module works with version 3 of that [API](https://brickset.com/api/v3.asmx).

**Pre-Requisites**

Created using Go 1.16

An API key from Brickset is required. Currently they are free and you can get one [here](https://brickset.com/tools/webservices/requestkey).
In order to use the inventory features of Brickset to track your own Lego collection a Brickset account is also required.

**Quick Start**

Use of non-inventory functions:

Supply your API key as a parameter for a function

```
package main

import (
	"fmt"

	"github.com/jonathanmedd/brickset/brickset"
)

func main() {

    test, err := brickset.GetThemes("4-e3wM-sWsw-Su3pI")

	if err != nil {
		fmt.Println(err)
		return
	}

    fmt.Println(test)
}
```

Use of all functions including inventory:

Supply your API key and also Brickset Website credentials to make a connection using the Login function. Use the returned hash in subsequent calls using inventory functions

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/jonathanmedd/brickset/brickset"
)

func main() {

	testConnection, err := brickset.Login("4-e3wM-sWsw-Su3pI", "test@test.com", "P@ssword!")

	if err != nil {
		fmt.Println(err)
		return
	}
	query, err := brickset.GetSets(testConnection.ApiKey, testConnection.Hash, 500, "Indiana Jones", "Raiders of the Lost Ark", "", "", 1, 0, "Pieces")
	if err != nil {
		fmt.Println(err)
		return
	}
	sets := query.Sets
	c, err := json.MarshalIndent(sets, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(c))
}
```

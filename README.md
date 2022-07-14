# goquery
A Go library to queries on different game servers

## Usage
```go
package main

import (
    "fmt"
    "github.com/Seyz123/goquery"
)

func main() {
    data, err := goquery.Query("play.symp.fr", 19132)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Server is running " + data.Software)
}
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the [Apache 2.0](https://choosealicense.com/licenses/apache-2.0/) license.

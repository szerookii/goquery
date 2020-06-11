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
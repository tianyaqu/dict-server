package main

import (
    "fmt"
    "flag"
)

var name = flag.String("name", "oald", "dict name")

func main() {
    flag.Parse()

    dt := NewStarDict("./" + *name)

    fmt.Printf("%s \n", dt)
}

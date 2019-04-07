package main

import (
    "flag"
    "github.com/gin-gonic/gin"
)

var (
    name = flag.String("name", "oald", "dict name")
    dictMap map[string]Dict
)

func init() {
    dictMap = make(map[string]Dict)
    dt := NewStarDict(*name)
    dictMap[dt.Name()] = dt
}

func main() {
    flag.Parse()

    router := gin.Default()
    router.GET("/word/:term", handleLookup)
    router.Run(":8080")
}

package main

import (
    "flag"
    "github.com/gin-gonic/gin"
)

var (
    name = flag.String("name", "oald", "dict name")
    dir = flag.String("dir", "../data/", "dict data path")
    dictMap map[string]Dict
    dictGuess map[string]GuessFunc
)

func init() {
    registerDictGuesses()
}

func main() {
    flag.Parse()

    dictMap = loadDictsFromDir(*dir)
    router := gin.Default()
    router.GET("/word/:term", handleLookup)
    router.GET("/suggestions/:term", handleSuggestions)
    router.Run(":8080")
}

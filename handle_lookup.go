package main

import(
    "net/http"
    "github.com/gin-gonic/gin"
)

func handleLookup(ctx *gin.Context) {
    term := ctx.Param("term")
    if term == "" {
        return
    }
    
    res := &Res{
        Term : term,
    }

    defs := []*Def{}
    for _, dict := range dictMap {
         name := dict.Name()
         desc := dict.Lookup(term)
         definition := &Def{
            Dict : name,
            Desc : desc,
         }
         defs = append(defs, definition)
    }

    res.Defs = defs

    ctx.JSON(http.StatusOK, res)
}

func handleSuggestions(ctx *gin.Context) {
    term := ctx.Param("term")
    if term == "" {
        return
    }
    res := &Res{
        Term : term,
    }

    maxCnt := 5
    sugs := []string{}
    filter := make(map[string]bool)
    for _, dict := range dictMap {
         suggestions := dict.Suggest(term, maxCnt)
         for _, suggestion := range suggestions {
            if _, ok := filter[suggestion]; !ok {
                sugs = append(sugs, suggestion)
                filter[suggestion] = true
            }
         }
    }

    length := len(sugs)
    if maxCnt < length {
        length = maxCnt
    }
    res.Sugs = sugs[:length]

    ctx.JSON(http.StatusOK, res)
}

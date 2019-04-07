package main

import (
    "path"
    "strings"
    "io/ioutil"
)

func loadDictsFromDir(dir string) map[string]Dict {
    dictMap = make(map[string]Dict)

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return dictMap
    }

    for _, file := range files {
        if file.IsDir() {
            sub := path.Join(path.Dir(dir), file.Name()) 
            typeName := GuessDictType(sub)
            dt := DictFactory(typeName, sub)
            //fmt.Printf("fi %s, name %s sub %s\n",file.Name(), typeName, sub)
            if dt != nil {
                //fmt.Printf("load dict %s dir %s\n", dt.Name(), sub)
                dictMap[dt.Name()] = dt
            }
        }
    }

    return dictMap
}

func DictFactory(typeName, dir string) Dict{
    var dt Dict
    switch typeName {
        case "stardict" : 
            dt = NewStarDict(dir)
    }

    return dt
}

func GuessDictType(dir string) string {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return ""
    }

    names := []string{}
    for _, file := range files {
        names = append(names, file.Name())
    }

    for dictName, f := range dictGuess {
        if guess := f(names); guess {
            return dictName
        }
    }

    return ""
}

type GuessFunc func([]string) bool

func registerDictGuesses() {
    dictGuess = make(map[string]GuessFunc)
    dictGuess["stardict"] = func(files []string) bool {
        hit := false
        for _, file := range files {
            if strings.HasSuffix(file, ".ifo") {
                //name = strings.Split(file, ".")[0]
                hit = true
            }
        }

       return hit
    }
}

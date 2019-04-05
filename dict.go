package main

import (
    "fmt"
    "os"
    "strconv"
    "bufio"
    "strings"
    "reflect"
    "encoding/json"
)

type Dict interface {
    Name() string
    Load(string) error
    Lookup(string)string
    String()string
}

type MetaInfo struct {
    Name string         `json:"bookname"`
    Desc string         `json:"description"` 
    Version string      `json:"version"`
    Author string       `json:"author"`  
    Email string        `json:"email"`
    Copyright string    `json:"copyright"`
    Website string      `json:"website"`
    Date string         `json:"date"`
    WordCount uint32    `json:"wordcount"`
    Synwordcount uint32 `json:"synwordcount"`
    Idxfilesize uint32  `json:"idxfilesize"`   
    Idxoffsetbits uint32    `json:"idxoffsetbits"`
    Sametypesequence string `json:"sametypesequence"`
    Dicttype uint32     `json:"dicttype"`
}

type StarDict struct {
   Meta *MetaInfo 
}

func NewStarDict(base string) *StarDict {
    d := &StarDict{
    }

    d.Load(base)

    return d
}

func (d *StarDict) Name() string {
    if d.Meta != nil {
        return d.Meta.Name
    }

    return ""
}

func CheckIfExist(file, backup string) bool {
    if _, err := os.Stat(file); !os.IsNotExist(err) {
        return true
    }

    if backup != "" {
        if _, err := os.Stat(backup); !os.IsNotExist(err) {
            return true
        }
    }

    return false
}

func (d *StarDict) String() string {
    str := ""
    if d.Meta != nil {
       b, _ := json.Marshal(d.Meta) 
       str = string(b)
    }

    return str
}

func (d *StarDict) Load(base string) error {
    ifo := base + ".ifo"
    idx := base + ".idx"
    dt := base + ".dict"

    if !CheckIfExist(ifo, "") || !CheckIfExist(idx, "") || !CheckIfExist(dt, dt + ".dz") {
        fmt.Printf("file not enough info %s, idx %s, dt %s\n", ifo, idx, dt)
        return nil
    }

    meta, err := d.loadMeta(ifo)
    d.Meta = meta

    return err
}

func (d *StarDict) loadMeta(file string) (*MetaInfo, error) {
    fr, err := os.Open(file)
    if err != nil {
        fmt.Printf("no file %s, exist\n", file)
        return nil, err 
    }
    defer fr.Close()

    metaMap := make(map[string]string)
    scanner := bufio.NewScanner(fr)
    for scanner.Scan() {
        text := scanner.Text()
        parts := strings.Split(text, "=")
        if len(parts) != 2 {
            continue
        }
        //fmt.Printf("%s %s\n", parts[0], parts[1])
        metaMap[parts[0]] = parts[1]
    }

    meta := MetaInfo{}
    ref_type := reflect.TypeOf(meta)
    ref_value := reflect.ValueOf(&meta)

    for i:= 0; i < ref_type.NumField(); i++ {
        field := ref_type.Field(i)
        tag := field.Tag.Get("json")
        if v, ok := metaMap[tag]; ok {
            switch field.Type.Kind() {
                case reflect.Uint32:
                    value, _ := strconv.ParseUint(v,10, 32)
                    fmt.Printf("tag %s , value %s %d\n", tag, v, value)
                    ref_value.Elem().FieldByName(field.Name).SetUint(value)
                case reflect.String:
                    ref_value.Elem().FieldByName(field.Name).SetString(v)
            }
        }
    }

    return &meta, nil
}

func (d *StarDict) Lookup(term string) string {
    return ""
}

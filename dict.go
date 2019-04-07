package main

import (
    "fmt"
    "os"
    "strconv"
    "bufio"
    "io"
    "path"
    "io/ioutil"
    "strings"
    "reflect"
    "encoding/json"
    "encoding/binary"
    "github.com/armon/go-radix"
)

type Dict interface {
    Name() string
    Load(string) error
    Lookup(string)string
    Suggest(string, int)[]string
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
   dict *radix.Tree
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

func (d *StarDict) GenDictName(base string) string {
    files, err := ioutil.ReadDir(base)
    if err != nil {
        return ""
    }

    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".ifo") {
            return strings.TrimSuffix(file.Name(), ".ifo")
        }
    }

    return ""
}

func (d *StarDict) Load(base string) error {
    //ifo := base + ".ifo"
    //idx := base + ".idx"
    //dt := base + ".dict"
    name := d.GenDictName(base)
    ifo := path.Join(base, name + ".ifo")
    idx := path.Join(base, name + ".idx")
    dt := path.Join(base, name + ".dict")

    if !CheckIfExist(ifo, "") || !CheckIfExist(idx, "") || !CheckIfExist(dt, dt + ".dz") {
        fmt.Printf("file not enough info %s, idx %s, dt %s\n", ifo, idx, dt)
        return nil
    }

    meta, err := d.loadMeta(ifo)
    if err != nil {
        return err
    }

    d.Meta = meta

    dict, err := d.loadDict(idx, dt)
    if err != nil {
        return err
    }

    d.dict = dict

    return nil
}

func (d *StarDict) loadDict(idx, data string) (*radix.Tree, error) {
    f_idx, err := os.Open(idx)
    if err != nil {
        return nil, err
    }
    defer f_idx.Close()

    f_data, err := os.Open(data)
    if err != nil {
        return nil, err
    }
    defer f_data.Close()

    r_idx := bufio.NewReader(f_idx)
    r_data := bufio.NewReader(f_data)

    tree := radix.New()

    //i := 0
    for {
        b, err := r_idx.ReadBytes(byte(0))
        if err != nil {
            //fmt.Printf("i %d b %v %v \n", i, b, err)
            break
        }
        w := string(b[:len(b) - 1])
        start, length := uint32(0), uint32(0)
        err = binary.Read(r_idx, binary.BigEndian, &start)
        err = binary.Read(r_idx, binary.BigEndian, &length)
        piece := make([]byte, length)

        //make sure read the exact bytes
        n, err := io.ReadFull(r_data, piece)
        if err != nil || uint32(n) < length {
            fmt.Printf("read err %v n%d, word %s, start %d len %d \n", err, n, w, start, length)
            break
        }
        //fmt.Printf("---read word %s, start %d len %d desc %s \n", w, start, length, string(piece))
        tree.Insert(w, string(piece))
    }

    return tree, nil
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
                    ref_value.Elem().FieldByName(field.Name).SetUint(value)
                case reflect.String:
                    ref_value.Elem().FieldByName(field.Name).SetString(v)
            }
        }
    }

    return &meta, nil
}

func (d *StarDict) Lookup(term string) string {
    if d.dict != nil {
        if desc, ok := d.dict.Get(term); ok {
            return desc.(string)
        }
    }

    return ""
}

func (d *StarDict) Suggest (term string, maxCnt int) []string {
    cnt := 0
    suggestions := []string{}
    if maxCnt == 0 {
        return suggestions
    }
    if d.dict != nil {
        d.dict.WalkPrefix(term, func(s string, v interface{}) bool {
            suggestions = append(suggestions, s)
            cnt += 1
            if cnt < maxCnt {
                return false
            } else {
                return true
            }
        })
    }

    return suggestions
}

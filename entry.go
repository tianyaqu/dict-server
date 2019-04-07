package main

type Entry struct {
    Desc string         `json:"desc"`
    Example string      `json:"example"`
}

type Morph struct {
    Name string         `json:"name"`
    Entries []*Entry    `json:"entries"`
}

type Word struct {
    Morphs []*Morph     `json:"morphs"`
}

func NewWord(str string) *Word {
    return nil
}

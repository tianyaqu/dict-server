package main

type Def struct {
    Dict string     `json:"dict"`
    Desc string     `json:"desc"`
}

type Res struct {
    Term string     `json:"term"`
    Defs []*Def     `json:"definition,omitempty"`
}

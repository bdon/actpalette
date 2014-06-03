package main

import (
  "encoding/json"
  "encoding/hex"
  "encoding/binary"
  "os"
  "io/ioutil"
  "bytes"
  "fmt"
  "path/filepath"
)

func main() {
  input_file := os.Args[1]
  s, err := ioutil.ReadFile(input_file)
  if err != nil {
    fmt.Println(err)
  }
  p := NewPalette(s)
  var extension = filepath.Ext(input_file)
  var name = input_file[0:len(input_file)-len(extension)]
  ioutil.WriteFile(name + ".act",p.Bytes(),0777)
}

type Palette struct {
  Colors []string `json:"colors"`
  TransparentColor string `json:"transparent_color"`
}

func NewPalette(str []byte) Palette {
  p := Palette{}
  json.Unmarshal(str, &p)
  return p
}

func (p Palette) IndexOfTColor() int {
  for i, c := range p.Colors {
    if c == p.TransparentColor {
      return i
    }
  }
  panic("transparent color not in colors")
}

func (p Palette) Bytes() []byte {
  retval := make([]byte, 772)
  for i := 0; i < 772; i++ {
    retval[i] = 0
  }

  for i, color := range p.Colors {
    s, _ := hex.DecodeString(color)
    retval[3*i] = s[0]
    retval[3*i+1] = s[1]
    retval[3*i+2] = s[2]
  }

  cbytes := new(bytes.Buffer)
  binary.Write(cbytes, binary.BigEndian, uint64(len(p.Colors)))

  // length 
  cbytes.Next(6)
  b1, _ := cbytes.ReadByte()
  retval[768] = b1
  b2, _ := cbytes.ReadByte()
  retval[769] = b2

  if p.TransparentColor != "" {
    tbytes := new(bytes.Buffer)
    err := binary.Write(tbytes, binary.BigEndian, uint64(p.IndexOfTColor()))
    if (err != nil) {
      fmt.Println(err)
    }
    // length 
    tbytes.Next(6)
    b3, _ := tbytes.ReadByte()
    retval[770] = b3
    b4, _ := tbytes.ReadByte()
    retval[771] = b4
  } else {
    // transparency
    retval[770] = 0xFF
    retval[771] = 0xFF
  }

  return retval
}


package main

import (
  "encoding/json"
  "encoding/hex"
  "encoding/binary"
)

func main() {

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

  bytes := make([]byte,2)
  binary.PutUvarint(bytes, uint64(len(p.Colors)))

  // length 
  retval[768] = bytes[0]
  retval[769] = bytes[1]

  if p.TransparentColor != "" {
    tbytes := make([]byte,2)
    binary.PutUvarint(tbytes, uint64(p.IndexOfTColor()))

    // length 
    retval[770] = tbytes[0]
    retval[771] = tbytes[1]
  } else {
    // transparency
    retval[770] = 0xFF
    retval[771] = 0xFF
  }

  return retval
}


package main

import (
  "bytes"
  "testing"
  "encoding/binary"
)

func TestLen(t *testing.T) {
  str := []byte("{}")
	a := NewPalette(str)
  r := bytes.NewReader(a.Bytes())
  if r.Len() != 772 {
    t.Error("palette should have fixed size of 772");
  }
}

func TestSetColors(t *testing.T) {
  str := []byte(`{
    "colors":["FFEEDD","CCBBAA"],
    "transparent_color":"FFEEDD"
    }`)
	a := NewPalette(str)
  r := bytes.NewReader(a.Bytes())
  r.Seek(0,0)
  b, _ := r.ReadByte()
  if b != 0xFF {
    t.Error("first color should be full red");
  }
  r.Seek(1,0)
  b, _ = r.ReadByte()
  if b != 0xEE {
    t.Error("first color should be EE green");
  }
  r.Seek(2,0)
  b, _ = r.ReadByte()
  if b != 0xDD {
    t.Error("first color should be DD blue");
  }
  r.Seek(3,0)
  b, _ = r.ReadByte()
  if b != 0xCC {
    t.Error("second color should be CC red");
  }
}

func TestNoTransparent(t *testing.T) {
  str := []byte(`{
    "colors":["FFEEDD","CCBBAA","000000"]
    }`)
	a := NewPalette(str)
  r := bytes.NewReader(a.Bytes())
  r.Seek(770,0)
  bnum := make([]byte,2)
  b1, _ := r.ReadByte()
  bnum[0] = b1
  b2, _ := r.ReadByte()
  bnum[1] = b2

  if bnum[0] != 0xFF || bnum[1] != 0xFF {
    t.Errorf("transparent index should be FFFF, got %d", bnum);
  }
}

func TestSetTransparent(t *testing.T) {
  str := []byte(`{
    "colors":["FFEEDD","CCBBAA","000000"],
    "transparent_color":"CCBBAA"
    }`)
	a := NewPalette(str)
  r := bytes.NewReader(a.Bytes())
  r.Seek(770,0)
  bnum := make([]byte,2)
  b1, _ := r.ReadByte()
  bnum[0] = b1
  b2, _ := r.ReadByte()
  bnum[1] = b2

  buf := bytes.NewReader(bnum)
  num, _ := binary.ReadUvarint(buf)
  if num != 1 {
    t.Errorf("transparent index should be 1, got %d", num);
  }
}

func TestSetPaletteLength(t *testing.T) {
  str := []byte(`{
    "colors":["FFEEDD","CCBBAA","F0F0F0"],
    "transparent_color":"F0F0F0"
    }`)
	a := NewPalette(str)
  r := bytes.NewReader(a.Bytes())
  r.Seek(768,0)
  bnum := make([]byte,2)
  b1, _ := r.ReadByte()
  bnum[0] = b1
  b2, _ := r.ReadByte()
  bnum[1] = b2

  buf := bytes.NewReader(bnum)
  num, _ := binary.ReadUvarint(buf)
  if num != 3 {
    t.Errorf("palette length should be 3, got %d", num);
  }

  // test max number of colors (255)
  bigarr := []string{}
  for i := 0; i < 255; i++ {
    bigarr = append(bigarr, "F0F0F0")
  }
  a.Colors = bigarr
  r = bytes.NewReader(a.Bytes())
  r.Seek(768,0)
  bnum = make([]byte,2)
  b1, _ = r.ReadByte()
  bnum[0] = b1
  b2, _ = r.ReadByte()
  bnum[1] = b2

  buf = bytes.NewReader(bnum)
  num, _ = binary.ReadUvarint(buf)
  if num != 255 {
    t.Errorf("palette length should be 255, got %d", num);
  }
}

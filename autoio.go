package autoio

import "os"
import "io"
import "bufio"

import "errors"

import "compress/gzip"
import "compress/bzip2"


// Wrap common stream file types into one for ease of scanning
//
type AutoioHandle struct {
  Fp *os.File
  Scanner *bufio.Scanner
  Writer *bufio.Writer

  Bz2Reader io.Reader
  GzReader *gzip.Reader
  FileType string
}


// Magic strings we look for at the beginning of the file to determine file type.
//
var magicmap map[string]string = map[string]string{ "\x1f\x8b" : ".gz" , "\x1f\x9d" : ".Z", "\x42\x5a" : ".bz2" , "\x50\x4b\x03\x04" : ".zip" }

func OpenScanner( fn string ) ( h AutoioHandle, err error ) {

  if fn == "-" {
    h.Fp = os.Stdin
    h.Scanner = bufio.NewScanner( h.Fp )
    return h, nil
  }

  var sentinalfp *os.File

  sentinalfp,err = os.Open( fn )
  if err != nil { return h, err }
  defer sentinalfp.Close()


  b := make( []byte, 2, 2 )
  n,err := sentinalfp.Read(b)
  if (n<2) || (err != nil) {
    h.Fp,err = os.Open( fn )
    if err != nil { return h, err }
    h.Scanner = bufio.NewScanner( h.Fp )
    return h, err
  }

  if typ,ok := magicmap[string(b)] ; ok {

    h.Fp,err = os.Open( fn )
    if err != nil { return h, err }

    if typ == ".gz" {
      h.FileType = "gz"

      h.GzReader,err = gzip.NewReader( h.Fp )
      if err != nil {
        h.Fp.Close()
        return h, err
      }
      h.Scanner = bufio.NewScanner( h.GzReader )
    } else if typ == ".bz2" {

      h.FileType = "bz2"

      h.Bz2Reader = bzip2.NewReader( h.Fp )
      h.Scanner = bufio.NewScanner( h.Bz2Reader )
    } else {
      err = errors.New(typ + "extension not supported")
    }

    return h, err
  }


  b2 := make( []byte, 2, 2)
  n,err = sentinalfp.Read(b2)
  if (n<2) || (err != nil) {
    h.Fp,err = os.Open( fn )
    if err != nil { return h, err }
    h.Scanner = bufio.NewScanner( h.Fp )
    return h, err
  }

  s := string(b) + string(b2)
  if typ,ok := magicmap[s]; ok {
    if typ == ".zip" {
      err = errors.New("zip extension not supported")
      return h, err
    }
    err = errors.New(typ + "extension not supported")
    return h, err
  }

  h.Fp,err = os.Open( fn )
  if err != nil { return h, err }
  h.Scanner = bufio.NewScanner( h.Fp )

  return h, err

}

func CreateWriter( fn string ) ( h AutoioHandle, err error ) {

  if fn == "-" {
    h.Fp = os.Stdout
  } else {
    h.Fp,err = os.Create( fn )
    if err != nil { return h, err }
  }

  h.Writer = bufio.NewWriter( h.Fp )
  return h, nil
}

func (h *AutoioHandle) Flush() {

  if h.Writer != nil {
    h.Writer.Flush()
  }

}

func ( h *AutoioHandle) Close() error {

  if h.FileType == "gz" {
    e := h.GzReader.Close()
    if e!=nil { return e }
  }

  e := h.Fp.Close()
  return e

}

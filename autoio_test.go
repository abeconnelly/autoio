package autoio

import "fmt"
import "testing"

var g_verbose bool = false


func TestOpenTextFile( t *testing.T ) {
  fn := "./magic/test.txt"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}

func TestOpenGzFile( t *testing.T ) {
  fn := "./magic/test.txt.gz"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}

func TestOpenBzip2File( t *testing.T ) {
  fn := "./magic/test.txt.bz2"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}



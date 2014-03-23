package gosk

import (
	"fmt"
	"testing"
)

const format string = `
return printf("%s%s%s", color, body, reset)
  color =
    switch .Typ
    case "text":
      <- $blue
    case "heading":
      if .IsDir
        <- $red
      else
        <- $green
    default:
      <- $yellow

  body =
    switch .Typ
    case "text":
      <- printf("text: %s", .Data)
    case "heading":
      <- printf("heading: %s", .Data)
    default:
      <- printf("%d", .Size)
`

const format2 = `
return "format2"
`

const format3 = `
return printf("%s %s", d, d2)
  d =
    <- .Data
  d2 =
    <- .Data [1:]
`

const format4 = `
return head + body_color + body
  head = "head"[1:]
  body_color = "color"
  body = "body"
`

const format5 = `
switch .typ
case "text":
  <- "text"
case "heading":
  if 4 <= 1
    <- "foo"
  else
    <- "heading"
default:
  <- "default"
`

type line struct {
	IsDir bool
	Typ   string
	Data  string
	Size  int
}

func TestGosk(t *testing.T) {
	// yyDebug = 3
	tmpl := NewTemplate(format).Vars(VarMap{
		"printf": fmt.Sprintf,
		"reset":  "rst",
		"blue":   "b",
		"red":    "r",
		"green":  "g",
		"yellow": "y",
	})
	tmpl2 := NewTemplate(format2)

	tmpl3 := NewTemplate(format3).Vars(VarMap{
		"printf": fmt.Sprintf,
	})

	tmpl4 := NewTemplate(format4).Vars(VarMap{
		"printf": fmt.Sprintf,
	})

	tmpl5 := NewTemplate(format5)

	// rootNode.Debug()

	line1 := &line{
		IsDir: true,
		Data:  "data1",
		Size:  14,
	}
	line2 := VarMap{
		"IsDir": false,
		"Typ":   "heading",
		"Data":  "data2",
		"Size":  16,
	}
	if tmpl.Execute(line1) != "y14rst" {
		t.Error("unexpected")
	}
	if tmpl.Execute(line2) != "gheading: data2rst" {
		t.Error("unexpected")
	}

	if tmpl2.Execute() != "format2" {
		t.Error("unexpected")
	}

	if tmpl3.Execute(line2) != "data2 ata2" {
		t.Error("unexpected")
	}

	if tmpl4.Execute() != "eadcolorbody" {
		t.Error("unexpected")
	}

	if tmpl5.Execute(map[string]string{"typ": "heading"}) != "heading" {
		t.Error("unexpected")
	}
}

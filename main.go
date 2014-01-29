package main

import (
	"fmt"
)

const format string = `
return printf("%s%s%s", color, body, reset)
  color =
    switch .type
    case "text":
      <- blue
    case "heading":
      if .is_dir
        <- red
      else
        <- green
    default:
      <- yellow

  body =
    switch .type
    case "text":
      <- printf("text: %s", .data)
    case "heading":
      <- printf("heading: %s", .data)
    default:
      <- printf("%d", .size)
`

const format2 = `
return printf("d: %d, d2: %d", .d, d2)
`

func main() {
	// yyDebug = 3
	t := NewTemplate(format2).Vars(VarMap{
		"printf": fmt.Sprintf,
		"d2":     5,
	})
	// rootNode.Debug()
	fmt.Println(t.Execute(VarMap{
		"d": 4,
	}))
}

package main

import (
	"fmt"
)

type CCC struct {
	val string
}

type someStr struct {
	value string
	ccc   *CCC
}

func main() {
	fmt.Println("MAP")
	store := make(map[string]*someStr)

	store["qwe"] = &someStr{
		value: "HELLO",
		ccc:   &CCC{val: "this is ccc"},
	}

	fmt.Println("VAL", store["qwe"], store["qwe"].ccc)

	val := store["qwe"]

	val.value = "BY"
	val.ccc.val = "another ccc"
	fmt.Println("VAL BY:", val, store["qwe"], store["qwe"].ccc)

	val.ccc.val = "another ccc"
	ccc := val.ccc
	ccc.val = "FINISH"
	fmt.Println("VAL CCC:", store["qwe"].ccc, ccc)

}

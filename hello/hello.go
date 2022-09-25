/**
 * @Author: Noaghzil
 * @Date:   2022-09-25 15:53:46
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2022-09-25 16:26:56
 */
package main

import (
	"fmt"

	"example.com/greetings"
)

func main() {
	message := greetings.Hello("Emaly")
	fmt.Println(message)
}

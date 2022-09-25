/**
 * @Author: Noaghzil
 * @Date:   2022-09-25 15:39:04
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2022-09-25 16:33:48
 */
package greetings

import "fmt"

func Hello(name string) string {
	message := fmt.Sprintf("Hello, %v. Welcome!", name)
	return message
}

package helper

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

// Centralize error handing, simple print message and exit
func LogError(err error, message string) {
	fmt.Println()
	fmt.Println()
	fmt.Println(color.RedString("Error: " + message))
	log.Fatalln(err)
}

// Warning but dont exit
func LogWarning(err error, message string) {
	fmt.Println()
	fmt.Println()
	fmt.Println(color.YellowString("Error: " + message))
}

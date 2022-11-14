package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
)

type Out struct {
	printer *message.Printer
}

func NewOut(tag language.Tag) Out { return Out{printer: message.NewPrinter(tag)} }

func (out *Out) Print(a ...any) {
	if _, err := out.printer.Print(a...); err != nil {
		log.Printf("unable to print message: '%s'", fmt.Sprint(a...))
		fmt.Print(a...)
	}
}
func (out *Out) Printf(format string, a ...any) {
	if _, err := out.printer.Printf(format, a...); err != nil {
		log.Printf("unable to print message: '%s'", fmt.Sprintf(format, a...))
		fmt.Printf(format, a...)
	}
}
func (out *Out) Println(a ...any) {
	if _, err := out.printer.Println(a...); err != nil {
		log.Printf("unable to print message: '%s'", fmt.Sprintln(a...))
		fmt.Println(a...)
	}
}

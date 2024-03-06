package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"syscall/js"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	fmt.Println("Go Web Assembly")

	js.Global().Set("formatJSON", jsonWrapper())
	js.Global().Set("ls", listFTP())

	<-make(chan struct{})

}

func listFTP() string {
	connect()

	ls := listDir(".")
	fmt.Println(ls)
	return ls
}

var _c *ftp.ServerConn

func connect() {
	var c *ftp.ServerConn
	c, err := ftp.Dial("localhost:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login("mrt", "23")
	if err != nil {
		log.Fatal(err)
	}

	// if err := c.Quit(); err != nil {
	// 	log.Fatal(err)
	// }
	_c = c
}

func storeFile(file string) {
	data := bytes.NewBufferString("Hello World")
	err := _c.Stor("test-file.txt", data)
	if err != nil {
		panic(err)
	}
}

func getFile(file string) string {
	r, err := _c.Retr("test-file.txt")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	println(string(buf))

	return string(buf)
}

func listDir(dir string) string {
	e, err := _c.List(".")
	if err != nil {
		return "!ERROR!" + err.Error()
	}

	var stringResult = ""
	for _, entry := range e {
		stringResult += entry.Name + "\n"
	}

	return stringResult
}

func deleteFile(file string) {
	fmt.Println("Deleting file", file)
}

func deleteDir(dir string) {
	fmt.Println("Deleting dir", dir)
}

func createDir(dir string) {
	fmt.Println("Creating dir", dir)
}

func renameFile(oldName string, newName string) {
	fmt.Println("Renaming file", oldName, "to", newName)
}

func renameDir(oldName string, newName string) {
	fmt.Println("Renaming dir", oldName, "to", newName)
}

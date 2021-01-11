package main

import (
	"JsonChecker/lib"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filePath := flag.String("file", "./correctJson/correct.json", "path to correct json file")
	flag.Parse()
	fmt.Printf("file path: %s\n", *filePath)

	targetJsons := flag.Args()
	var targetJson string
	if len(targetJsons) >= 3 { // フラグの他に入力があればtargetが入力されている
		targetJson = targetJsons[0]
	} else {
		fmt.Println("input target json...")
		targetJsonByte, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("error has occured at stdin")
			return
		}
		targetJson = string(targetJsonByte)
	}
	if targetJson == "" {
		fmt.Printf("target json is not found")
		return
	}
	fmt.Printf("targetJson: %s\n", targetJson)

	if err := lib.JsonCheck(targetJson, *filePath); err != nil {
		fmt.Printf("error has occurred: %s\n", err.Error())
		return
	}

	return
}

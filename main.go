package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening to port :6379")

	listen, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listen.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)

		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Value must be an array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Array length cannot be 0")
			continue
		}

		writer := NewWriter(conn)

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		commandFunc, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command:", command)
			writer.Write(Value{typ: "error", str: "Invalid Command"})
			continue
		}

		response := commandFunc(args)

		writer.Write(response)
	}

}

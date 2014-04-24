# Clickatell apis for Go

## SMS

### Usage
```go
package main

import "github.com/gabstv/clickatellgo/sms"
import "fmt"

func main(){
	//// method 1
	auth := sms.Auth{"user", "password", "apiid"}
	m := sms.New(auth)
	// send
	msgid, err := m.Send("17775550000", "Hello, world.")
	fmt.Println(msgid, err)

	//// method 2
	msgid, err = sms.SendSMS("user", "password", "apiid", "17775550000", "Hello, world.")
	fmt.Println(msgid, err)
}
```
package sms

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Auth struct {
	User     string
	Password string
	APIID    string
}

type sender struct {
	Auth Auth
}

func (s *sender) Send(phonenumber, message string) (id string, err error) {
	id, err = SendSMS(s.Auth.User, s.Auth.Password, s.Auth.APIID, phonenumber, message)
	return
}

func New(auth Auth) *sender {
	output := &sender{auth}
	return output
}

func SendSMS(user, password, apiid, phonenumber, message string) (id string, err error) {
	functimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(30*time.Second))
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            functimeout,
	}
	client := &http.Client{Transport: tr}
	dv0 := url.Values{"user": {user}, "password": {password}, "api_id": {apiid}, "to": {phonenumber}, "text": {message}}
	buff := new(bytes.Buffer)
	wc, err := charset.NewWriter("ascii", buff)
	if err != nil {
		fmt.Println("ERROR ascii (latin1)", err)
		return
	}
	wc.Write([]byte(dv0.Encode()))
	wc.Close()
	r0, err := http.NewRequest("POST", "https://api.clickatell.com/http/sendmsg", buff)
	if err != nil {
		return
	}
	r0.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=ascii")
	//resp, err := client.PostForm("https://api.clickatell.com/http/sendmsg", url.Values{"user": {user}, "password": {password}, "api_id": {apiid}, "to": {phonenumber}, "text": {message}})
	resp, err := client.Do(r0)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	defer resp.Body.Close()
	io.Copy(&buffer, resp.Body)
	respstr := buffer.String()
	if strings.HasPrefix(respstr, "ID: ") {
		id = respstr[4:]
	} else {
		err = errors.New(respstr)
	}
	return
}

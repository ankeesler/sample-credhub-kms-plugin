package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pivotal/sample-credhub-kms-plugin/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func die(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 4 {
		die(fmt.Sprintf("usage: %s <socket> <method> <data>", os.Args[0]))
	}

	config := tls.Config{InsecureSkipVerify: true}
	creds := credentials.NewTLS(&config)
	conn, err := grpc.Dial(
		os.Args[1],
		grpc.WithTransportCredentials(creds),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		die(fmt.Sprintf("dial: %s", err.Error()))
	}
	defer conn.Close()

	c := v1beta1.NewKeyManagementServiceClient(conn)
	ctx := context.Background()
	switch os.Args[2] {
	case "encrypt":
		req := v1beta1.EncryptRequest{
			Plain: []byte(os.Args[3]),
		}
		rsp, err := c.Encrypt(ctx, &req)
		if err != nil {
			die(fmt.Sprintf("encrypt: %s", err.Error()))
		}

		//fmt.Println("encrypt:", base64.StdEncoding.EncodeToString(rsp.Cipher))
		fmt.Println("encrypt:", string(rsp.Cipher))

	case "decrypt":
		req := v1beta1.DecryptRequest{
			Cipher: []byte(os.Args[3]),
		}
		//_, err := base64.StdEncoding.Decode(req.Cipher, []byte(os.Args[3]))
		//if err != nil {
		//	die(fmt.Sprintf("decode: %s", err.Error()))
		//}

		rsp, err := c.Decrypt(ctx, &req)
		if err != nil {
			die(fmt.Sprintf("decrypt: %s", err.Error()))
		}

		fmt.Println("decrypt:", string(rsp.Plain))
	}
}

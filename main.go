package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"net"
	"log"
	"bytes"
	"bufio"
	"strings"
	"io/ioutil"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	sh *shell.Shell
	data []byte
	gateways = []string{
		"http://localhost:8080/ipfs/",
		"https://gateway.ipfs.io/ipfs/",
		"https://cloudflare-ipfs.com/ipfs/",
	}
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDaemon(daemon string) *shell.Shell {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", daemon, timeout)
	HandleErr(err)

	if conn != nil {
		defer conn.Close()
		fmt.Println("Connecting to ", daemon)
	}

	return shell.NewShell(daemon)
}

func stdinRead() {
	scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        b := scanner.Bytes()
		if bytes.Compare(b, []byte("EOF")) == 0 {
			return
		} else {
			data = append(data, b[:]...)
			data = append(data, []byte("\n")[:]...)
		}
    }
    if err := scanner.Err(); err != nil {
        HandleErr(err)
    }
}

func main() {
	daemon := flag.String("daemon", "localhost:5001", "Address and port of IPFS daemon.")
	eof := flag.Bool("eof", false, "Create encrypted paste from STDIN.")
	fname := flag.String("input", "", "File to encrypt & upload.")
	password := flag.String("password", "", "Password/key for encryption. This will be SHA1'd.")

	sh = getDaemon(*daemon)

	flag.Parse()

	if *password == "" {
		log.Fatal("Must supply a password/key with --password")
	}

	if *fname != "" {
		f, err := os.Open(*fname)
		HandleErr(err)
		data, err = ioutil.ReadAll(f)
		HandleErr(err)
	} else if *eof {
		fmt.Println("Paste your message here & finish by typing EOF + enter")
		stdinRead()
	} else {
		log.Fatal("Must supply either paste text, or a file name with --eof or --input")
	}

	key := Key(*password)

	_ = fmt.Sprintf(*daemon)

	b64, err := Encrypt([]byte(key), data)
	HandleErr(err)

	var html []byte

	if *fname != "" {
		html = DecrypterFromFile(*fname, b64)
	} else if *eof {
		html = DecrypterFromPaste(b64)
	}

	cid, err := sh.Add(strings.NewReader(string(html)))
	HandleErr(err)

	node, err := sh.Resolve("")
	HandleErr(err)

	err = sh.Publish(node, cid)
	HandleErr(err)

	fmt.Println("GATEWAYS\n============================")

	for _, gateway := range gateways {
		fmt.Printf("%s%s#%s\n", gateway, cid, key)
	}
}
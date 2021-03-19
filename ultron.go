package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const POOL_ADDR string = "51.15.127.80:2811"
const MINER_NAME string = "ultronM1N3R"
const MINER_VERSION string = "1.0.1"
const RIG_ID string = "None"

var USERNAME string = "samedamci"

const PROXIES_FILE string = "proxies.txt"

var REJECTED, ACCEPTED int

const MAX_REJECTIONS = 100
const TIMEOUT = 1.5 // seconds

func worker(dialer proxy.Dialer, thread_index int) {
	conn, err := dialer.Dial("tcp", POOL_ADDR)
	if err != nil {
		return
	}
	buffer := make([]byte, 3)
	_, err = conn.Read(buffer)
	if err != nil {
		return
	}

	for {
		_, err = conn.Write([]byte("JOB," + USERNAME + ",AVR"))
		if err != nil {
			continue
		}

		buffer = make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			time.Sleep(1 * time.Second)
			fmt.Printf("error getting the job (%d)\n", thread_index)
			continue
		}

		job := strings.Split(string(buffer), ",")
		if len(job) != 3 {
			continue
		}
		prefix_bytes := job[0]
		target_bytes := job[1]
		diff, _ := strconv.Atoi(strings.Replace(job[2], "\x00", "", -1))

		// Calculate hash
		for i := 0; i <= diff*100; i++ {
			h := sha1.New()
			h.Write([]byte(prefix_bytes + strconv.Itoa(i)))
			sum := h.Sum(nil)
			hash := hex.EncodeToString(sum[:])

			if hash == target_bytes {
				time.Sleep(TIMEOUT * 1000000000 * time.Nanosecond)
				_, err = conn.Write(
					[]byte(
						strconv.Itoa(i) + ",," + MINER_NAME + " " + MINER_VERSION + "," + RIG_ID,
					),
				)
				if err != nil {
					fmt.Println(err)
				}

				buffer := make([]byte, 6)
				_, err = conn.Read(buffer)
				buffer = bytes.Trim(buffer, "\x00")
				feedback := string(buffer)

				if feedback == "GOOD" || feedback == "BLOCK" {
					ACCEPTED++
				} else if feedback == "BAD" {
					REJECTED++
					buffer = make([]byte, 1024)
					_, err = conn.Read(buffer)
				}
			}
		}
	}
}

func main() {
	if len(os.Args)-1 >= 1 {
		USERNAME = os.Args[1]
	}
	proxies_bytes, err := ioutil.ReadFile(PROXIES_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxies := strings.Split(string(bytes.Trim(proxies_bytes, "\n")), "\n")

	fmt.Printf(">> Miner: %s %s\n", MINER_NAME, MINER_VERSION)
	fmt.Printf(">> Identifier: %s\n", RIG_ID)
	fmt.Printf(">> Username: %s\n", USERNAME)
	fmt.Printf(">> Using %d SOCKS5 proxies\n", len(proxies))
	fmt.Printf(">> Pool: %s\n", POOL_ADDR)
	fmt.Println()
	fmt.Println("initializing")

	thread_index := 0
	for i := 0; i < len(proxies); i++ {
		proxy_addr := proxies[i]
		fmt.Printf("creating workers for %s proxy connection\n", proxy_addr)
		dialer, _ := proxy.SOCKS5("tcp", proxy_addr, nil, nil)

		for j := 0; j < 23; j++ {
			thread_index++
			go worker(dialer, thread_index)
			time.Sleep(100000000 * time.Nanosecond) // 0.1 ms
		}
	}
	fmt.Printf("created %d workers\n", thread_index)
	fmt.Println("initialized")

	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("shares: (%d/%d)\n", ACCEPTED, REJECTED)
		if int32((REJECTED/(ACCEPTED+1))*100) >= 10 || REJECTED >= MAX_REJECTIONS {
			break
		}
	}
}

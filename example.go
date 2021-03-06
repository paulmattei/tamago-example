// https://github.com/f-secure-foundry/tamago-example
//
// Copyright (c) F-Secure Corporation
// https://foundry.f-secure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

// Basic test example for tamago/arm running on supported i.MX6 targets.

package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	mathrand "math/rand"
	"os"
	"runtime"
	"time"

	"github.com/f-secure-foundry/tamago/soc/imx6"
)

var Build string
var Revision string
var banner string

var verbose = true

var exit chan bool

func init() {
	banner = fmt.Sprintf("%s/%s (%s) • %s %s",
		runtime.GOOS, runtime.GOARCH, runtime.Version(),
		Revision, Build)

	log.SetFlags(0)

	// imx6 package debugging
	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	model := imx6.Model()
	_, family, revMajor, revMinor := imx6.SiliconVersion()

	if !imx6.Native {
		banner += fmt.Sprintf(" • %s %d MHz (emulated)", model, imx6.ARMFreq()/1000000)
		return
	}

	if err := imx6.SetARMFreq(900); err != nil {
		log.Printf("WARNING: error setting ARM frequency: %v", err)
	}

	banner += fmt.Sprintf(" • %s %d MHz", model, imx6.ARMFreq()/1000000)

	log.Printf("imx6_soc: %s (%#x, %d.%d) @ %d MHz - native:%v",
		model, family, revMajor, revMinor, imx6.ARMFreq()/1000000, imx6.Native)
}

func example(init bool) {
	start := time.Now()
	exit = make(chan bool)
	n := 0

	log.Println("-- begin tests -------------------------------------------------------")

	n += 1
	go func() {
		log.Println("-- fs ----------------------------------------------------------------")
		TestFile()
		TestDir()

		exit <- true
	}()

	sleep := 100 * time.Millisecond

	n += 1
	go func() {
		log.Println("-- timer -------------------------------------------------------------")

		t := time.NewTimer(sleep)
		log.Printf("waking up timer after %v", sleep)

		start := time.Now()

		for now := range t.C {
			log.Printf("woke up at %d (%v)", now.Nanosecond(), now.Sub(start))
			break
		}

		exit <- true
	}()

	n += 1
	go func() {
		log.Println("-- sleep -------------------------------------------------------------")

		log.Printf("sleeping %s", sleep)
		start := time.Now()
		time.Sleep(sleep)
		log.Printf("slept %s (%v)", sleep, time.Since(start))

		exit <- true
	}()

	n += 1
	go func() {
		log.Println("-- rng ---------------------------------------------------------------")

		size := 32

		for i := 0; i < 10; i++ {
			rng := make([]byte, size)
			rand.Read(rng)
			log.Printf("%x", rng)
		}

		count := 1000
		start := time.Now()

		for i := 0; i < count; i++ {
			rng := make([]byte, size)
			rand.Read(rng)
		}

		log.Printf("retrieved %d random bytes in %s", size*count, time.Since(start))

		seed, _ := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
		mathrand.Seed(seed.Int64())

		exit <- true
	}()

	n += 1
	go func() {
		log.Println("-- ecdsa -------------------------------------------------------------")
		TestSignAndVerify()
		exit <- true
	}()

	n += 1
	go func() {
		log.Println("-- btc ---------------------------------------------------------------")

		ExamplePayToAddrScript()
		ExampleExtractPkScriptAddrs()
		ExampleSignTxOutput()

		exit <- true
	}()

	if imx6.Native && imx6.Family == imx6.IMX6ULL {
		n += 1
		go func() {
			log.Println("-- i.mx6 dcp ---------------------------------------------------------")
			TestDCP()
			exit <- true
		}()
	}

	log.Printf("launched %d test goroutines", n)

	for i := 1; i <= n; i++ {
		<-exit
	}

	log.Printf("----------------------------------------------------------------------")
	log.Printf("completed %d goroutines (%s)", n, time.Since(start))

	runs := 9
	chunksMax := 50
	chunks := mathrand.Intn(chunksMax) + 1
	fillSize := 160 * 1024 * 1024
	chunkSize := fillSize / chunks

	log.Printf("-- memory allocation (%d runs) ----------------------------------------", runs)
	testAlloc(runs, chunks, chunkSize)

	if imx6.Native {
		count := 10 * 1024 * 1024
		readSize := 0x7fff

		if init {
			// Pre-USB use the entire iRAM, accounting for required
			// alignments which take additional space.
			readSize = 0x20000 - 512
		}

		log.Println("-- memory cards -------------------------------------------------------")

		for _, card := range cards {
			TestUSDHC(card, count, readSize)
		}
	}
}

func main() {
	start := time.Now()

	log.Println(banner)

	example(true)

	if imx6.Native && (imx6.Family == imx6.IMX6UL || imx6.Family == imx6.IMX6ULL) {
		log.Println("-- i.mx6 usb ---------------------------------------------------------")
		StartUSB()
	}

	log.Printf("Goodbye from tamago/arm (%s)", time.Since(start))
}

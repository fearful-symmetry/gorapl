# goRAPL

[![GoDoc](https://godoc.org/github.com/fearful-symmetry/gorapl?status.svg)](https://godoc.org/github.com/fearful-symmetry/gorapl)  [![Go Report Card](https://goreportcard.com/badge/github.com/fearful-symmetry/gorapl)](https://goreportcard.com/report/github.com/fearful-symmetry/gorapl)

A dead-simple, low-level API for accessing the Intel [RAPL API](https://www.phoronix.com/scan.php?page=news_item&px=MTcxMjY).


## Install

`gorapl` uses modules for dependency management. This means you'll need at least go 1.11. With that in mind:

```
$ go get github.com/fearful-symmetry/gorapl
```

This is a rather experimental library that's early on in it's life, so it isn't very useful yet! It has a lot of shortcomings, and a lot of things that needs to do, but can't. Most notably, it can only handle systems with one CPU socket. TODOs and progress is being tracked on the issues page.

## Usage & caveats

Unless you're using something like [msr-safe](https://github.com/llnl/msr-safe), `gorapl` requres root to access the underlying msr device at `/dev/cpu/$CPU/msr`. 

If that device doesn't exist, you might need to load the kernel module:

```bash
$ sudo modprobe msr
```

Not all Intel CPUs support RAPL. The feature was introduced in Sandy Bridge, and available domains vary by processor. A quick way to check:

```bash
$ sudo rdmsr -a 0x611 #rdmsr is part of the msr-tools package on most linuxes

#You can check to see if the kernel picked up an RAPL domains too:

$ dmesg | grep rapl

```

Now that' you're ready to go, using `gorapl` is easy:

```go

    h, err := NewRAPL(0)
    if err != nil {
        // ...
    }

	dat, err := h.ReadPowerLimit(DRAM)
	if err != nil {
		//
	}
	fmt.Printf("Current RAPL power limit settings: %#v\n", dat)

```
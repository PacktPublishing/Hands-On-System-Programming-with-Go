package main

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/tatsushid/go-fastping"
)

const wait = time.Millisecond * 250

type pingPool chan *fastping.Pinger

func (p pingPool) Get() *fastping.Pinger {
	select {
	case v := <-p:
		return v
	case <-time.After(wait):
		return fastping.NewPinger()
	}
}

func (p pingPool) Put(v *fastping.Pinger) {
	select {
	case p <- v:
	case <-time.After(wait):
	}
	return
}

func help(ifaces []net.IP) {
	log.Println("please specify a valid network interface number")
	for i, f := range ifaces {
		mask, _ := f.DefaultMask().Size()
		fmt.Printf("%d - %s/%v\n", i, f, mask)
	}
	os.Exit(0)
}

func main() {
	ifaces, err := getInterfaces()
	if err != nil {
		log.Fatalln(err)
	}
	if len(os.Args) != 2 {
		help(ifaces)
	}
	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	if i < 0 || i > len(ifaces) {
		help(ifaces)
	}
	m := ifaces[i].DefaultMask()
	ip := ifaces[i].Mask(m)
	log.Printf("Lookup in %s", ip)
	done := make(chan struct{})
	address := make(chan net.IP)
	ones, bits := m.Size()
	pool := make(pingPool, 10)
	for i := 0; i < 1<<(uint(bits-ones)); i++ {
		go func(i int) {
			p := pool.Get()
			defer func() {
				pool.Put(p)
				done <- struct{}{}
			}()
			p.AddIPAddr(&net.IPAddr{IP: makeIP(ip, i)})
			p.OnRecv = func(a *net.IPAddr, _ time.Duration) { address <- a.IP }
			p.Run()
		}(i)
	}
i = 0
for {
    select {
    case ip := <-address:
        log.Printf("Found %s", ip)
    case <-done:
        if i >= bits-ones {
            return
        }
        i++
    }
}
}

func makeIP(ip net.IP, i int) net.IP {
	addr := make(net.IP, len(ip))
	copy(addr, ip)
	b := new(big.Int)
	b.SetInt64(int64(i))
	v := b.Bytes()
	copy(addr[len(addr)-len(v):], v)
	return addr
}

func getInterfaces() (result []net.IP, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip = ip.To4(); ip != nil {
				result = append(result, ip)
			}
		}
	}
	return
}

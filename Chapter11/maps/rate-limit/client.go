package main

import (
	"log"
	"os"
	"time"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch11/maps"
)

func main() {
	c := maps.NewClient(maps.DailyCap, os.Getenv("MAPS_APIKEY"))
	defer c.Close()
	start := time.Now()
	for _, l := range [][2]float64{
		{40.4216448, -3.6904040},
		{40.4163111, -3.7047328},
		{40.4123388, -3.7096724},
		{40.4145150, -3.7064412},
	} {
		locs, err := c.ReverseGeocode(l[0], l[1])
		e := time.Since(start)
		if err != nil {
			log.Println(e, l, err)
			continue
		}
		if len(locs) != 0 {
			locs = locs[:1]
		}
		log.Println(e, l, locs)
	}
}

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/2qif49lt/dhtlistener"
	flag "github.com/2qif49lt/pflag"
	"net/http"
	_ "net/http/pprof"
)

var srvaddr = flag.StringP("addr", "a", "", "address ip:port")

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func main() {

	flag.Parse()
	if *srvaddr == "" {
		*srvaddr = ":0"
	}
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	w := dhtlistener.NewWire(1024, 256)
	go func() {
		for resp := range w.Response() {
			info := map[string]interface{}{}

			err := dhtlistener.Decode(resp.MetadataInfo, &info)
			if err != nil {
				continue
			}

			if _, ok := info["name"]; !ok {
				continue
			}

			bt := bitTorrent{
				InfoHash: hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}

			data, err := json.Marshal(bt)
			if err == nil {
				fmt.Printf("%s\n\n", data)
			}
		}
	}()
	go w.Run()

	d := dhtlistener.NewDht(*srvaddr)

	d.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}

	d.Run()
}

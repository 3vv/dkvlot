// h 20181127
//
// Redis server/cluster synchronizer

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/redis.v5"
)

func main() {
	for {
		//
		// Command line handling
		//
		// Init flags
		var keys = flag.Bool("keys", false, keysD)
		var sync = flag.Bool("sync", false, syncD)
		var l = flag.String("l", "", lD)
		var f = flag.String("f", "*", fD)
		var s = flag.String("s", "", sD)
		var d = flag.String("d", "", dD)
		// Parse flags
		flag.Parse()
		// Ensure Cmds
		if *keys == false && *sync == false {
			help()
			break
		}
		//
		// Redis server/cluster connecting
		//
		// Connect to the source host
		if len(*s) > 0 {
			sHosts = strings.Split(*s, ",")
			//fmt.Println(sHosts)
			connectSourceHost()
		} else {
			help()
			break
		}
		if !sConnected {
			help()
			break
		}
		// Connect to the destination host
		if len(*d) > 0 {
			dHosts = strings.Split(*d, ",")
			//fmt.Println(dHosts)
			connectDestinationHost()
		}
		//
		// Redis server/cluster commanding
		//
		// Cmd: keys
		if *keys {
			// Filter source keys
			var keys = filterSourceKeys(*f)
			n := len(keys)
			if n <= 0 {
				fmt.Println(noKeyD)
				break
			}
			// Retrieve key(s)
			for i := 0; i < len(keys); i++ {
				fmt.Println(keys[i])
			}
		}
		// Cmd: sync
		if *sync {
			if !dConnected {
				help()
				break
			}
			// Lister key(s)
			if len(*l) > 0 {
				if *f != "*" {
					log.Fatalln("Option [l]ister conflict with [f]ilter!")
				}
				var lister, err = os.Open(*l)
				if err != nil {
					log.Fatalln("Can not open [l]ister file!")
				}
				var listerScanner = bufio.NewScanner(lister)
				// Scan lister
				// Synchronize key(s)
				for listerScanner.Scan() {
					// Fetch the text
					// Synchronize the key from source to destination
					syncKey(listerScanner.Text())
				}
				// Filter key(s)
			} else {
				// Filter source keys
				var keys = filterSourceKeys(*f)
				n := len(keys)
				if n <= 0 {
					fmt.Println(noKeyD)
					break
				}
				// Synchronize key(s)
				for i := 0; i < n; i++ {
					syncKey(keys[i])
				}
			}
		}
		//
		// Finally
		if true {
			break
		}
	}
	// Count key(s)
	fmt.Println("Done " + strconv.FormatInt(keysCnt, 10) + " key(s).")
}

// Synchronize a key from the source host to the deestination host
func syncKey(key string) {
	keysCnt = keysCnt + 1
	var val string
	var ttl time.Duration
	if !sIsCluster {
		val = sHost.Dump(key).Val()
		ttl = sHost.PTTL(key).Val()
	} else {
		val = sBeam.Dump(key).Val()
		ttl = sBeam.PTTL(key).Val()
	}
	if ttl < -1 {
		// NoOp
		// !
		ttl = 0 // Rotten
	} else {
		if ttl == -1 {
			ttl = 0
		}
		// DoOp
		// ...
	}
	if !dIsCluster {
		dHost.Restore(key, ttl, val)
	} else {
		dBeam.Restore(key, ttl, val)
	}
}

// Filter keys of the source server/cluster
func filterSourceKeys(f string) []string {
	var keys *redis.StringSliceCmd
	if !sIsCluster {
		keys = sHost.Keys(f)
	} else {
		keys = sBeam.Keys(f)
	}
	return keys.Val()
}

// Ping to cluster
func clusterPingTest(redisClient *redis.ClusterClient) {
	var pingTest = redisClient.Ping()
	var pingMessage, pingError = pingTest.Result()
	if pingError != nil {
		log.Fatalf("%v m=%v e=%v", "Error when pinging the host!", pingMessage, pingError)
	}
}

// Ping to server
func serverPingTest(redisClient *redis.Client) {
	var pingTest = redisClient.Ping()
	var pingMessage, pingError = pingTest.Result()
	if pingError != nil {
		log.Fatalf("%v m=%v e=%v", "Error when pinging the host!", pingMessage, pingError)
	}
}

// Connects to the destination server/cluster
func connectDestinationHost() {
	if len(dHosts) == 1 {
		dHost = redis.NewClient(&redis.Options{
			Addr: dHosts[0],
		})
		dIsCluster = false
		serverPingTest(dHost)
	} else {
		dBeam = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: dHosts,
		})
		dIsCluster = true
		clusterPingTest(dBeam)
	}
	dConnected = true
}

// Connects to the source server/cluster
func connectSourceHost() {
	if len(sHosts) == 1 {
		sHost = redis.NewClient(&redis.Options{
			Addr: sHosts[0],
		})
		sIsCluster = false
		serverPingTest(sHost)
	} else {
		sBeam = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: sHosts,
		})
		sIsCluster = true
		clusterPingTest(sBeam)
	}
	sConnected = true
}

// Help
func help() {
	fmt.Println(`
RCS - Redis server/Cluster Synchronizer

Cmds:
  -keys=false ` + keysD + `
  -sync=false ` + syncD + `

Opts:
  -l=""  ` + lD + `
  -f="*" ` + fD + `
  -s=""  ` + sD + `
  -d=""  ` + dD + `
	`)
	os.Exit(0)
}

var (
	// Redis server/cluster client
	sHost *redis.Client
	dHost *redis.Client
	sBeam *redis.ClusterClient
	dBeam *redis.ClusterClient
	// Redis server/cluster connecting
	sHosts []string
	dHosts []string
	// Redis to connect is cluster or not
	sIsCluster = false
	dIsCluster = false
	// Redis server/cluster connected
	sConnected = false
	dConnected = false
	// Count of keys
	keysCnt int64
)

const (
	// Descriptions
	noKeyD = "No key in source host."
	keysD  = "Retrieve [keys] from source host."
	syncD  = "[Sync]hronize keys specified by -l/-f from source host to destination host."
	lD     = "The [l]ist of keys specified by a file lines"
	fD     = "The [f]ilter of keys specified by a pattern"
	sD     = "The [s]ource list of host:port server/cluster (Ex. 127.0.0.1:6379)"
	dD     = "The [d]estination list of host:port server/cluster (Ex. 127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002)"
)

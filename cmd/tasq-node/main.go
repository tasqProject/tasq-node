// Command tasq-node is the operator reference client for TasQ.
//
// It registers a machine, advertises an offer, and serves jobs from the
// coordinator. The scheduling and enclave backends are provided at build time;
// this reference wiring returns clear errors where the private parts belong.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tasqProject/tasq-node/internal/offer"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "https://api.tasqnetwork.io", "coordinator base URL")
		gpu      = flag.String("gpu", "", "GPU model to advertise, for example \"H100 80GB\"")
		region   = flag.String("region", "", "region to advertise")
		price    = flag.String("price", "", "price per hour in network credits")
		modes    = flag.String("modes", "R", "assurance modes to offer, any of A R P")
	)
	flag.Parse()

	if *gpu == "" || *region == "" || *price == "" {
		flag.Usage()
		os.Exit(2)
	}

	o := offer.Offer{
		MachineID:    deriveMachineID(),
		GPU:          *gpu,
		Modes:        parseModes(*modes),
		Region:       *region,
		PricePerHour: *price,
	}
	if err := o.Validate(); err != nil {
		log.Fatalf("invalid offer: %v", err)
	}

	fmt.Printf("tasq-node: advertising %s in %s at %s/hr, modes %v\n", o.GPU, o.Region, o.PricePerHour, o.Modes)
	fmt.Printf("coordinator: %s\n", *endpoint)

	// TODO: register with the coordinator, publish the offer, then serve matched
	// jobs with internal/job.Runner. The serve loop lives in the private build
	// alongside the enclave backend.
	log.Fatal("serve loop is not included in the reference client")
}

func parseModes(s string) []offer.Mode {
	out := make([]offer.Mode, 0, len(s))
	for _, c := range s {
		switch c {
		case 'A':
			out = append(out, offer.Attested)
		case 'R':
			out = append(out, offer.Redundant)
		case 'P':
			out = append(out, offer.Proven)
		}
	}
	return out
}

// deriveMachineID is a placeholder. The real id ties the operator key to the
// hardware, see tasq-protocol for how the ledger id is derived.
func deriveMachineID() string { return "machine-local" }

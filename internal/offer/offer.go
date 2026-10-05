// Package offer describes what a machine advertises to the market.
package offer

// Mode is an assurance mode a machine can serve.
type Mode string

const (
	Attested  Mode = "A"
	Redundant Mode = "R"
	Proven    Mode = "P"
)

// Offer is what the coordinator publishes on behalf of a machine. Price is in
// network credits per hour, as a decimal string so it is exact on the wire.
type Offer struct {
	MachineID    string `json:"machineId"`
	GPU          string `json:"gpu"`
	Modes        []Mode `json:"modes"`
	Region       string `json:"region"`
	PricePerHour string `json:"pricePerHour"`

	// X25519 and ML-KEM public keys a client seals input against. In mode A
	// these must come from the current attestation, not a stored value.
	KeyX25519 []byte `json:"keyX25519"`
	KeyMLKEM  []byte `json:"keyMlkem"`
}

// Supports reports whether the machine advertises a given mode.
func (o Offer) Supports(m Mode) bool {
	for _, x := range o.Modes {
		if x == m {
			return true
		}
	}
	return false
}

// Validate checks an offer is well formed before it is published.
func (o Offer) Validate() error {
	if o.MachineID == "" {
		return errField("machineId")
	}
	if o.GPU == "" {
		return errField("gpu")
	}
	if len(o.Modes) == 0 {
		return errField("modes")
	}
	if o.Supports(Attested) && (len(o.KeyX25519) == 0 || len(o.KeyMLKEM) == 0) {
		return errMsg("mode A offer must carry attested keys")
	}
	return nil
}

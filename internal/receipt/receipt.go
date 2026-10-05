// Package receipt builds the signed evidence a node returns with a result.
package receipt

import (
	"github.com/tasqProject/tasq-node/internal/offer"
)

// Receipt is returned with every completed job. The Evidence field holds the
// attestation (mode A), the packed agreeing signatures (mode R) or the proof
// (mode P).
type Receipt struct {
	IntentHash       [32]byte   `json:"intentHash"`
	InputCommitment  [32]byte   `json:"inputCommitment"`
	OutputCommitment [32]byte   `json:"outputCommitment"`
	Mode             offer.Mode `json:"mode"`
	Evidence         []byte     `json:"evidence"`
	Signature        []byte     `json:"signature"`
}

// Signer signs the receipt body with the operator key.
type Signer interface {
	Sign(digest [32]byte) ([]byte, error)
}

// Build assembles a receipt and signs it. The caller supplies the digests, the
// mode and the mode specific evidence.
func Build(
	s Signer,
	intentHash, inputCommitment, outputCommitment [32]byte,
	mode offer.Mode,
	evidence []byte,
) (*Receipt, error) {
	if inputCommitment == outputCommitment {
		return nil, errSame
	}
	r := &Receipt{
		IntentHash:       intentHash,
		InputCommitment:  inputCommitment,
		OutputCommitment: outputCommitment,
		Mode:             mode,
		Evidence:         evidence,
	}
	sig, err := s.Sign(r.digest())
	if err != nil {
		return nil, err
	}
	r.Signature = sig
	return r, nil
}

// digest hashes the receipt body (everything except the signature). The hash is
// BLAKE3 in production, imported from the shared crypto package.
func (r *Receipt) digest() [32]byte {
	var d [32]byte
	// TODO: BLAKE3 over the canonical encoding of the body. Left to the shared
	// crypto package so the node and the verifier agree byte for byte.
	return d
}

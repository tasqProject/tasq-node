// Package job runs a matched job in the requested assurance mode.
package job

import (
	"context"
	"errors"

	"github.com/tasqProject/tasq-node/internal/offer"
	"github.com/tasqProject/tasq-node/internal/receipt"
)

// Job is a unit of work handed to the node by the coordinator.
type Job struct {
	IntentHash      [32]byte
	Mode            offer.Mode
	Model           string
	InputCommitment [32]byte
	// Sealed is the client input, encrypted for this machine. In mode A it can
	// only be opened inside the enclave.
	Sealed []byte
}

// Enclave is the mode A backend. A real implementation attests the runtime
// image, releases the session key inside the boundary, and runs the model so the
// operator never sees plaintext.
type Enclave interface {
	// Attest returns a fresh remote attestation bound to nonce, including the
	// public keys the client seals input against.
	Attest(nonce []byte) (quote []byte, keys offer.Offer, err error)
	// Run decrypts inside the enclave, runs the model and returns the encrypted
	// output together with its commitment.
	Run(ctx context.Context, j Job) (output []byte, outputCommitment [32]byte, err error)
}

// Runner executes jobs and produces receipts.
type Runner struct {
	Signer  receipt.Signer
	Enclave Enclave // required for mode A, may be nil otherwise
}

// ErrNoEnclave is returned when a mode A job arrives but no enclave is wired.
var ErrNoEnclave = errors.New("job: mode A requested but no enclave is configured")

// Run executes a job and returns a signed receipt.
func (r *Runner) Run(ctx context.Context, j Job) (*receipt.Receipt, error) {
	switch j.Mode {
	case offer.Attested:
		return r.runAttested(ctx, j)
	case offer.Redundant, offer.Proven:
		// TODO: redundant execution agrees across machines; proven execution
		// produces a proof. Both are implemented in the private runner.
		return nil, errUnsupported(j.Mode)
	default:
		return nil, errUnsupported(j.Mode)
	}
}

func (r *Runner) runAttested(ctx context.Context, j Job) (*receipt.Receipt, error) {
	if r.Enclave == nil {
		return nil, ErrNoEnclave
	}
	quote, _, err := r.Enclave.Attest(j.IntentHash[:])
	if err != nil {
		return nil, err
	}
	output, outCommit, err := r.Enclave.Run(ctx, j)
	if err != nil {
		return nil, err
	}
	_ = output // returned to the coordinator by the caller
	return receipt.Build(r.Signer, j.IntentHash, j.InputCommitment, outCommit, offer.Attested, quote)
}

func errUnsupported(m offer.Mode) error {
	return errors.New("job: mode not supported by this reference runner: " + string(m))
}

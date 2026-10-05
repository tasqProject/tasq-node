# tasq-node

The operator side of TasQ. Run this on a machine with spare CPU and GPU to join
the network: register the machine, advertise what it can do, accept jobs from
the coordinator, run them, and return a signed receipt.

This is a reference client. It shows the shape of an operator: the offer format,
the job lifecycle, the attestation hook for mode A, and how a receipt is built.
The scheduling, pricing and the enclave specific glue that we run in production
are not in this repository. The functions where they belong are marked and
return a clear error.

## What a node does

1. **Register.** Prove control of the machine and derive the ledger id that ties
   the operator key to the hardware.
2. **Advertise.** Publish an offer: GPU model, the modes it supports, region and
   price. See `internal/offer`.
3. **Run.** Pull a matched job, run it in the right mode, and in mode A only
   after the enclave has produced an attestation. See `internal/job`.
4. **Receipt.** Return the output with a signed receipt: BLAKE3 digests, the
   mode, and the mode specific evidence. See `internal/receipt`.

## Build

```
go build ./...
./tasq-node --help
```

## Status

Pre-launch. The coordinator endpoint and the on-chain addresses are not final.
The enclave runner is an interface here, not a working mode A backend.

## License

Apache License 2.0. See `LICENSE`. Copyright The TasQ Project.

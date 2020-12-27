# GO Implementation of BW6
A new secure and optimized pairing-friendly elliptic curve that is suitable for one layer proof composition. The curve is defined over a 761-bit prime field. It is at least five times faster to verify a Groth proof, compared to Zexe.

## TODO
- [x] Sparse multiplications for Fp6
- [x] Pairing implementation
    - [x] Precomputation of line evaluation for addition and doubling steps
    - [x] Optimal Ate Miller loop
    - [x] Final exponentiation
        - [ ] Cyclotomic exponentiation
    - [x] Multi-pair evaluation
    - [x] Pairing tests
- [ ] Scalar field
- [ ] GLV multiplication

## References
- [Optimized and secure pairing-friendly elliptic curves suitable for one layer proof composition](https://eprint.iacr.org/2020/351.pdf)
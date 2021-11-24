Demo Atomic Cross-chain swap between Ethereum and Hyperledger Fabric

- Ethereum: 
simple HTLC smart contract for transferring ether from Alice to Bob
  - deploy contract, paramers, revipient, hashlock, timelock in Unix epoch
  - commit: password: finalizing the transfer
  - reverting: if timelock is active, the transfer can be reverted to Alice

- HLF:
Hyperledger Fabric, simple value transfer between Bob and Alice and a HTLC implementation of the value transfer
  - MintToken: minting token to an account
  - Balance: balance token of an account
  - BurnToken: burning token from an account
  - Transfer: normal transfer, transferring from one account to another one
  - TransferConditional: transferring conditional with timelock and hashlock
  - Commit: finalize transfer with the password
  - Revert: revert the password

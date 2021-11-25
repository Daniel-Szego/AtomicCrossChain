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

Atomic Cross-chain swap steps:

1. Alice creates a secret passwords and creates a sha256(password) hash of it

2. Alice sends sha256(password) to Bob, Alice abd Bob agress on a t timelock

3. Both Alice and Bob create HTLC with sha256(password) and t timelock. Alice on Ethereum, Bob on Hyperledger Fabric

4. Alice Commit the transaction on Hyperledger Fabric, so she gets the HLF token and reveals the "password"

5. As the "password" is committed, Bob see it and he can commit the transaction on Ethereum, so he can get the ether. 

+1. If the transaction is not committed by Alice in step 4, after the timeout explires, all parties can revoke the transactions



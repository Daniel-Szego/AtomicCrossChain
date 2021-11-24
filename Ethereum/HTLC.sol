// SPDX-License-Identifier: MIT

pragma solidity >=0.7.0;

// Alice initiates a hased timelock conditional payment to Bob on chain 1
contract AliceHTLC_Chain1 {

    // test setup
    // deployer: 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4
    // Alice: 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4
    // Bob: 0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2

    // from address
    address payable public fromAlice;
    // to address
    address payable public toBob;
    // timeout
    uint256 public timeOut;
    // hashlock
    bytes32 public hashLock;    

    event CommitEvent(string preimage);

    // contructor
    constructor(address payable _toBob, bytes32 _hashLock, uint256 _timeOut) {
        fromAlice =  payable(msg.sender);
        toBob = _toBob;
        hashLock = _hashLock;
        timeOut = _timeOut;            
    }    
        
    // allow payments
    fallback () payable external {}

    receive() external payable {}

    // getting contract balance
    function getBalance() public view returns (uint256){
        return address(this).balance;
    } 
    
    // executing the transaction -> Bob gets the payment
    // if valid secretHash presented
    // if timeout still not reached
    function commit(string memory _secretHash) public {
       require(hashLock == sha256(abi.encodePacked(_secretHash)), "password is wrong");
       require(block.timestamp <= timeOut, "timelock already activated");
       toBob.transfer(address(this).balance);  

       emit CommitEvent(_secretHash);      
    }

    // reverting the transaction -> Alice gets the payment back
    // only if timeout still already reached    
    function reverting() public {
      require(block.timestamp > timeOut, "timelock still not active");
     fromAlice.transfer(address(this).balance);
    }

}
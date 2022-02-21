pragma solidity ^0.8.11;

/**
 * @title EnablePayGas
 * @author Bui Duc Huy<duchuy.124dk@gmail.com>
 * @dev Implementation of the contract enable pay gas.
 */
contract EnablePayGas {
  event EnablePayGas(address indexed contractAddress, address indexed payer);

  mapping(address => bool) enableContracts;
  uint256 _minimumBalance;

  constructor() public {
  }

  function enable(address _contract) public {
    enableContracts[_contract] = true;

    emit EnablePayGas(_contract, msg.sender);
  }
}

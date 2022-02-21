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

  function enable(address _contract) payable public {
    require(msg.value < _minimumBalance, "Coin98 EnablePayGas: Exceed Value");
    enableContracts[_contract] = true;

    emit EnablePayGas(_contract, msg.sender);
  }
}

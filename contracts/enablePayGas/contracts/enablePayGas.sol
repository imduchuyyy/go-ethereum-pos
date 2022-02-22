pragma solidity ^0.8.11;

/**
 * @title EnablePayGas
 * @author Bui Duc Huy<duchuy.124dk@gmail.com>
 * @dev Implementation of the contract enable pay gas.
 */
contract EnablePayGas {
  event EnablePayGas(address indexed contractAddress, bytes method, address indexed payer);

  mapping(address => mapping(bytes => bool)) public enableContracts;

  function enable(address _contract, bytes memory _method) payable public {
    // require(msg.value < _minimumBalance, "Coin98 EnablePayGas: Exceed Value");
    enableContracts[_contract][_method] = true;

    emit EnablePayGas(_contract, _method, msg.sender);
  }
}

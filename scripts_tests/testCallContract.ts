import Web3 from "web3"
require("dotenv").config()

const web3: Web3 = new Web3(new Web3.providers.HttpProvider('http://localhost:8545'));

const privateKey: string = process.env.PRIVATE_KEY ? process.env.PRIVATE_KEY : ""
const address: string = process.env.ADDRESS ? process.env.ADDRESS : ""

const bytecode: string = "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220d4dbe87ca9034d34315e9ed83d4d46fe7767b201977764ef3ca5a71de14ba55464736f6c634300080b0033"
const abi: string = `
[
	{
		"inputs": [],
		"name": "retrieve",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "num",
				"type": "uint256"
			}
		],
		"name": "store",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
`

const sendTransaction = async (transactionParameter: any) => {
  const chainId: number = await web3.eth.net.getId();
  const nonce: number = await web3.eth.getTransactionCount(address);

  const gas = 8000000

  const signedTx = await web3.eth.accounts.signTransaction({
    ...transactionParameter,
    nonce,
    gas,
    gasPrice: 100,
    chainId,
  }, privateKey)

  console.log(signedTx)

  const receipt = await web3.eth.sendSignedTransaction(signedTx.rawTransaction ? signedTx.rawTransaction : "")
  console.log(receipt)

  return receipt
}

const deployContract = async () => {
  const chainId: number = await web3.eth.net.getId();
  const nonce: number = await web3.eth.getTransactionCount(address);

  const gas = 8000000

  const receipt = await sendTransaction({
    from: address,
    data: "0x" + bytecode,
    gas,
  })
  console.log(receipt)
}

const callContract = async () => {
    const contractAddress: string = "0xA4aC15D64E5d718E03614Fb0DC566d1616E5dc7a"

    const contract: any = new web3.eth.Contract(JSON.parse(abi), contractAddress)

    const data = contract.methods.store(1).encodeABI()

    const receipt = await sendTransaction({
        from: address,
        data,
        to: contractAddress
    })

    console.log(receipt)
}

// deployContract()
callContract()

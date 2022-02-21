import Web3 from "web3"
require("dotenv").config()

const web3: Web3 = new Web3(new Web3.providers.HttpProvider('http://localhost:8545'));

const privateKey: string = process.env.PRIVATE_KEY ? process.env.PRIVATE_KEY : "";
const address: string = process.env.ADDRESS ? process.env.ADDRESS : "";

const sendtx = async () => {
  const chainId: number = await web3.eth.net.getId();
  const nonce: number = await web3.eth.getTransactionCount(address);

  const gas = 21000

  const signedTx = await web3.eth.accounts.signTransaction({
    from: address,
    to: "0x379860934aa5D9A87c3164533848D065A3cf9B0f",
    data: "0x",
    gas,
    gasPrice: 100,
    nonce,
    chainId,
    value: 1
  }, privateKey)

  console.log(signedTx)

  const receipt = await web3.eth.sendSignedTransaction(signedTx.rawTransaction ? signedTx.rawTransaction : "")
  console.log(receipt)
}

const getInfo = async () => {
  const transaction = await web3.eth.getTransaction("0xaf8e77ec3a35f90083f5e002595c3288e598ca9b385ae2ef1050a4f51198550b")
  console.log(transaction)
}

getInfo()

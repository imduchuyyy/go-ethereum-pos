import Web3 from "web3"

(async () => {
    const web3: Web3 = new Web3(new Web3.providers.HttpProvider('http://localhost:8545'));

    let address = "000000000000000000000000025622aea16db104c1040e42251b59922d33e9cf"
    let method = "a6f9dae1"
    let pos = "0000000000000000000000000000000000000000000000000000000000000000"

    let key1 = web3.utils.keccak256(address + pos)
    if (key1 == null) {
        return
    }
    console.log("key 1", key1)
    const finalKey = web3.utils.keccak256(method + key1.slice(2, key1.length))

    console.log("finalKey", finalKey)
    console.log("exactKey", "0xf05110a882f897cdb8f48b569417b8e3a5263e6e10787653cb8c8aa00a6d4064")

    const myKeyData = await web3.eth.getStorageAt("0x0000000000000000000000000000000000001000", finalKey)
    const data = await web3.eth.getStorageAt("0x0000000000000000000000000000000000001000", "0xf05110a882f897cdb8f48b569417b8e3a5263e6e10787653cb8c8aa00a6d4064")

    console.log(data)
    console.log("my key data", myKeyData)
})()

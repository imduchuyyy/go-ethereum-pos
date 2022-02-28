import web3 from "web3"

(async () => {
    let address = "000000000000000000000000856C2afd19b368DD38449C751f831FedEa9542fc"
    let pos = "0000000000000000000000000000000000000000000000000000000000000000"
    console.log(address + pos)
    let key = web3.utils.sha3(address + pos)
    if (key == null) {
        return
    }
    console.log(key)
    console.log(key.slice(2, key.length))
    console.log(web3.utils.sha3("00" + key.slice(2, key.length)))
})()

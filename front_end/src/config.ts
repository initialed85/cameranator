export let apiHttpUrl = `https://cameranator.initialed85.cc/v1/graphql`
export let apiWsUrl = `wss://cameranator.initialed85.cc/api/v1/graphql`
export let fileHttpUrl = `https://cameranator.initialed85.cc/`

if (window?.location?.port === "3000") {
    apiHttpUrl = `http://localhost:8080/v1/graphql`
    apiWsUrl = `ws://localhost:8080/v1/graphql`
    fileHttpUrl = `https://cameranator.initialed85.cc/`
}

console.log(`apiHttpUrl = ${apiHttpUrl}`)
console.log(`apiWsUrl = ${apiWsUrl}`)
console.log(`fileHttpUrl = ${fileHttpUrl}`)

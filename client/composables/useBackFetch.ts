export const useBackFetch = async (uri: string, method: string, body?: object, headers?: object) => {
    return await $fetch.raw('http://localhost:8080/'+uri, { method, body, headers })
}
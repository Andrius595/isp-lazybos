import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const body = await readBody(event)

    let response
    try {
        const res = await useBackFetch('bet-user/register', 'POST', body)

        response = { status: true, data:  res._data}
    } catch(e) {
        response = { status: false, message: e.data.message}
    }

    return response
})

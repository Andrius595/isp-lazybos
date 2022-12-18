import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    let response
    try {
        const res = await useBackFetch('bet-user/logout', 'POST')

        response = { status: true, data:  res._data}
    } catch(e) {
        response = { status: false, message: e.data?.message ?? 'Something went wrong'}
    }

    deleteCookie(event, 'sessionup')
    deleteCookie(event, 'user')

    return response
})

import {getCookie} from "h3";
import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const sessionCookie = getCookie(event, 'sessionup')

    if (!sessionCookie) {
        return false
    }

    const headers = {
        'cookie': event.req.headers.cookie,
    }

    let response
    try {
        const res = await useBackFetch('bet-user/me', 'GET', undefined, headers)

        response = { status: true, data:  res._data}
    } catch(e) {
        response = { status: false, message: e.data?.message ?? 'Something went wrong'}
    }

    return response
})

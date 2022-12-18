import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const headers = {
        'cookie': event.req.headers.cookie,
    }

    let response
    try {
        const res = await useBackFetch('bet-user/bets', 'GET', undefined, headers )

        response = { status: true, data:  res._data}
    } catch(e) {
        console.log(e)
        response = { status: false, message: e.data?.message ?? 'Something went wrong'}
    }

    console.log('rrr', response)

    return response
})

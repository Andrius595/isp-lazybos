import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const body = await readBody(event)
    const headers = {
        'cookie': event.req.headers.cookie,
    }

    let response
    try {
        const res = await useBackFetch(`bet-user/autobet/${body.bet_uuid}`, 'DELETE', undefined, headers )

        response = { status: true, data:  res._data}
    } catch(e) {

        response = { status: false, message: e.data?.message ?? 'Something went wrong'}
    }

    return response
})

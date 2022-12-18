import {useBackFetch} from "~/composables/useBackFetch";
import {readBody} from "h3";

export default defineEventHandler(async (event) => {
    const body = await readBody(event)
    const headers = {
        'cookie': event.req.headers.cookie,
    }

    let response
    try {
        const res = await useBackFetch(`admin/report/admins`, 'POST', body, headers )

        response = { status: true, data:  res._data}
    } catch(e) {
        console.log(e)
        response = { status: false, message: e.data?.message ?? 'Something went wrong'}
    }

    return response
})

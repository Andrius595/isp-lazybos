import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const headers = {
        'cookie': event.req.headers.cookie,
    }

    let response
    try {
        const res = await useBackFetch('admin/identity-verifications', 'GET', undefined, headers)

        response = { status: true, data:  res._data}
    } catch(e) {
        response = { status: false, message: e.data.message}
    }

    return response
})

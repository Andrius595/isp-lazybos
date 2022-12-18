import {useBackFetch} from "~/composables/useBackFetch";

export default defineEventHandler(async (event) => {
    const body = await readBody(event)

    let response
    try {
        const res = await useBackFetch('admin/login', 'POST', body)

        const cookieHeader = res.headers.get('set-cookie');
        const nameValue = cookieHeader.split(';')[0]
        const value = nameValue.split('=')[1]

        setCookie(event, 'sessionup', value)

        response = { status: true, data:  res._data}
    } catch(e) {
        response = { status: false, message: e.data.message}
    }

    return response
})

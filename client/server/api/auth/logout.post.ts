export default defineEventHandler(async (event) => {
    try {
        await useBackFetch('bet-user/logout', 'POST')
    } catch (e) {}

    setCookie(event, 'sessionup', undefined)
})

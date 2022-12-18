import Routes from "~/types/routes";

export default defineNuxtRouteMiddleware(async () => {
    const user = useCookie('user')

    if (!user.value) {
        return navigateTo({ name: Routes.Auth.Login })
    }
})
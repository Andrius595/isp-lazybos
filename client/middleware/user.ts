import getRouteUrl from "~/utils/getRouteUrl";
import Routes from "~/types/routes";

export default defineNuxtRouteMiddleware(async () => {
    const response = await $fetch('/api/auth/me', { method: 'POST' })

    if (!response.status) {
        return navigateTo(getRouteUrl(Routes.Auth.Login))
    }
})
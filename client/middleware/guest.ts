import Routes from "~/types/routes";

export default defineNuxtRouteMiddleware(async () => {
    const response = await $fetch('/api/auth/me', { method: 'POST' })

    if (response.status) {
        return navigateTo({ name: Routes.Main })
    }
})
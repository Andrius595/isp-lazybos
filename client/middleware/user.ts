import Routes from "~/types/routes";
import checkUserRole from "~/utils/checkUserRole";

export default defineNuxtRouteMiddleware(async () => {
    const user = useCookie('user')

    const isUser = checkUserRole(user.value, 'user')

    if (!isUser) {
        return navigateTo({ name: Routes.Main })
    }
})
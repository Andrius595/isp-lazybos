import Routes from "~/types/routes";
import checkUserRole from "~/utils/checkUserRole";

export default defineNuxtRouteMiddleware(async () => {
    const user = useCookie('user')

    const isSalesAdmin = checkUserRole(user.value, 'users')

    if (!isSalesAdmin) {
        return navigateTo({ name: Routes.Main })
    }
})
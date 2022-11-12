import { LocationQueryRaw, RouteParams, useRouter } from "vue-router";

export default (name: string, params: RouteParams = {}, query: LocationQueryRaw = {}) => {
  const router = useRouter()

  return router.resolve({ name, params, query })
}
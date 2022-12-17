import { useRoute } from "vue-router";

export default (routeName: string) => {
  const route = useRoute()
  return route?.name === routeName
}
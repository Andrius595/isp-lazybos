<template>
  <div>
    <AuthenticatedLayout>
      <h1>Paskyros langas</h1>
      <ul v-if="user" class="list-group mb-3">
        <li class="list-group-item">Name: {{ `${user.first_name} ${user.last_name}` }}</li>
        <li class="list-group-item">Email: {{ user.email }}</li>
        <li class="list-group-item">Balance: {{ user.balance }} Eur</li>
        <li class="list-group-item">Identity status: {{ user.identitity_verified ? 'Verified' : 'Not verified' }}</li>
      </ul>
      <div v-if="user && !user.identitity_verified" class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Identity.Request)">Request Identity Verification</NuxtLink>
        </ul>
      </div>

    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";

const user = ref(null)

await fetchUser()

async function fetchUser() {
  const response = await $fetch('/api/auth/me', { method: 'POST' })

  if (!response.status) {
    return navigateTo({ name: Routes.Auth.Login })
  }

  user.value = response.data
}
</script>

<style scoped lang="scss">

</style>
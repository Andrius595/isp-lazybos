<template>
  <div>
    <AuthenticatedLayout>
      <h1>Paskyros langas</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Identity.Request)">Prašymo langas</NuxtLink>
        </ul>
      </div>

    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";

await $fetch('/api/auth/me', { method: 'POST' })
await $fetch('/api/events/user-bets', { method: 'POST' })
</script>

<style scoped lang="scss">

</style>
<template>
  <div>
    <AuthenticatedLayout>
      <h1>Rungtynių sąrašas</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">Rungtynių užbaigimo patvirtinimas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Rungtynių pašalinimo patvirtinimas</NuxtLink>
        </ul>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";
</script>
<style scoped lang="scss">

</style>
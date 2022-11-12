<template>
  <div>
    <AuthenticatedLayout>
      <h1>Prašymų sąrašas</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Identity.Confirm, {id: 1})">Prašymo patvirtinimo langas</NuxtLink>
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
<template>
  <div>
    <AuthenticatedLayout>
      <h1>Ataskaitų pasirinkimo langas</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Profit)">Pelningumo ataskaitos langas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Matches)">Varžybų ataskaitos langas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Taxes)">Mokesčių ataskaitos langas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Users)">Vartotojų ataskaitos langas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.AdminsActions)">Administratorių veiksmų ataskaitos langas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.ScheduledReport)">Suplanuotos ataskaitos langas</NuxtLink>
        </ul>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";
import {definePageMeta} from "#imports";

definePageMeta({
  middleware: 'sales-admin',
})
</script>

<style scoped lang="scss">

</style>
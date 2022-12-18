<template>
  <div>
    <AuthenticatedLayout>
      <h1>Select report</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Profit)">Profit Report</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Matches)">Matches Report</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Taxes)">Taxes Report</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.Users)">Users Report</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.AdminsActions)">Administrators actions Report</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.Reports.ScheduledReport)">Scheduled Report</NuxtLink>
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
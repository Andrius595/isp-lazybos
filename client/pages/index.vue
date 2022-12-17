<template>
  <div>
    <AuthenticatedLayout>
      <h1>Pagrindinis langas</h1>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";
import { useRouter } from "vue-router";
import getRouteUrl from "~/utils/getRouteUrl";
import Routes from "~/types/routes";

const router = useRouter()
router.replace(getRouteUrl(Routes.Main))
</script>
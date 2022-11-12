<template>
  <div>
    <AuthenticatedLayout>
      <h1>Rungtynių duomenų forma</h1>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import { useRoute, useRouter } from "vue-router";
const router = useRouter()
const route = useRoute()
console.log(router.resolve(route))
</script>
<style scoped lang="scss">

</style>
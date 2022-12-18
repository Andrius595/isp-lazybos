<template>
  <div>
    <AuthenticatedLayout>
      <h1>Statymų sąrašas</h1>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

await fetchBets()

async function fetchBets() {
  await $fetch('/api/events/user-bets', { method: 'POST' })
}
</script>
<style scoped lang="scss">

</style>
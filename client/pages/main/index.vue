<template>
  <div>
    <AuthenticatedLayout>
      <h1>Pagrindinis langas</h1>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import {definePageMeta} from "#imports";

definePageMeta({
  middleware: 'logged-in',
})
</script>

<style scoped lang="scss">

</style>
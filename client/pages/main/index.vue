<template>
  <div>
    <AuthenticatedLayout>
      <h1 v-if="user">Hello, {{user.first_name}}</h1>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import {definePageMeta} from "#imports";

definePageMeta({
  middleware: 'logged-in',
})

const user = useCookie('user')
</script>

<style scoped lang="scss">

</style>
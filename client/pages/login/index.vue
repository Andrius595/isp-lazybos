<template>
  <div>
    <GuestLayout>
      <h1 class="mb-4">Login</h1>
      <div class="d-flex justify-content-center">
        <form class="col-sm-12 col-lg-6">
          <span v-if="errorMessage.length" class="text-danger">{{ errorMessage }}</span>
        <!-- Email input -->
        <div class="form-outline mb-4">
          <input v-model="formData.email" type="email" class="form-control" />
          <label class="form-label">Email</label>
        </div>

        <!-- Password input -->
        <div class="form-outline mb-4">
          <input v-model="formData.password" type="password" class="form-control" />
          <label class="form-label">Password</label>
        </div>


        <!-- Submit button -->
        <button type="button" class="btn btn-primary btn-block mb-4" @click="handleLogin">Login</button>

        <!-- Register buttons -->
        <div class="text-center">
          <p>Not a member? <NuxtLink :href="getRouteUrl(Routes.Auth.Register)">Register</NuxtLink></p>
        </div>
      </form>
      </div>
    </GuestLayout>
  </div>
</template>

<script setup lang="ts">
import GuestLayout from "~/layouts/GuestLayout.vue";
import getRouteUrl from '~/utils/getRouteUrl'
import Routes from "~/types/routes";
import {definePageMeta} from "#imports";

const formData = ref({
  email: '',
  password: '',
})

const errorMessage = ref('')

async function handleLogin() {
  errorMessage.value = ''

  const res = await $fetch('api/auth/login', {
    method: 'POST',
    body: formData.value,
  })

  if (!res.status) {
    errorMessage.value = res.message

    return
  }

  return navigateTo(getRouteUrl(Routes.Main))
}

definePageMeta({
  middleware: 'guest',
})
</script>

<style scoped lang="scss">

</style>
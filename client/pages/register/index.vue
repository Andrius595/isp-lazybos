<template>
  <div>
    <GuestLayout>
      <h1>Registration</h1>
      <div class="d-flex justify-content-center">
        <form class="col-sm-12 col-lg-6">
          <span v-if="errorMessage.length" class="text-danger">{{ errorMessage }}</span>
          <!-- Email input -->
          <div class="row">
            <div class="col-sm-12 col-md-6">
              <div class="form-outline mb-4">
                <label class="form-label">First Name</label>
                <input v-model="formData.first_name" type="text" class="form-control"/>
              </div>
            </div>
            <div class="col-sm-12 col-md-6">
              <div class="form-outline mb-4">
                <label class="form-label">Last Name</label>
                <input v-model="formData.last_name" type="text" class="form-control"/>
              </div>
            </div>
          </div>
          <div class="form-outline mb-4">
            <label class="form-label">Email</label>
            <input v-model="formData.email" type="email" class="form-control"/>
          </div>

          <!-- Password input -->
          <div class="form-outline mb-4">
            <label class="form-label">Password</label>
            <input v-model="formData.password" type="password" class="form-control"/>
          </div>

          <!-- Submit button -->
          <button type="button" class="btn btn-primary btn-block mb-4" @click="handleRegistration">Register</button>

          <!-- Register buttons -->
          <div class="text-center">
            <p>Already have an account?
              <NuxtLink :href="getRouteUrl(Routes.Auth.Login)">Log in</NuxtLink>
            </p>
          </div>
        </form>
      </div>

    </GuestLayout>
  </div>
</template>

<script setup lang="ts">
import GuestLayout from "~/layouts/GuestLayout.vue";
import getRouteUrl from "~/utils/getRouteUrl";
import Routes from "~/types/routes";
import {definePageMeta} from "#imports";

const formData = ref({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
})

const errorMessage = ref('')

async function handleRegistration() {
  errorMessage.value = ''
  const response = await $fetch('/api/auth/register', {
    method: 'POST',
    body: formData.value,
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  return navigateTo(getRouteUrl(Routes.Auth.Login))
}


definePageMeta({
  middleware: 'guest',
})
</script>

<style scoped lang="scss">

</style>
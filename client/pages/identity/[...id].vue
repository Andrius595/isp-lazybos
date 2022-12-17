<template>
  <div>
    <AuthenticatedLayout>
      <h1>Identity Request Preview</h1>
      <h5>User information:</h5>
      <ul class="list-group my-2">
        <li class="list-group-item">{{ `Name: ${request.user.first_name} ${request.user.first_name}` }}</li>
        <li class="list-group-item">{{ `Email: ${request.user.first_name} ${request.user.first_name}` }}</li>
        <li class="list-group-item">{{ `Created at: ${formatDate(request.created_at)}` }}</li>
      </ul>
      <div class="d-flex gap-2 border border-dark rounded p-4">
        <div>
          <h5>ID Photo:</h5>
          <img :src="request.id_photo_base_64" alt="ID Photo" class="w-100">
        </div>
        <div>
          <h5>Portrait Photo:</h5>
          <img :src="request.portrait_photo_base_64" alt="ID Photo" class="w-100">
        </div>
      </div>
      <div class="d-flex gap-2 mt-2">
        <button class="btn btn-success" @click="handleAction(true)">Confirm</button>
        <button class="btn btn-danger" @click="handleAction(false)">Deny</button>

      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import formatDate from "~/utils/formatDate";

const route = useRoute()

const errorMessage = ref('')
const uuid = computed(() => Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)

const request = ref(await fetchData(uuid.value))

async function fetchData(uuid: string) {
  const response = await $fetch('/api/identity-verifications/identity-requests', {method: 'POST'})

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  const request = response.data.find((request) => request.uuid === uuid)

  if (!request) {
    return navigateTo({ name: Routes.Identity.List })
  }

  return request
}

async function handleAction(action: boolean) {
  const response = await $fetch('/api/identity-verifications/action', {
    method: 'POST',
    body: {
      verification_uuid: uuid.value,
      accept: action,
    }
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  return navigateTo({ name: Routes.Identity.List })
}

watch(
    () => uuid.value,
    async (newValue) => {
      await fetchData(newValue)
    }
)
</script>

<style scoped lang="scss">

</style>
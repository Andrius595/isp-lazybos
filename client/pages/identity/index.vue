<template>
  <div>
    <AuthenticatedLayout>
      <h1>Identity verification requests</h1>
      <table v-if="requests.length" class="table">
        <thead>
        <tr>
          <td>No.</td>
          <td>User</td>
          <td>Status</td>
          <td>Created at</td>
          <td></td>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(request, index) in requests" :key="index">
          <td>{{ index }}</td>
          <td>{{ `${request.user.first_name} ${request.user.last_name}` }}</td>
          <td>{{ request.status }}</td>
          <td>{{ formatDate(request.created_at) }}</td>
          <td>
            <button class="btn btn-outline-primary" @click="handleRedirect(request.uuid)">View</button>
          </td>
        </tr>
        </tbody>
      </table>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import formatDate from "~/utils/formatDate";

const router = useRouter()

const errorMessage = ref('')
const requests = ref([])

await fetchData()

async function fetchData() {
  const response = await $fetch('/api/identity-verifications/identity-requests', { method: 'POST' })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  requests.value = response.data.filter((request) => request.status === 'pending')
}

function handleRedirect(uuid: string) {
  return navigateTo({ name: Routes.Identity.Confirm, params: { id: uuid } })
}
</script>

<style scoped lang="scss">

</style>
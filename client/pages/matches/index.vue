<template>
  <div>
    <AuthenticatedLayout>
     <h1>Varžybų sąrašas</h1>
      <table v-if="matches && matches.length" class="table">
        <thead>
        <tr>
          <td>No.</td>
          <td>Name</td>
          <td>Status</td>
          <td>Begins at</td>
          <td></td>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(match, index) in matches" :key="index">
          <td>{{ index }}</td>
          <td>{{ match.name }}</td>
          <td>{{ match.finished ? 'Finished' : 'In progress' }}</td>
          <td>{{ formatDate(match.begins_at) }}</td>
          <td>
            <div v-if="!match.finished" class="btn btn-primary" @click="handleBetButton(match.uuid)">Bet</div>
          </td>
        </tr>
        </tbody>
      </table>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import formatDate from "~/utils/formatDate";
import Routes from "~/types/routes";

const errorMessage = ref('')
const matches = ref([])

await fetchData()

async function fetchData() {
  const response = await $fetch('/api/events/all', { method: 'POST' })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  matches.value = response.data ?? []
}

function handleBetButton(uuid) {
 return navigateTo({ name: Routes.Matches.Information, params: { id: uuid}})
}
</script>


<style scoped lang="scss">

</style>
<template>
  <div>
    <AuthenticatedLayout>
      <h1>Users Report</h1>
      <div class="container">
        <div class="mt-2 mx-5">
          <h5>Select client:</h5>
        </div>
        <div class="mx-5" style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <select class="form-select" aria-label="adminSelect" @change="handleUserSelect">
            <option selected>Select client...</option>
            <option v-for="(user, index) in users" :key="index" :value="user.uuid">{{ user.email }}</option>
          </select>
        </div>
        <button type="button" style="width: 250px; margin: 20px 45px; border: 2px solid black;
              border-radius: 8px; background-color: lightgreen;" @click="handleShowBets">Show client's behaviour</button>
        <div class="row">
          <!-- Atsiranda paspaudus mygtuka -->
          <div v-if="showTable" class="mt-5 mx-auto"
               style="border: 2px solid black; border-radius: 8px; padding: 15px 10px 0 10px">
            <h4>Client bets</h4>
            <table class="table">
              <thead style="background-color: #4CBB17;">
              <tr>
                <th scope="col">Sum</th>
                <th scope="col">Date</th>
                <th scope="col">Coefficient</th>
                <th scope="col">Bet status</th>
                <th scope="col">Selection winner</th>
              </tr>
              </thead>
              <tbody class="table-secondary">
              <tr v-for="(bet, index) in bets" :key="index">
                <td>{{ bet.stake }} Eur</td>
                <td>{{ formatDate(bet.timestamp) }}</td>
                <td>{{ bet.odds }}</td>
                <td>{{ bet.status }}</td>
                <td>{{ bet.selection.winner }}</td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import formatDate from "~/utils/formatDate";

const users = ref([])
const bets = ref([])

const showTable = ref(false)

const selectedUser = ref(null)

await fetchUsers()
async function fetchUsers() {
  const response = await $fetch('/api/users/all', { method: 'POST' })

  if (!response.status) {
    // TODO show error

    return
  }

  users.value = response.data
}

function handleUserSelect(event) {
  showTable.value = false
  bets.value = []
  const user = users.value.find((u) => u.uuid === event.target.value)

  if (!user) {
    selectedUser.value = null
    // TODO show error

    return
  }

  selectedUser.value = user
}

async function handleShowBets() {
  const response = await $fetch('/api/reports/users', {
    method: 'POST',
    body: {
      user_uuid: selectedUser.value?.uuid
    }
  })

  if (!response.status) {
    // TODO show error
    showTable.value = false

    return
  }

  showTable.value = true
  bets.value = response.data
}
</script>

<style scoped lang="scss">

</style>
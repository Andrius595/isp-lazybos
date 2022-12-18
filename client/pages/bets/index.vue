<template>
  <div>
    <AuthenticatedLayout>
      <h1>Bets list</h1>

      <table class="table table-striped">
        <thead>
        <tr>
          <th scope="col">Sum</th>
          <th scope="col">Date</th>
          <th scope="col">Coefficient</th>
          <th scope="col">Bet status</th>
          <th scope="col">Selection winner</th>
        </tr>
        </thead>
        <tbody>
          <tr v-for="(bet, index) in bets" :key="index">
            <td>{{ bet.stake }} Eur</td>
            <td>{{ formatDate(bet.timestamp) }}</td>
            <td>{{ bet.odds }}</td>
            <td>{{ bet.status == "tbd" ? "To be determined" : bet.status }}</td>
            <td>{{ bet.winner == "away" ? bet.event.away_team.name : bet.event.home_team.name }}</td>
          </tr>
        </tbody>
      </table>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import formatDate from "~/utils/formatDate";


const bets = ref([])
await fetchBets()

async function fetchBets() {
  const response = await $fetch('/api/events/user-bets', { method: 'POST' })

  if (!response.status) {
    // TODO show error

    return
  }

  bets.value = response.data
}
</script>
<style scoped lang="scss">

</style>

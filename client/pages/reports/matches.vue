<template>
  <div>
    <AuthenticatedLayout>
      <h1>Matches report</h1>
      <div class="border border-2 border-dark rounded p-4">
        <div>
          <div>
            <h5>Select time period:</h5>
          </div>
          <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
            <label for="timePeriod1" class="form-label">From:</label>
            <input v-model="dateFrom" type="date" class="form-control" id="timePeriod1">
            <label for="timePeriod2" class="form-label">To:</label>
            <input v-model="dateTo" type="date" class="form-control" id="timePeriod2">
            <button type="button" style="width: 150px; margin: 0 300px; border: 2px solid black;
            border-radius: 8px; background-color: lightgreen;" @click="handleShowBets">Show bets</button>
          </div>
        </div>

        <!-- Atsiranda paspaudus mygtuka -->
        <div v-if="bets" class="mt-5" style="width: auto; border: 2px solid black; border-radius: 8px;
        padding: 15px 10px 0 10px">
          <h4>Bets list</h4>
          <table class="table">
            <thead style="background-color: gray;">
            <tr>
              <th scope="col">Sum</th>
              <th scope="col">Date</th>
              <th scope="col">Coefficient</th>
              <th scope="col">Bet status</th>
              <th scope="col">Selection Winner</th>
            </tr>
            </thead>
            <tbody class="table-secondary">
            <tr v-for="(bet, index) in bets" :key="index">
              <td>{{ bet.stake }} Eur</td>
              <td>{{ formatDate(bet.timestamp) }}</td>
              <td>{{ bet.odds }}</td>
              <td>{{ bet.status == "tbd" ? "To be determined" : bet.status }}</td>
              <td>{{ bet.winner == "away" ? bet.event.away_team.name : bet.event.home_team.name }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import formatDate from "~/utils/formatDate";

const dateFrom = ref('')
const dateTo = ref('')

const bets = ref([])

async function handleShowBets() {
  const response = await $fetch('/api/reports/matches', {
    method: 'POST',
    body: {
      from: new Date(dateFrom.value).toISOString(),
      to: new Date(dateTo.value).toISOString(),
    }
  })

  if (!response.status) {
    // TODO show error
    return
  }

  bets.value = response.data
}
</script>

<style scoped lang="scss">

</style>

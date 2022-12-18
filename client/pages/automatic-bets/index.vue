<template>
  <div>
    <AuthenticatedLayout>
      <h1>Automatinių statymų sąrašas</h1>
      <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.AutomaticBets.Create)">Create Automatic Bet</NuxtLink>
        </ul>
      </div>
      <table class="table table-striped">
        <thead>
        <tr>
          <td>Is High Risk</td>
          <td>Balance Fraction</td>
          <td>Action</td>
        </tr>
        </thead>
        <tbody>
          <tr v-for="(bet, index) in bets" :key="index">
            <td>{{ bet.high_risk.toString() }}</td>
            <td>{{ bet.balance_fraction.toString() }}</td>
            <td>
              <div class="btn btn-danger" @click="handleDeleteBet(bet)">Delete</div>
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
import getRouteUrl from "~/utils/getRouteUrl";

const bets = ref([])

await fetchBets()

async function fetchBets() {
  const response = await $fetch('/api/autobets/all', {
    method: 'POST',
  })

  if (!response.status) {
    // TODO show error

    return
  }

  bets.value = response.data
}

async function handleDeleteBet(bet) {
  const response = await $fetch('/api/autobets/delete', {
    method: 'POST',
    body: {
      bet_uuid: bet.uuid,
    }
  })

  if (!response.status) {
    // TODO show error

    return
  }

  await fetchBets()
}
</script>
<style scoped lang="scss">

</style>
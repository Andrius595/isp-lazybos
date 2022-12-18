<template>
  <div>
    <AuthenticatedLayout>
      <div class="m-5">
        <h1>Profit report</h1>
      </div>
      <div class="mx-5">
        <h5>Select time period:</h5>
      </div>
      <div class="mx-5" style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
        <label for="timePeriod1" class="form-label">From:</label>
        <input v-model="dateFrom" type="date" class="form-control" id="timePeriod1">
        <label for="timePeriod2" class="form-label">To:</label>
        <input v-model="dateTo" type="date" class="form-control" id="timePeriod2">
        <button type="button" style="width: 150px; margin: 0 300px; border: 2px solid black;
            border-radius: 8px; background-color: lightgreen;" @click="handleShowProfit">Show profit/loss</button>
      </div>
      <!-- Atsiranda paspaudus mygtuka -->
      <div v-if="isTableShowing" class="mt-4 mx-5">
        <h5>Profit assessment:</h5>
        <h5>{{ results.profit }} EUR</h5><br>
        <h5>Loss assessment:</h5>
        <h5>{{ results.loss }} EUR</h5><br>
        <h5>Final profit/loss assessment:</h5>
        <h5>{{ results.final }} EUR</h5><br>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const dateFrom = ref('')
const dateTo = ref('')
const isTableShowing = ref(false)

const results = ref({
  profit: '0',
  loss: '0',
  final: '0',
})

async function handleShowProfit() {
  const response = await $fetch('/api/reports/profit', {
    method: 'POST',
    body: {
      from: new Date(dateFrom.value).toISOString(),
      to: new Date(dateTo.value).toISOString(),
    }
  })

  console.log(response)

  if (!response.status) {
    isTableShowing.value = false
    // TODO show error
    return
  }

  isTableShowing.value = true
  results.value = response.data
}
</script>

<style scoped lang="scss">

</style>
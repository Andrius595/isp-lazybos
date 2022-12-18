<template>
  <div>
    <AuthenticatedLayout>
      <h1>Taxes report</h1>
      <div class="border border-2 border-dark rounded p-5">
        <div>
          <h5>Enter tax rate (%):</h5>
        </div>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <input v-model.number="taxRate" type="number" min="0" max="100" class="form-control" id="taxRate">
        </div>
        <div class="mt-2">
          <h5>Select year:</h5>
        </div>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <select @change="year = Number($event.target.value)" class="form-select" aria-label="yearSelect">
            <option selected>Year...</option>
            <option value="2022">2022</option>
            <option value="2021">2021</option>
            <option value="2020">2020</option>
            <option value="2019">2019</option>
            <option value="2018">2018</option>
            <option value="2017">2017</option>
            <option value="2016">2016</option>
          </select>
        </div>
        <button type="button" style="width: 150px; margin: 0 350px; border: 2px solid black;
        border-radius: 8px; background-color: lightgreen;" @click="handleShowTax">Show tax</button>
        <!-- Atsiranda paspaudus mygtuka -->
        <div v-if="showTax" class="mt-4">
          <h5>Year's net profit:</h5>
          <h5>{{ netProfit }} EUR</h5><br>
          <h5>Year's tax:</h5>
          <h5>{{ yearTax }} EUR</h5><br>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const showTax = ref(false)
const year = ref(2022)
const taxRate = ref(0)

const results = ref({
  profit: '0',
  loss: '0',
  final: '0',
})

const netProfit = computed(() => {
  const final = Number(results.value.final)

  if (final > 0) {
    return final - final * (taxRate.value/100)
  }

  return final
})

const yearTax = computed(() => {
  const final = Number(results.value.final)

  if (final > 0) {
    return final * (taxRate.value/100)
  }

  return 0
})


async function handleShowTax() {
  const from = Number(year.value)

  const response = await $fetch('/api/reports/profit', {
    method: 'POST',
    body: {
      from: new Date(from.toString()).toISOString(),
      to: new Date((from+1).toString()).toISOString(),
    }
  })

  if (!response.status) {
    showTax.value = false
    // TODO show error
    return
  }

  showTax.value = true
  results.value = response.data
}
</script>

<style scoped lang="scss">

</style>
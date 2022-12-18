<template>
  <div>
    <AuthenticatedLayout>
      <h1>Scheduled report</h1>
      <div class="border border-2 border-dark rounded p-4">

        <h5>Enter email address:</h5>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <input v-model="email" type="email" class="form-control" id="email" placeholder="name@example.com">
        </div>
        <div class="mt-2">
          <h5>Select type of report:</h5>
        </div>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <select class="form-select" aria-label="yearSelect" @change="handleReportSelect">
            <option value="" selected>Select report...</option>
            <option value="profit">Profit report</option>
            <option value="deposit">Deposit report</option>
          </select>
        </div>
        <button type="button" style="width: 150px; margin: 10px 350px; border: 2px solid black;
        border-radius: 8px; background-color: lightgreen;" @click="handleSendReport">Send report</button>
        <div v-if="showSuccessMessage" class="mt-4">
          <h5 class="text-success">Report scheduled successfully!</h5>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const showSuccessMessage = ref(false)

const email = ref('')
const selectedReport = ref(null)

async function handleSendReport() {
  const response = await $fetch('/api/reports/auto', {
    method: 'POST',
    body: {
      send_to: email.value,
      type: selectedReport.value
    }
  })

  if (!response.status) {
    // TODO show error

    return
  }

  showSuccessMessage.value = true
}

function handleReportSelect(event) {
  selectedReport.value = event.target.value
}
</script>

<style scoped lang="scss">

</style>
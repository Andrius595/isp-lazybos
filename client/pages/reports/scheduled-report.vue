<template>
  <div>
    <AuthenticatedLayout>
      <h1>Scheduled report</h1>
      <div class="border border-2 border-dark rounded p-4">

        <h5>Enter email address:</h5>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <input type="email" class="form-control" id="email" placeholder="name@example.com">
        </div>
        <div class="mt-2">
          <h5>Select type of report:</h5>
        </div>
        <div style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <select class="form-select" aria-label="yearSelect">
            <option selected>Report...</option>
            <option value="profit">Profit report</option>
            <option value="taxes">Taxes report</option>
            <option value="adminActions">Administrators' actions report</option>
            <option value="clients">Clients report</option>
            <option value="matches">Matches report</option>
          </select>
        </div>
        <div class="mt-2">
          <h5>Select scheduled period:</h5>
        </div>
        <div  style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
          <select class="form-select" aria-label="yearSelect">
            <option selected>Period...</option>
            <option value="1">1 day</option>
            <option value="7">7 days</option>
            <option value="14">14 days</option>
          </select>
        </div>
        <button type="button" style="width: 150px; margin: 10px 350px; border: 2px solid black;
        border-radius: 8px; background-color: lightgreen;" @click="handleSendReport">Send report</button>
        <div v-if="showSuccessMessage" class="mt-4">
          <h5 class="text-success">Report has been successfully sent</h5>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const showSuccessMessage = ref(false)

function handleSendReport() {
  showSuccessMessage.value = !showSuccessMessage.value
}
</script>

<style scoped lang="scss">

</style>
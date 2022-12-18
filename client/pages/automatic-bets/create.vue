<template>
  <div>
    <AuthenticatedLayout>
      <h1>Create Automatic Bet</h1>
      <div>
        <span>Balance: {{ user.balance }} Eur</span>
      </div>
      <div class="mb-2">
        <label class="form-label">Balance fraction</label>
        <input v-model.number="balanceFraction" class="form-control" type="number" min="0" max="1" step="0.01" />
      </div>
      <div class="form-check">
        <input v-model="isHighRisk" class="form-check-input" type="checkbox">
        <label class="form-check-labell">Is High Risk</label>
      </div>
      <div class="btn btn-primary" @click="handleCreateBet">Create</div>

    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";

const user = ref(null)

const isHighRisk = ref(false)
const balanceFraction = ref(0)

await fetchUser()

async function fetchUser() {
  const response = await $fetch('/api/auth/me', { method: 'POST' })

  if (!response.status) {
    return navigateTo({ name: Routes.Auth.Login })
  }

  user.value = response.data
}

async function handleCreateBet() {
  const response = await $fetch('/api/autobets/create', {
    method: 'POST',
    body: {
      high_risk: isHighRisk.value,
      balance_fraction: balanceFraction.value,
    }
  })

  if (!response.status) {
    // TODO show error

    return
  }

  return navigateTo({ name: Routes.AutomaticBets.List })
}
</script>
<style scoped lang="scss">

</style>
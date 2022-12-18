<template>
  <div>
    <AuthenticatedLayout>
      <h1>Withdraw money</h1>
      <div v-if="successMessage.length" class="text-success">{{ successMessage }}</div>
      <div v-if="errorMessage.length" class="text-danger">{{ errorMessage }}</div>
      <div v-if="selectedUser">
        Selected user balance: {{ selectedUser.balance }}
      </div>
      <div class="d-flex gap-2">
        <div class="w-100">
          <div class="mb-3">
            <label class="form-label">User</label>
            <select class="form-select" aria-label="Select user" @change="handleUserSelect">
              <option selected value="">Select user</option>
              <option v-for="(user, index) in users" :key="index" :value="user.uuid">
                {{ `${user.first_name} ${user.last_name}` }}
              </option>
            </select>
          </div>
        </div>
        <div class="w-100">
          <div class="mb-3">
            <label class="form-label">Amount</label>
            <input v-model.number="amount" type="number" step="0.01" class="form-control" placeholder="0.0">
          </div>
        </div>
      </div>
      <button class="btn btn-success" @click="handleWithdraw">Withdraw</button>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const errorMessage = ref('')
const successMessage = ref('')
const users = ref([])
const amount = ref(0)
const selectedUser = ref(null)

await fetchUsers()


async function fetchUsers() {
  const response = await $fetch('/api/users/all', {method: 'POST'})

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  users.value = response.data ?? []
}

function handleUserSelect(event) {
  selectedUser.value = users.value.find((user) => user.uuid === event.target.value) ?? null
}

async function handleWithdraw() {
  successMessage.value = ''
  errorMessage.value = ''

  const response = await $fetch('/api/wallet/withdraw', {
    method: 'POST',
    body: {
      user_uuid: selectedUser.value?.uuid,
      amount: amount.value,
    }
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  successMessage.value = 'Withdraw was successful'
  selectedUser.value = null

  await fetchUsers()
}

watch(
    () => amount.value,
    (newAmount) => {
      if (newAmount < 0) {
        amount.value = Math.max(amount.value, 0)
      }

      amount.value = Math.min(amount.value, selectedUser.value.balance)
    }
)
</script>

<style scoped lang="scss">

</style>
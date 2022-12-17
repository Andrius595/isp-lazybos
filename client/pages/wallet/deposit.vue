<template>
  <div>
    <AuthenticatedLayout>
      <h1>Deposit Money</h1>
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
      <button class="btn btn-success" @click="handleDeposit">Deposit</button>

    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const errorMessage = ref('')
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

async function handleDeposit() {
  const response = await $fetch('/api/wallet/deposit', {
    method: 'POST',
    body: {
      user_uuid: selectedUser?.uuid,
      amount: amount.value,
    }
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  users.value = response.data ?? []
}

watch(
    () => amount.value,
    (newAmount) => {
      if (newAmount < 0) {
        amount.value = 0
      }
    }
)
</script>

<style scoped lang="scss">

</style>
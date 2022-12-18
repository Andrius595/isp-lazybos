<template>
<div>
    <AuthenticatedLayout>
      <h1>Administrators actions report</h1>
      <div class="d-flex">
        <div><div class="mt-2 mx-5">
          <h5>Select administrator:</h5>
        </div>
          <div class="mx-5" style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
            <select @change="handleSelectAdmin" class="form-select" aria-label="adminSelect">
              <option value="" selected>Select Administrator...</option>
              <option v-for="(admin, index) in admins" :key="index" :value="admin.uuid">{{ admin.email }}</option>
            </select>
          </div>
          <button type="button" style="width: 250px; margin: 20px 45px; border: 2px solid black;
        border-radius: 8px; background-color: lightgreen;" @click="handleShowAdminActions">Show administrator's actions</button>
        </div>
        <div v-if="isShowingTable">
          <div class="mx-auto mt-3" style="width: 450px; border: 2px solid black; border-radius: 8px; padding: 15px 10px 0 10px">
            <h4>Administrator actions list</h4>
            <table class="table">
              <thead style="background-color: #9a9a9a;">
              <tr>
                <th scope="col">Date</th>
                <th scope="col">IP address</th>
                <th scope="col">Action type</th>
              </tr>
              </thead>
              <tbody class="table-secondary">
                <tr v-for="(action, index) in actions" :key="index">
                  <td>{{ action.timestamp}}</td>
                  <td>{{ action.action}}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";

const actions = ref([])

const admins = ref([])
const selectedAdmin = ref(null)
const isShowingTable = ref(false)

await fetchAdmins()

async function fetchAdmins() {
  const res = await $fetch('/api/reports/admins', { method: 'POST' })

  if (!res.status) {
    // TODO show error
    return
  }

  admins.value = res.data
}

function handleSelectAdmin(event) {
  const admin = admins.value.find((user) => user.uuid === event.target.value)
  isShowingTable.value = false

  selectedAdmin.value = admin ?? null
}

async function handleShowAdminActions() {
  const res = await $fetch('/api/reports/logs', {
    method: 'POST',
    body: {
      admin_uuid: selectedAdmin.value.uuid,
    }
  })

  if (!res.status) {
    // TODO show error
    return
  }

  actions.value = res.data
  isShowingTable.value = true
}

</script>

<style scoped lang="scss">

</style>
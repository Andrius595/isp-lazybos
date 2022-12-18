<template>
<div>
    <AuthenticatedLayout>
      <h1>Administratorių veiksmų ataskaitos langas</h1>
      <div class="d-flex">
        <div><div class="mt-2 mx-5">
          <h5>Select administrator:</h5>
        </div>
          <div class="mx-5" style="width: 250px; border: 2px solid black; border-radius: 8px; padding: 5px">
            <select class="form-select" aria-label="adminSelect">
              <option selected>Administrator...</option>
              <option value="admin1">Admin 1</option>
              <option value="admin2">Admin 2</option>
              <option value="admin3">Admin 3</option>
            </select>
          </div>
          <button type="button" style="width: 250px; margin: 20px 45px; border: 2px solid black;
        border-radius: 8px; background-color: lightgreen;" @click="handleShowAdminActions">Show administrator's actions</button>
        </div>
        <div v-if="selectedAdmin">
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
              <tr>
                <td><b>1</b></td>
                <td>Test</td>
                <td>Test</td>
              </tr>
              <tr>
                <td><b>2</b></td>
                <td>Test</td>
                <td>Test</td>
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

const admin = ref([])
const selectedAdmin = ref(null)

function fetchAdmins() {

}

function handleShowAdminActions() {
  selectedAdmin.value = !selectedAdmin.value
}

</script>

<style scoped lang="scss">

</style>
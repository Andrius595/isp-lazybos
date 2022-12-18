<template>
  <div>
    <AuthenticatedLayout>
      <h1>Match edit form</h1>
      <div class="container">
        <div style="margin-bottom:2rem" class="card bg-white mt-5">
          <div class="row">
            <div class="col-sm">
              <div class="card-title text-black">
                <h5 class="text-center py-2 font-weight-bold mb-0 mt-2">Match editing</h5>
              </div>
              <div v-if="match" class="card-body">
                <label for="name">Match name</label>
                <input type="text" name="name" v-model="match.name" class="form-control mb-2">
                <label style="margin-top:1rem" for="outcomeselection">Select match outcome</label>
                <select @change="handleSelectSelection" class="form-select form-select-sm" aria-label=".form-select-sm example">
                  <option value="" selected>Select outcome</option>
                  <option v-for="(selection, index) in match.selections" :key="index" :value="selection.uuid">{{ selection.name }}</option>
                </select>
                <div style="margin-top:1rem;" id="hide" v-if="selectedOutcome" class="col-sm">
                  <div class="row">
                    <div class="col-sm">
                      <label>Odds Home Team</label>
                      <input v-model.number="selectedOutcome.odds_home" step="0.01" type="number" placeholder="" class="form-control mb-2">
                      <label>Odds Away Team</label>
                      <input v-model.number="selectedOutcome.odds_away" step="0.01" type="number" placeholder="" class="form-control mb-2">
                    </div>
                  </div>
                </div>
                <button class="btn btn-success btn-block mt-3" name="createMatch">Save match</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import { useRoute } from "vue-router";
import Routes from "~/types/routes";

const route = useRoute()

const errorMessage = ref('')
const uuid = computed(() => Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)

const match = ref(await fetchData(uuid.value))

const selectedOutcome = ref(null)

async function fetchData(uuid: string) {
  const response = await $fetch('/api/events/all', {method: 'POST'})

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  const match = response.data.find((match) => match.uuid === uuid)

  if (!match) {
    return navigateTo({ name: Routes.Identity.List })
  }

  return match
}


function handleSelectSelection(e){
  const selection = match.value.selections.find((sel) => sel.uuid === e.target.value)

  if (!selection) {
    selectedOutcome.value = null

    return
  }

  selectedOutcome.value = selection
}
</script>


<style scoped lang="scss">

</style>
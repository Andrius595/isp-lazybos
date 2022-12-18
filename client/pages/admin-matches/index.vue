<template>
  <div>
    <AuthenticatedLayout>
      <h1>Match list</h1>
      <div class="container">
        <div class="row">
          <div class="col-lg-7 mx-auto">
            <div class="card border-0 shadow">
              <div class="card-body p-5">
                <div class="table-responsive">
                  <table id="toptable" class="table m-0 table-hover">
                    <thead>
                    <tr>
                      <th scope="col">Match ID.</th>
                      <th scope="col">Match name</th>
                      <th scope="col">Match status</th>
                    </tr>
                    </thead>
                    <tbody>
                    <!-- Generate unique ids for each row to target collapse -->
                    <template v-for="(match, index) in matches" :key="index">
                      <tr data-bs-toggle="collapse" :data-bs-target="`#outcomeRow${index}`">
                        <!-- Rungtynių duomenys -->
                        <th scope="row">{{ index + 1}}</th>
                        <td>{{ match.name }}</td>
                        <td>{{ match.finished ? 'Finished' : 'In Progress' }}</td>

                      </tr>
                      <tr class="collapse" :id="`outcomeRow${index}`">
                        <td id="nohover" colspan="3">
                          <table class="table m-0">
                            <thead>
                            <tr>
                              <th scope="col">Outcome</th>
                              <th scope="col">Winner</th>
                              <th scope="col">Pick winner</th>
                              <th scope="col">Actions</th>
                              <th scope="col">
                                <NuxtLink name="editmatch" class="btn btn-info btn-sm"
                                          :href="getRouteUrl(Routes.AdminMatches.Edit, {id: match.uuid})">Edit match
                                </NuxtLink>
                              </th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr v-for="(selection, selectionIndex) in match.selections" :key="selectionIndex">
                              <td>{{ selection.name}}</td>
                              <td>{{ selection.winner }}</td>
                              <td v-if="selection.winner === 'tbd'">
                                <div class="form-check">
                                  <!-- Keep same name for all radios -->
                                  <input @change="handleSelectionRadio($event, selection)" class="form-check-input" type="radio" :name="`flexRadioDefault${selectionIndex}`"
                                         value="home">
                                  <label class="form-check-label">
                                    Home Team won
                                  </label>
                                </div>
                                <div class="form-check">
                                  <input @change="handleSelectionRadio($event, selection)" class="form-check-input" type="radio" :name="`flexRadioDefault${selectionIndex}`" value="away">
                                  <label class="form-check-label">
                                    Away Team won
                                  </label>
                                </div>
                              </td>
                              <td v-if="selection.winner === 'tbd'">
                                <div v-if="selection.new_winner" class="btn btn-warning btn-sm" @click="handleEndMatch(selection)">End match</div>
                                <div class="btn btn-danger btn-sm" @click="handleTerminateMatch(selection)">Terminate match</div>
                              </td>
                            </tr>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </template>
                    </tbody>
                  </table>

                </div>
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
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";

const matches = ref([])

await fetchMatches()

async function fetchMatches() {
  const response = await $fetch('/api/events/all', {method: 'POST'})
  if (!response.status) {
    // TODO show error

    return
  }

  matches.value = response.data
}

async function handleEndMatch(selection) {
 await $fetch('/api/events/resolve', {
   method: 'POST',
   body: {
     selection_uuid: selection.uuid,
     winner: selection.new_winner,
   }
 })

  await fetchMatches()
}

async function handleTerminateMatch(selection) {
  await $fetch('/api/events/resolve', {
    method: 'POST',
    body: {
      selection_uuid: selection.uuid,
      winner: 'none',
    }
  })

  await fetchMatches()
}

function handleSelectionRadio(event, selection) {
  selection.new_winner = event.target.value
}
</script>

<style scoped lang="scss">
#nohover:hover {
  box-shadow: none;
}
</style>
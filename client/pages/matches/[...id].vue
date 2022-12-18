<template>
  <div>
    <AuthenticatedLayout>
      <h1>Bet on match {{ match?.name }}</h1>
      <div v-if="errorMessage.length" class="text-danger">{{ errorMessage }}</div>
      <div v-if="successMessage.length" class="text-success">{{ successMessage }}</div>
      <div v-if="match" class="d-flex flex-column">
        <label style="margin-top:1rem" for="outcomeselection">Select match outcome</label>
        <select @change="handleSelectSelection" class="form-select form-select-sm" aria-label=".form-select-sm example">
          <option value="" selected>Select outcome</option>
          <option v-for="(selection, index) in match.selections" :key="index" :value="selection.uuid">{{ selection.name }}</option>
        </select>
        <div v-if="selectedOutcome">
          <div class="form-check mt-2">
            <!-- Keep same name for all radios -->
            <input @change="handleSelectionRadio($event)" class="form-check-input" type="radio" :name="`flexRadioDefault`" value="home">
            <label class="form-check-label">
              Home Team won (odds: {{ selectedOutcome.odds_home }})
            </label>
          </div>
          <div class="form-check">
            <input @change="handleSelectionRadio($event, selection)" class="form-check-input" type="radio" :name="`flexRadioDefault`" value="away">
            <label class="form-check-label">
              Away Team won (odds: {{ selectedOutcome.odds_away }})
            </label>
          </div>
          <div class="form-outline mt-2">
            <label class="form-label">Stake amount</label>
            <input v-model.number="selectedOutcome.amount" min="0" type="number" class="form-control" step="0.01">
          </div>

          <div class="btn btn-primary mt-3" @click="handleBetButton">Submit bet</div>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";

const route = useRoute()

const uuid = computed(() => Array.isArray(route.params.id) ? route.params.id[0] : route.params.id)

const errorMessage = ref('')
const successMessage = ref('')

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

  console.log(selection)
  if (!selection) {
    selectedOutcome.value = null

    return
  }

  selectedOutcome.value = selection
}

function handleSelectionRadio(event) {
  selectedOutcome.value.new_winner = event.target.value
}

async function handleBetButton() {
  successMessage.value = ''
  errorMessage.value = ''

  const body = {
    selection_uuid: selectedOutcome.value.uuid,
    stake: selectedOutcome.value.amount,
    winner: selectedOutcome.value.new_winner,
  }
  const response = await $fetch('/api/events/bet', {
    method: 'POST',
    body
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  successMessage.value = 'Bet was placed successfully!'
}
</script>

<style scoped lang="scss">

</style>
<template>
  <div>
    <AuthenticatedLayout>
      <h1>Match data form</h1>
      <div class="container">
        <div style="margin-bottom:2rem" class="card bg-white mt-5">
          <div class="row">
            <div class="col-sm">
              <div class="card-title text-black">
                <h5 class="text-center py-2 font-weight-bold mb-0 mt-2">Match creation</h5>
              </div>
              <div class="card-body">
                <label>Match name</label>
                <input v-model="matchForm.name" type="text" placeholder="" class="form-control mb-3">
                <label>Match start date</label>
                <input v-model="matchForm.begins_at" type="datetime-local" placeholder="" class="form-control mb-3">
                <label>Sport</label>
                <select class="form-select mb-3" @change="handleSportSelect">
                  <option disabled value="">Select sport..</option>ß
                  <option v-for="(sportOption, index) in sportsOptions" :key="index" :value="sportOption">{{ sportOption }}</option>
                </select>
                <div class="row">
                  <div class="col-sm">
                    <button @click='isShowingOutcome = !isShowingOutcome' class="btn btn-primary btn-block mt-3">Add match outcome</button>
                  </div>
                </div>
                <div id="hide" v-if="isShowingOutcome" class="col-sm">
                  <div class="row">
                    <div class="col-sm">
                      <label>Outcome name</label>
                      <input v-model="outcomeForm.name" type="text" placeholder="" class="form-control mb-2">
                      <label>Odds home team</label>
                      <input id ="odds_home" v-model.number="outcomeForm.odds_home" type="number" step="0.01" placeholder="" class="form-control mb-2">
                      <label>Odds away team</label>
                      <input id="odds_away" v-model.number="outcomeForm.odds_away" type="number" step="0.01" placeholder="" class="form-control mb-2">
                      <p><input type="checkbox" id="checkbox" v-model="autoCoff" />
                        <label for="checkbox"> Automatic odds calculation</label>
                      </p>
                      <button @click="addOutcome" class="btn btn-primary btn-block mt-3">Create match outcome</button>
                    </div>
                  </div>
                </div>
                <label style="margin-top:1rem" for="outcomeselection">Match outcomes</label>
                <div class="overflow-auto" style="width:20rem;height:10rem">
                  <ul class="list-group">
                    <li v-for="(outcome, index) in outcomes" :key="index" class="list-group-item">{{ outcome.name }}</li>
                  </ul>
                </div>
                <button class="btn btn-success btn-block mt-3" @click="handleCreateMatch">Create match</button>
              </div>
            </div>
            <div class="col-sm">
              <div class="card-title text-black">
                <h5 class="text-center py-2 font-weight-bold mb-0 mt-2">Team creation</h5>
              </div>
              <div class="card-body">
                <ul class="nav nav-tabs" id="myTab" role="tablist">
                  <li class="nav-item" role="presentation">
                    <button class="nav-link active" id="teamA-tab" data-bs-toggle="tab" data-bs-target="#teamA"
                      type="button" role="tab" aria-controls="teamA" aria-selected="true">Home Team</button>
                  </li>
                  <li class="nav-item" role="presentation">
                    <button class="nav-link" id="teamB-tab" data-bs-toggle="tab" data-bs-target="#teamB" type="button"
                      role="tab" aria-controls="teamB" aria-selected="false">Away Team</button>
                  </li>
                </ul>
                <div class="tab-content">
                  <div class="tab-pane active" id="teamA" role="tabpanel" aria-labelledby="teamA-tab">
                    <label>Team's name</label>
                    <input v-model="homeTeamForm.name" type="text" placeholder="" class="form-control mb-2">
                    <label>Team's country</label>
                    <input v-model="homeTeamForm.country" type="text"  placeholder="" class="form-control mb-2">
                    <label>Team's league</label>
                    <input v-model="homeTeamForm.league" type="text" placeholder="" class="form-control mb-2">
                    <hr />
                    <label>Player's name</label>
                    <input v-model="homePlayerForm.first_name" type="text" placeholder="" class="form-control mb-2">
                    <label>Player's last name</label>
                    <input v-model="homePlayerForm.last_name" type="text" placeholder="" class="form-control mb-2">
<!--                    <label>Player's birth date</label>-->
<!--                    <input v-model="homePlayerForm.birth_date" type="date" placeholder="" class="form-control mb-2">-->
<!--                    <label>Player's nationality</label>-->
<!--                    <input v-model="homePlayerForm.nationality" type="text" name="playerAnationality" placeholder="" class="form-control mb-2">-->
                    <button class="btn btn-success btn-block mt-3" @click="addHomePlayer">Add player to Home team</button>
                    <p style="margin-top:0.5rem;">List of players:</p>
                    <div class="overflow-auto" style="width:20rem;height:10rem">
                      <ul class="list-group">
                        <li v-for="(player, index) in homeTeamForm.players" :key="index" class="list-group-item">{{ `${player.first_name} ${player.last_name}`}}</li>
                      </ul>
                    </div>
                  </div>
                  <div class="tab-pane" id="teamB" role="tabpanel" aria-labelledby="teamB-tab">
                    <label>Team's name</label>
                    <input v-model="awayTeamForm.name" type="text" placeholder="" class="form-control mb-2">
                    <label>Team's country</label>
                    <input v-model="awayTeamForm.country" type="text" placeholder="" class="form-control mb-2">
                    <label>Team's league</label>
                    <input v-model="awayTeamForm.league" type="text" placeholder="" class="form-control mb-2">
                    <hr />
                    <label>Player's name</label>
                    <input v-model="awayPlayerForm.first_name" type="text" placeholder="" class="form-control mb-2">
                    <label>Player's last name</label>
                    <input v-model="awayPlayerForm.last_name" type="text" placeholder="" class="form-control mb-2">
<!--                    <label>Player's birth date</label>-->
<!--                    <input v-model="awayPlayerForm.birth_date" type="date" placeholder="" class="form-control mb-2">-->
<!--                    <label>Player's nationality</label>-->
<!--                    <input v-model="awayPlayerForm.nationality" type="text" placeholder="" class="form-control mb-2">-->
                    <button class="btn btn-success btn-block mt-3" @click="addAwayPlayer">Add player to Away team</button>
                    <p style="margin-top:0.5rem;">List of players:</p>
                    <div class="overflow-auto" style="width:20rem;height:10rem">
                      <ul class="list-group">
                        <li v-for="(player, index) in awayTeamForm.players" :key="index" class="list-group-item">{{ `${player.first_name} ${player.last_name}`}}</li>
                      </ul>
                    </div>
                  </div>
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

const autoCoff = ref<boolean>(false)
const isShowingOutcome = ref<boolean>(false)


const matchForm = ref({
  name: '',
  sport: '',
  begins_at: '',
})


const sportsOptions = ['football', 'basketball']

const outcomeErrorMessage = ref('')
const outcomes = ref([])
const emptyOutcomeForm = {
  name: '',
  odds_home: 0,
  odds_away: 0,
}
const outcomeForm = ref({...emptyOutcomeForm})

function addOutcome() {
  if (outcomeForm.value.name.length === 0) {
    outcomeErrorMessage.value = 'Outcome must have name'

    return
  }

  outcomes.value.push({...outcomeForm.value})
  outcomeForm.value = {...emptyOutcomeForm}
}


const teamErrorMessage = ref('')
const emptyTeamForm = {
  name: '',
  country: '',
  league: '',
  players: []
}

const homeTeamForm = ref({...emptyTeamForm})
const awayTeamForm = ref({...emptyTeamForm})

const playerErrorMessage = ref('')
const emptyPlayerForm = {
  first_name: '',
  last_name: '',
  birth_date: '',
  nationality: '',
}

const awayPlayerForm = ref({...emptyPlayerForm})
const homePlayerForm = ref({...emptyPlayerForm})

function addHomePlayer() {
  playerErrorMessage.value = ''
  const playerForm = homePlayerForm.value

  if (
      playerForm.first_name.length === 0 ||
      playerForm.last_name.length === 0
  ) {
    playerErrorMessage.value = 'Player form fields can\'t be empty'

    return
  }

  homeTeamForm.value.players = [ ...homeTeamForm.value.players, {...homePlayerForm.value}]
  homePlayerForm.value = {...emptyPlayerForm}

  console.log(homeTeamForm.value, awayTeamForm.value)
}

function addAwayPlayer() {
  playerErrorMessage.value = ''
  const playerForm = awayPlayerForm.value

  if (
      playerForm.first_name.length === 0 ||
      playerForm.last_name.length === 0
  ) {
    playerErrorMessage.value = 'Player form fields can\'t be empty'

    return
  }

  awayTeamForm.value.players = [ ...awayTeamForm.value.players, {...awayPlayerForm.value}]
  awayPlayerForm.value = {...emptyPlayerForm}

  console.log(homeTeamForm.value, awayTeamForm.value)
}

function handleSportSelect(event) {
  matchForm.value.sport = event.target.value
}

async function handleCreateMatch() {
  const date = new Date(matchForm.value.begins_at).toISOString()
  console.log('date', date)
  const body = {
    ...matchForm.value,
    begins_at: date,
    home_team: homeTeamForm.value,
    away_team: awayTeamForm.value,
    selections: outcomes.value,
  }

  body.home_team.players = body.home_team.players.map((player) => `${player.first_name} ${player.last_name}`)
  body.away_team.players = body.away_team.players.map((player) => `${player.first_name} ${player.last_name}`)

  console.log(body)

  const response = await $fetch('/api/events/create', { method: 'POST', body })

  if (!response.status) {
    // TODO show error
    return
  }

  return navigateTo({ name: Routes.AdminMatches.List })
}
</script>


<style scoped lang="scss">

</style>
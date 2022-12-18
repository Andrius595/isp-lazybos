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
                <label for="name">Match name</label>
                <input type="text" name="name" placeholder="" class="form-control mb-2">
                <label for="date">Match start date</label>
                <p><input type="date" name="date" placeholder="" class="form-control mb-2"><input type="time"
                    name="time" placeholder="" class="form-control mb-2"></p>
                <label for="sport">Sport</label>
                <select class="form-select mb-3">
                  <option disabled value="">Select sport..</option>
                </select>
                <label for="team1">Team A</label>
                <select class="form-select mb-3">
                  <option disabled value="">Select team A..</option>
                  <!-- Komandos -->
                </select>
                <label for="team2">Team B</label>
                <select class="form-select mb-3">
                  <option disabled value="">Select team B..</option>
                  <!-- Komandos -->
                </select>
                <p><input type="checkbox" id="checkbox" v-model="autoCoff" />
                  <label for="checkbox"> Automatic odds calculation</label>
                </p>
                <div class="row">
                  <div class="col-sm">
                    <button @click='isShowingOutcome = !isShowingOutcome' class="btn btn-primary btn-block mt-3"
                      name="createMatch">Add match outcome</button>
                  </div>
                </div>
                <div id="hide" v-show='isShowingOutcome' class="col-sm">
                  <div class="row">
                    <div class="col-sm">
                      <label for="name">Outcome name</label>
                      <input type="text" name="name" placeholder="" class="form-control mb-2">
                      <label for="name">Odds team A</label>
                      <input type="number" name="name" placeholder="" class="form-control mb-2">
                      <label for="name">Odds team B</label>
                      <input type="number" name="name" placeholder="" class="form-control mb-2">
                      <button @click='isShowingOutcome = !isShowingOutcome' class="btn btn-primary btn-block mt-3"
                        name="createMatch">Create match outcome</button>
                    </div>
                  </div>
                </div>
                <label style="margin-top:1rem" for="outcomeselection">Match outcomes</label>
                <div class="overflow-auto" style="width:20rem;height:10rem">
                  <ul class="list-group">
                    <li class="list-group-item">Outcome1</li>
                    <li class="list-group-item">Outcome2</li>
                    <li class="list-group-item">Outcome3</li>
                    <li class="list-group-item">Outcome4</li>
                    <li class="list-group-item">Outcome5</li>
                  </ul>
                </div>
                <button class="btn btn-success btn-block mt-3" name="createMatch">Create match</button>
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
                      type="button" role="tab" aria-controls="teamA" aria-selected="true">Team A</button>
                  </li>
                  <li class="nav-item" role="presentation">
                    <button class="nav-link" id="teamB-tab" data-bs-toggle="tab" data-bs-target="#teamB" type="button"
                      role="tab" aria-controls="teamB" aria-selected="false">Team B</button>
                  </li>
                </ul>
                <div class="tab-content">
                  <div class="tab-pane active" id="teamA" role="tabpanel" aria-labelledby="teamA-tab">
                    <label for="teamAname">Team's name</label>
                    <input type="text" name="teamAname" placeholder="" class="form-control mb-2">
                    <label for="teamAcountry">Team's country</label>
                    <input type="text" name="teamAcountry" placeholder="" class="form-control mb-2">
                    <label for="teamAleague">Team's league</label>
                    <input type="text" name="teamAleague" placeholder="" class="form-control mb-2">
                    <label for="playerAname">Player's name</label>
                    <input type="text" name="playerAname" placeholder="" class="form-control mb-2">
                    <label for="playerAlastname">Player's last name</label>
                    <input type="text" name="playerAlastname" placeholder="" class="form-control mb-2">
                    <label for="playerAdate">Player's birth date</label>
                    <input type="date" name="playerAdate" placeholder="" class="form-control mb-2">
                    <label for="playerAnationality">Player's nationality</label>
                    <input type="text" name="playerAnationality" placeholder="" class="form-control mb-2">
                    <button class="btn btn-success btn-block mt-3" name="addPlayerTeamA">Add player to Team A</button>
                    <p style="margin-top:0.5rem;">List of players:</p>
                    <div class="overflow-auto" style="width:20rem;height:10rem">
                      <ul class="list-group">
                        <li class="list-group-item">Player1</li>
                        <li class="list-group-item">Player2</li>
                        <li class="list-group-item">Player3</li>
                        <li class="list-group-item">Player4</li>
                        <li class="list-group-item">Player5</li>
                      </ul>
                    </div>
                    <button class="btn btn-success btn-block mt-3" name="createTeamA">Create Team A</button>
                  </div>
                  <div class="tab-pane" id="teamB" role="tabpanel" aria-labelledby="teamB-tab">
                    <label for="teamBname">Team's name</label>
                    <input type="text" name="teamBname" placeholder="" class="form-control mb-2">
                    <label for="teamBcountry">Team's country</label>
                    <input type="text" name="teamAcountry" placeholder="" class="form-control mb-2">
                    <label for="teamBleague">Team's league</label>
                    <input type="text" name="teamAleague" placeholder="" class="form-control mb-2">
                    <label for="playerBname">Player's name</label>
                    <input type="text" name="playerBname" placeholder="" class="form-control mb-2">
                    <label for="playerBlastname">Player's last name</label>
                    <input type="text" name="playerBlastname" placeholder="" class="form-control mb-2">
                    <label for="playerBdate">Player's birth date</label>
                    <input type="date" name="playerBdate" placeholder="" class="form-control mb-2">
                    <label for="playerBnationality">Player's nationality</label>
                    <input type="text" name="playerBnationality" placeholder="" class="form-control mb-2">
                    <button class="btn btn-success btn-block mt-3" name="addPlayerTeamB">Add player to Team B</button>
                    <p style="margin-top:0.5rem;">List of players:</p>
                    <div class="overflow-auto" style="width:20rem;height:10rem">
                      <ul class="list-group">
                        <li class="list-group-item">Player1</li>
                        <li class="list-group-item">Player2</li>
                        <li class="list-group-item">Player3</li>
                        <li class="list-group-item">Player4</li>
                        <li class="list-group-item">Player5</li>
                      </ul>
                    </div>
                    <!-- if player list empty button disabled -->
                    <button class="btn btn-success btn-block mt-3" name="createTeamB">Create Team B</button>
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
import { useRoute, useRouter } from "vue-router";
const router = useRouter()
const route = useRoute()
console.log(router.resolve(route))
//
const autoCoff = ref<boolean>(false)
const isShowingOutcome = ref<boolean>(false)


</script>


<style scoped lang="scss">

</style>
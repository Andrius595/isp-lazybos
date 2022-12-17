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
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="date">Match start date</label>
                    <p><input  type="date" name="date" placeholder="" class="form-control mb-2"><input  type="time" name="time" placeholder="" class="form-control mb-2"></p>
                    <label for="sport">Sport</label>
                    <select   class="form-select mb-3">
                      <option disabled value ="">Select sport..</option>
                    </select>
                    <label for="team1">Team A</label>
                    <select   class="form-select mb-3">
                      <option disabled value ="">Select team A..</option>
                      <!-- Komandos -->
                    </select>
                    <label for="team2">Team B</label>
                    <select   class="form-select mb-3">
                      <option disabled value ="">Select team B..</option>
                      <!-- Komandos -->
                    </select>
                    <p><input type="checkbox" id="checkbox" v-model="autoCoff" />
                    <label for="checkbox"> Automatic odds calculation</label></p>
                    <div class="row">
                      <div class="col-sm">
                        <button @click='isShowingOutcome = !isShowingOutcome' class="btn btn-primary btn-block mt-3" name="createMatch">Add match outcome</button>
                      </div>
                      <div id="hide" v-show='isShowingOutcome' class="col-sm">
                        <div class="row">
                          <div class="col-sm">
                            <label for="name">Outcome name</label>
                            <input  type="text" name="name" placeholder="" class="form-control mb-2">
                            <label for="name">Odds team A</label>
                            <input  type="number" name="name" placeholder="" class="form-control mb-2">
                            <label for="name">Odds team B</label>
                            <input  type="number" name="name" placeholder="" class="form-control mb-2">
                            <button @click='isShowingOutcome = !isShowingOutcome' class="btn btn-primary btn-block mt-3" name="createMatch">Create match outcome</button>
                          </div>
                        </div>
                      </div>

                    </div>
                    <label style="margin-top:1rem"  for="outcomeselection">Match outcome</label>
                    <select style="margin-top:1rem"  class="form-select mb-3" multiple>
                      <option disabled value ="">Select match outcome..</option>
                      <option value ="">test1</option>
                      <option value ="">test2</option>
                    </select>
                    <p style="font-size:0.8rem">Hold down the Ctrl (windows) or Command (Mac) button to select multiple outcomes.</p>
                    <button class="btn btn-success btn-block mt-3" name="createMatch">Create match</button>
                  </div>
                </div>
                <div class="col-sm">
                  <div class="card-title text-black">
                    <h5 class="text-center py-2 font-weight-bold mb-0 mt-2">Team creation</h5>
                  </div>
                  <div class="card-body">
                    <label for="name">Team's name</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="date">Team's country</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="date">Team's league</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <button class="btn btn-success btn-block mt-3" name="createMatch">Add team</button>
                    <p style="margin-top:0.5rem;">List of teams:</p>
                    <div class="overflow-auto" style="width:20rem;height:10rem">
                      <ul class="list-group">
                        <li class="list-group-item">Team1</li>
                        <li class="list-group-item">Team2</li>
                        <li class="list-group-item">Team3</li>
                        <li class="list-group-item">Team4</li>
                        <li class="list-group-item">Team5</li>
                      </ul>
                    </div>
                  </div>
                </div>
                <div class="col-sm">
                  <div class="card-title text-black">
                    <h5 class="text-center py-2 font-weight-bold mb-0 mt-2">Player creation</h5>
                  </div>
                  <div class="card-body">
                    <label for="name">Player's name</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="date">Player's last name</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="date">Player's birth date</label>
                    <input  type="date" name="date" placeholder="" class="form-control mb-2">
                    <label for="date">Player's nationality</label>
                    <input  type="text" name="name" placeholder="" class="form-control mb-2">
                    <label for="team2">Player's team</label>
                    <select class="form-select mb-3">
                      <option disabled value ="">Select team..</option>
                      <!-- Komandos -->
                    </select>
                    <button class="btn btn-success btn-block mt-3" name="createMatch">Add player</button>
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
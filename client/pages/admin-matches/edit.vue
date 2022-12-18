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
              <div class="card-body">
                <label for="name">Match name</label>
                <input type="text" name="name" placeholder="" class="form-control mb-2">
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
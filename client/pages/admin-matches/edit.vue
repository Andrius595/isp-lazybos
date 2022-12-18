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
                <label style="margin-top:1rem" for="outcomeselection">Select match outcome</label>
                <select v-on:change="onChangeSite($event)" class="form-select form-select-sm" aria-label=".form-select-sm example">
                  <option value="null" selected></option>
                  <option value="1">One</option>
                  <option value="2">Two</option>
                  <option value="3">Three</option>
                </select>
                <div style="margin-top:1rem;" id="hide" v-if='isShowingOutcome' class="col-sm">
                  <div class="row">
                    <div class="col-sm">
                      <label for="name">Odds team A</label>
                      <input type="number" name="name" placeholder="" class="form-control mb-2">
                      <label for="name">Odds team B</label>
                      <input type="number" name="name" placeholder="" class="form-control mb-2">
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
import { useRoute, useRouter } from "vue-router";
const router = useRouter()
const route = useRoute()
console.log(router.resolve(route))
//
let isShowingOutcome = ref<boolean>(false);

function onChangeSite(e){
  if(e.target.value != "null")
  {
     isShowingOutcome.value = true;
  }
  else
  {
    isShowingOutcome.value = false;
  }
}
</script>


<style scoped lang="scss">

</style>
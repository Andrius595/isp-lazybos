<template>
  <div>
    <AuthenticatedLayout>
      <h1>Match ending confirmation window</h1>
      <div class="container">
        <div style="margin-bottom:2rem" class="card bg-white mt-5">
          <div class="row">
              <div class="col-sm">
                <h1 style="text-align:center;margin:2rem;">Are you sure you want to end match MATCH_NAME?</h1>
                <div class="row">
                  <div class="col-sm">
                    <button style="margin-left:25rem;margin-bottom:2rem;" class="btn btn-danger btn-block mt-3" name="createMatch">End match</button>
                    <NuxtLink class="btn btn-primary btn-block mt-3" :href="getRouteUrl(Routes.AdminMatches.List)" style="margin-left:15rem;margin-bottom:2rem;"  name="endMatch">Go back</NuxtLink>
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
</script>
<style scoped lang="scss">

</style>
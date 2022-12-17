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
                            <table class="table m-0">
                                <thead>
                                    <tr>
                                        <th scope="col">Match ID.</th>
                                        <th scope="col">Match name</th>
                                        <th scope="col">Time left</th>
                                        <th scope="col">Match status</th>
                                        <th scope="col"></th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <!-- Rungtynių duomenys -->
                                        <th scope="row"></th>
                                        <td></td>
                                        <td></td>
                                        <td></td>
                                        <td>
                                            <!-- Call to action buttons -->
                                            <ul class="list-inline m-0">
                                                <li class="list-inline-item">
                                                    <NuxtLink class="btn btn-warning btn-sm rounded-0" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)" type="button" data-toggle="tooltip" data-placement="top" title="End match"><i class="fa fa-trash"></i>End match</NuxtLink>
                                                </li>
                                                <li class="list-inline-item">
                                                    <NuxtLink class="btn btn-danger btn-sm rounded-0" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)" type="button" data-toggle="tooltip" data-placement="top" title="Terminate match"><i class="fa fa-trash"></i>Terminate match</NuxtLink>
                                                </li>
                                            </ul>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>

                        </div>
                    </div>
                </div>
            </div>
        </div>
      </div>
      <!-- <div class="d-flex">
        <ul class="list-group flex-shrink">
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">Rungtynių užbaigimo patvirtinimas</NuxtLink>
          <NuxtLink class="list-group-item list-group-item-action" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Rungtynių pašalinimo patvirtinimas</NuxtLink>
        </ul>
      </div> -->
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
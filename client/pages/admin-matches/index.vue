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
                            <table id = "toptable" class="table m-0 table-hover">
                                <thead>
                                    <tr>
                                        <th scope="col">Match ID.</th>
                                        <th scope="col">Match name</th>
                                        <th scope="col">Match status</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <!-- Generate unique ids for each row to target collapse -->
                                    <tr data-bs-toggle="collapse" data-bs-target="#outcomeRow" >
                                        <!-- Rungtynių duomenys -->
                                        <th scope="row">1</th>
                                        <td>Name</td>
                                        <td>In progress</td>

                                    </tr>
                                    <tr class="collapse" id="outcomeRow">
                                        <td id="nohover" colspan="3">
                                            <table class=" table m-0 collapse" id="outcomeRow">
                                        <thead>
                                            <tr>
                                                <th scope="col">Outcome</th>
                                                <th scope="col">Winner</th>
                                                <th scope="col">Pick winner</th>
                                                <th scope="col">Actions</th>
                                                <th scope="col">
                                                    <NuxtLink name="editmatch" class="btn btn-info btn-sm" :href="getRouteUrl(Routes.AdminMatches.Edit)">Edit match</NuxtLink>
                                                </th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            <tr>
                                                <td>Outcome name 1</td>
                                                <td>Team A won</td>
                                                <td>
                                                    <div class="form-check">
                                                        <!-- Keep same name for all radios -->
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="flexRadioDefault1">
                                                        <label class="form-check-label" for="flexRadioDefault1">
                                                            Team A won
                                                        </label>
                                                    </div>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="flexRadioDefault1">
                                                        <label class="form-check-label" for="flexRadioDefault1">
                                                            Team B won
                                                        </label>
                                                    </div>
                                                </td>
                                                <td>
                                                    <NuxtLink style="margin-right:1rem; margin-bottom:1rem;" name="endmatch" class="btn btn-warning btn-sm" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">End match</NuxtLink>
                                                    <NuxtLink class="btn btn-danger btn-sm" name="terminatematch" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Terminate match</NuxtLink>
                                                </td>
                                            </tr>
                                            <tr>
                                                <td>Outcome name 1</td>
                                                <td>Team A won</td>
                                                <td>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="checkBoxTeamA">
                                                        <label class="form-check-label" for="checkBoxTeamA">
                                                            Team X won
                                                        </label>
                                                    </div>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="checkBoxTeamB">
                                                        <label class="form-check-label" for="checkBoxTeamB">
                                                            Team S won
                                                        </label>
                                                    </div>
                                                </td>
                                                <td>
                                                    <NuxtLink style="margin-right:1rem; margin-bottom:1rem;" name="endmatch" class="btn btn-warning btn-sm" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">End match</NuxtLink>
                                                    <NuxtLink class="btn btn-danger btn-sm" name="terminatematch" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Terminate match</NuxtLink>
                                                </td>
                                            </tr>
                                                       
                                        </tbody>
                                    </table>
                                        </td>
                                    </tr>
                                    <!-- Generate unique ids for each row to target collapse -->
                                    <tr data-bs-toggle="collapse" data-bs-target="#outcomeRow2" >
                                        <!-- Rungtynių duomenys -->
                                        <th scope="row">1</th>
                                        <td>Name</td>
                                        <td>In progress</td>
                                    </tr>
                                    <tr class="collapse" id="outcomeRow2">
                                        <td id="nohover" colspan="3">
                                            <table class=" table m-0 collapse" id="outcomeRow2">
                                        <thead>
                                            <tr>
                                                <th scope="col">Outcome</th>
                                                <th scope="col">Winner</th>
                                                <th scope="col">Pick winner</th>
                                                <th scope="col">Actions</th>
                                                <th scope="col">
                                                    <NuxtLink name="editmatch" class="btn btn-info btn-sm" :href="getRouteUrl(Routes.AdminMatches.Edit)">Edit match</NuxtLink>
                                                </th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            
                                            <tr>
                                                
                                                <td>Outcome name 1</td>
                                                <td>Team A won</td>
                                                <td>
                                                    <div class="form-check">
                                                        <!-- Keep same name for all radios -->
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="flexRadioDefault1">
                                                        <label class="form-check-label" for="flexRadioDefault1">
                                                            Team A won
                                                        </label>
                                                    </div>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="flexRadioDefault1">
                                                        <label class="form-check-label" for="flexRadioDefault1">
                                                            Team B won
                                                        </label>
                                                    </div>
                                                </td>
                                                <td>
                                                    <NuxtLink style="margin-right:1rem; margin-bottom:1rem;" name="endmatch" class="btn btn-warning btn-sm" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">End match</NuxtLink>
                                                    <NuxtLink class="btn btn-danger btn-sm" name="terminatematch" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Terminate match</NuxtLink>
                                                </td>
                                            </tr>
                                            <tr>
                                                <td>Outcome name 1</td>
                                                <td>Team A won</td>
                                                <td>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="checkBoxTeamA">
                                                        <label class="form-check-label" for="checkBoxTeamA">
                                                            Team X won
                                                        </label>
                                                    </div>
                                                    <div class="form-check">
                                                        <input class="form-check-input" type="radio" name="flexRadioDefault" id="checkBoxTeamB">
                                                        <label class="form-check-label" for="checkBoxTeamB">
                                                            Team S won
                                                        </label>
                                                    </div>
                                                </td>
                                                <td>
                                                    <NuxtLink style="margin-right:1rem; margin-bottom:1rem;" name="endmatch" class="btn btn-warning btn-sm" :href="getRouteUrl(Routes.AdminMatches.EndConfirmation)">End match</NuxtLink>
                                                    <NuxtLink class="btn btn-danger btn-sm" name="terminatematch" :href="getRouteUrl(Routes.AdminMatches.RemoveConfirmation)">Terminate match</NuxtLink>
                                                </td>
                                            </tr>
                                                       
                                        </tbody>
                                    </table>
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
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";


</script>
<style scoped lang="scss">

#nohover:hover {
    box-shadow: none;
}


</style>
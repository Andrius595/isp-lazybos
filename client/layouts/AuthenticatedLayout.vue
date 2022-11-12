<template>
  <div>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container-fluid">
        <NuxtLink class="navbar-brand" :href="getRouteUrl(Routes.Main)">Lažybos</NuxtLink>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav mb-2 mb-lg-0">
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Profile)}" :href="getRouteUrl(Routes.Profile)">Paskyra</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Identity.List)}" :href="getRouteUrl(Routes.Identity.List)">Prašymų sąrašas</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Matches.List)}" :href="getRouteUrl(Routes.Matches.List)">Varžybų sąrašas</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AdminMatches.List)}" :href="getRouteUrl(Routes.AdminMatches.List)">Rungtynių sąrašas</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AdminMatches.Create)}" :href="getRouteUrl(Routes.AdminMatches.Create)">Rungtynių duomenų forma</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Bets.List)}" :href="getRouteUrl(Routes.Bets.List)">Statymų sąrašas</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Wallet.Deposit)}" :href="getRouteUrl(Routes.Wallet.Deposit)">Pinigų įnešimo forma</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AutomaticBets.List)}" :href="getRouteUrl(Routes.AutomaticBets.List)">Automatinių statymų sąrašas</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Wallet.Activate)}" :href="getRouteUrl(Routes.Wallet.Activate)">Įnešimo aktyvavimo forma</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Wallet.Withdraw)}" :href="getRouteUrl(Routes.Wallet.Withdraw)">Išgryninimo forma</NuxtLink>
            </li>
            <li class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Reports.List)}" :href="getRouteUrl(Routes.Reports.List)">Ataskaitų pasirinkimo langas</NuxtLink>
            </li>
          </ul>

        </div>
      </div>
    </nav>
    <div class="container mt-4">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import Routes from "~/types/routes";
import getRouteUrl from "~/utils/getRouteUrl";
import isRouteActive from "~/utils/isRouteActive";
</script>

<style scoped lang="scss">

</style>
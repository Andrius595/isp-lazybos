<template>
  <div>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container-fluid">
        <NuxtLink class="navbar-brand" :href="getRouteUrl(Routes.Main)">Lažybos</NuxtLink>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse justify-content-end" id="navbarSupportedContent">
          <ul class="navbar-nav mb-2 mb-lg-0">
            <li v-if="isUser" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Profile)}" :href="getRouteUrl(Routes.Profile)">Profile</NuxtLink>
            </li>
            <li v-if="isUsersAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Identity.List)}" :href="getRouteUrl(Routes.Identity.List)">Identity Verifications</NuxtLink>
            </li>
            <li v-if="isUserVerified" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Matches.List)}" :href="getRouteUrl(Routes.Matches.List)">Matches</NuxtLink>
            </li>
            <li v-if="isMatchesAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AdminMatches.List)}" :href="getRouteUrl(Routes.AdminMatches.List)">Admin Matches</NuxtLink>
            </li>
            <li v-if="isMatchesAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AdminMatches.Create)}" :href="getRouteUrl(Routes.AdminMatches.Create)">Create Match</NuxtLink>
            </li>
            <li v-if="isUserVerified" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Bets.List)}" :href="getRouteUrl(Routes.Bets.List)">Bets</NuxtLink>
            </li>
            <li v-if="isUsersAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Wallet.Deposit)}" :href="getRouteUrl(Routes.Wallet.Deposit)">Deposit</NuxtLink>
            </li>
            <li v-if="isUserVerified" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.AutomaticBets.List)}" :href="getRouteUrl(Routes.AutomaticBets.List)">Automatic Bets</NuxtLink>
            </li>
            <li v-if="isUsersAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Wallet.Withdraw)}" :href="getRouteUrl(Routes.Wallet.Withdraw)">Withdraw</NuxtLink>
            </li>
            <li v-if="isSalesAdmin" class="nav-item">
              <NuxtLink class="nav-link" :class="{ active: isRouteActive(Routes.Reports.List)}" :href="getRouteUrl(Routes.Reports.List)">Reports</NuxtLink>
            </li>
            <template v-if="isGuest">
              <li class="nav-item">
                <NuxtLink class="nav-link" :href="getRouteUrl(Routes.Auth.Login)">Login</NuxtLink>
              </li>
              <li class="nav-item">
                <NuxtLink class="nav-link" :href="getRouteUrl(Routes.Auth.Register)">Register</NuxtLink>
              </li>
            </template>
            <li v-if="!isGuest" class="nav-item" @click="handleLogout">
              <span class="nav-link">Logout</span>
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
import checkUserRole from "~/utils/checkUserRole";

const user = computed(() => (useCookie('user')).value)

const isGuest = computed(() => !user.value)
const isUser = computed(() => checkUserRole(user.value, 'user'))
const isUserVerified = computed(() => isUser.value && user.value.identitity_verified)
const isUsersAdmin = computed(() => checkUserRole(user.value, 'users'))
const isMatchesAdmin = computed(() => checkUserRole(user.value, 'matches'))
const isSalesAdmin = computed(() => checkUserRole(user.value, 'sales'))

async function handleLogout() {
  await $fetch('/api/auth/logout', { method: 'POST' })

  return navigateTo({ name: Routes.Auth.Login })
}
</script>

<style scoped lang="scss">

</style>
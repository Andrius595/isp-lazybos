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
</script>

<style scoped lang="scss">

</style>
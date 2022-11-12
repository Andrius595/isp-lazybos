<template>
  <div>
    <AuthenticatedLayout>
      <div class="first">
        <div class="second">Red</div>
      </div>
      <div class="second">Blue</div>
      <div class="third">
        <div class="second">Green</div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">

import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";</script>

<style scoped lang="scss">

</style>
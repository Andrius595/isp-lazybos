<template>
  <div>
    <AuthenticatedLayout>
      <div>
        <h1>Request Identity verification</h1>
        <div class="border border-2 border-dark rounded p-3">
          <span v-if="errorMessage.length" class="text-danger">{{ errorMessage }}</span>
          <div class="mb-3">
            <label for="formFile" class="form-label">Upload ID Photo</label>
            <input @change="handleIdPhotoChange" class="form-control" type="file" id="formFile">
          </div>
          <div class="mb-3">
            <label for="formFile" class="form-label">Upload Portrait Photo</label>
            <input @change="handlePortraitPhotoChange" class="form-control" type="file" id="formFile">
          </div>

          <button @click="handleRequest" class="btn btn-success">Request verification</button>
        </div>
      </div>
    </AuthenticatedLayout>
  </div>
</template>

<script setup lang="ts">
import AuthenticatedLayout from "~/layouts/AuthenticatedLayout.vue";
import Routes from "~/types/routes";

const errorMessage = ref('')

const formData = ref({
  id_photo: '',
  portrait_photo: '',
})


async function handleRequest() {
  const response = await $fetch('/api/identity-verifications/create-request', {
    method: 'POST',
    body: formData.value,
  })

  if (!response.status) {
    errorMessage.value = response.message

    return
  }

  return navigateTo({name: Routes.Profile})
}

async function handleIdPhotoChange(event) {
  formData.value.id_photo = await toBase64(event.target.files[0])
}

async function handlePortraitPhotoChange(event) {
  formData.value.portrait_photo = await toBase64(event.target.files[0])
}

const toBase64 = file => new Promise((resolve, reject) => {
  const reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = () => resolve(reader.result);
  reader.onerror = error => reject(error);
});
</script>


<style scoped lang="scss">

</style>
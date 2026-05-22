<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { usePocketbaseStore } from '@/stores/pb'
import NewForm from '@/views/NewForm.vue'

const { user } = usePocketbaseStore()
const isLoaded = ref(false)
const newStorageKey = `formkit-n-u_${user?.username}-f`
const values = ref({})
onMounted(async () => {
  const cache = localStorage.getItem(newStorageKey)
  if (cache) {
    values.value = JSON.parse(cache).data
  } else {
    values.value = {}
  }
  isLoaded.value = true
})
</script>

<template>
  <ProgressBar v-if="!isLoaded" />
  <NewForm v-if="isLoaded" :values="values"></NewForm>
</template>

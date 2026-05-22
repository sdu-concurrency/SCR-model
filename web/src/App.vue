<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { usePocketbaseStore } from '@/stores/pb'
import Toast from 'primevue/toast'
import { storeToRefs } from 'pinia'
import ConfirmDialog from 'primevue/confirmdialog'

const pocketbaseStore = usePocketbaseStore()
const { isBusy } = storeToRefs(pocketbaseStore)
</script>

<template>
  <TheHeader></TheHeader>
  <Suspense>
    <div class="container mx-auto mt-4 px-4 md:px-10 print:m-0 print:max-w-full">
      <RouterView />
    </div>
  </Suspense>

  <Dialog
    v-model:visible="isBusy"
    modal
    :pt="{
      root: 'border-none',
      mask: {
        style: 'backdrop-filter: blur(2px)'
      }
    }"
  >
    <template #container=""> <ProgressSpinner /></template>
  </Dialog>
  <Toast position="bottom-right" group="br" />
  <ConfirmDialog></ConfirmDialog>
</template>

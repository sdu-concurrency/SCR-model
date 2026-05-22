<script setup>
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { usePocketbaseStore } from '@/stores/pb'
import { useAppStore } from '@/stores/app'

const pocketbaseStore = usePocketbaseStore()
const { setUserDialogVisibility } = useAppStore()

const { isLoggedIn } = storeToRefs(pocketbaseStore)
const severity = computed(() => {
  return isLoggedIn.value ? 'success' : 'secondary'
})
const outlined = computed(() => {
  return isLoggedIn.value ? false : true
})
function openDialog() {
  setUserDialogVisibility(true)
}
</script>

<template>
  <Button
    icon="pi pi-user"
    :severity="severity"
    rounded
    :outlined="outlined"
    aria-label="User"
    @click="openDialog"
    size="large"
    class="!text-light-gray !border-light-gray"
  />
  <UserDialog></UserDialog>
</template>

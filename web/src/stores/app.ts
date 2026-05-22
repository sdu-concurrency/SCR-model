import { ref } from 'vue'
import { defineStore } from 'pinia'
export const useAppStore = defineStore('app', () => {
  const userDialogVisibility = ref<boolean>(false)
  async function setUserDialogVisibility(val: boolean) {
    userDialogVisibility.value = val
  }
  return { setUserDialogVisibility, userDialogVisibility }
})

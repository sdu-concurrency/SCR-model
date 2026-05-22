<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { usePocketbaseStore } from '@/stores/pb'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
const router = useRouter()
const pocketbaseStore = usePocketbaseStore()
const { isLoggedIn, user } = storeToRefs(pocketbaseStore)
const { api, logout } = pocketbaseStore
const toast = useToast()
const password = ref('')

const confirm_password = ref('')

async function onSubmit() {
  try {
    await api.send('/auth-set', {
      method: 'POST',
      body: {
        password: password.value
      }
    })
    toast.add({
      severity: 'success',
      summary: 'Password reset',
      detail: 'Please re-login with your new password',
      group: 'br',
      life: 2000
    })

    await logout()
    await router.push({
      name: 'home',
      replace: true
    })
  } catch (e: any) {
    console.error(e)
    toast.add({
      severity: 'error',
      summary: e.message,
      detail: e.response?.data?.password?.message,
      group: 'br',
      life: 2000
    })
  }
}
</script>

<template>
  <Card class="w-2/3 mx-auto">
    <template #title>Change password</template>
    <template #content>
      <form @submit.prevent="onSubmit">
        <div class="grid grid-cols-3 grid-rows-2 mb-2 gap-4 content-center items-center">
          <label for="password" class="font-semibold">{{
            $t('change_password.password_label')
          }}</label>
          <Password
            id="password"
            name="password"
            v-model="password"
            :feedback="false"
            class="col-span-2"
            autocomplete="off"
            size="small"
            :fluid="true"
          />
          <label for="confirm_password" class="font-semibold">{{
            $t('change_password.password_confirm_label')
          }}</label>
          <div class="col-span-2">
            <Password
              id="confirm_password"
              name="confirm_password"
              v-model="confirm_password"
              :feedback="false"
              autocomplete="off"
              size="small"
              :fluid="true"
              :invalid="confirm_password !== password || password.length < 5"
            />
            <p v-show="confirm_password !== password" class="text-sm text-red-700">
              Password does not match
            </p>
            <p v-show="password.length > 0 && password.length < 5" class="text-sm text-red-700">
              Password must be at least 5 characters
            </p>
          </div>
        </div>
        <div class="flex flex-row-reverse gap-2">
          <Button
            type="submit"
            label="Submit"
            :disabled="password === '' || confirm_password !== password || password.length < 5"
          ></Button>
        </div>
      </form>
    </template>
  </Card>
</template>

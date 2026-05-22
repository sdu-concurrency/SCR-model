<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { usePocketbaseStore } from '@/stores/pb'
import { useRouter } from 'vue-router'
const router = useRouter()
const pocketbaseStore = usePocketbaseStore()
const { adminLogin } = pocketbaseStore

const { isLoggedIn, user } = storeToRefs(pocketbaseStore)

const username = ref('')

const password = ref('')

onMounted(() => {
  if (isLoggedIn.value && user.value?.collectionName == '_superusers') {
    router.push({ name: 'session-create', replace: true })
  }
})

async function onLoginClicked() {
  await adminLogin(username.value, password.value)
  if (isLoggedIn.value) {
    router.push({ name: 'session-create', replace: true })
  }
}
</script>

<template>
  <Card class="w-full md:w-2/3 mx-auto">
    <template #title><span class="text-xl md:text-2xl">Admin panel</span></template>
    <template #content>
      <form @submit.prevent="onLoginClicked">
        <div
          class="flex flex-col md:grid md:grid-cols-3 md:grid-rows-2 mb-2 gap-4 content-center items-start md:items-center"
        >
          <label for="username" class="font-semibold">{{ $t('user_dialog.username') }}</label>
          <InputText
            id="username"
            name="username"
            v-model="username"
            class="w-full md:col-span-2"
            autocomplete="off"
          />
          <label for="password" class="font-semibold">{{ $t('user_dialog.password') }}</label>
          <Password
            id="password"
            name="password"
            v-model="password"
            :feedback="false"
            class="w-full md:col-span-2"
            autocomplete="off"
            size="small"
            :fluid="true"
          />
        </div>
        <div class="flex flex-col md:flex-row-reverse gap-2">
          <Button type="submit" label="Login" class="w-full md:w-auto"></Button>
        </div>
      </form>
    </template>
  </Card>
</template>

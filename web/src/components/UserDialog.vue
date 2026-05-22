<script setup>
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { usePocketbaseStore } from '@/stores/pb'
import { useAppStore } from '@/stores/app'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
const router = useRouter()
const pocketbaseStore = usePocketbaseStore()
const appStore = useAppStore()
const { t } = useI18n()

const { login, logout } = pocketbaseStore

const { setUserDialogVisibility } = appStore

const { userDialogVisibility } = storeToRefs(appStore)
const { isLoggedIn, user, session } = storeToRefs(pocketbaseStore)

const username = ref('')

const password = ref('')
function closeDialog() {
  setUserDialogVisibility(false)
}
async function onLoginClicked() {
  await login(username.value, password.value)
  closeDialog()
  if (isLoggedIn.value && user.value.admin) {
    router.push({ name: 'session-admin', replace: true })
  } else if (isLoggedIn.value && user.value.role && user.value.role === 'super') {
    router.push({
      name: 'session-overview',
      params: { sessionName: session.value.name },
      replace: true
    })
  } else if (isLoggedIn.value && user.value.role && user.value.role === 'normal') {
    router.push({
      name: 'session-form-welcome',
      params: { sessionName: session.value.name },
      replace: true
    })
  }
}
async function onLogoutClicked() {
  await logout()
  closeDialog()
  await router.push({
    name: 'home',
    replace: true
  })
}

const isAdmin = computed(() => {
  return user.value.collectionName === '_superusers'
})
const dialogHeader = computed(() => {
  if (isLoggedIn.value) {
    if (isAdmin.value) {
      // admin
      return t('user_dialog.username') + ': ' + user.value.email
    }
    return t('user_dialog.username') + ': ' + user.value.username
  } else {
    return t('user_dialog.header_login')
  }
})
</script>

<template>
  <Dialog
    v-model:visible="userDialogVisibility"
    modal
    :header="dialogHeader"
    :style="{ width: isLoggedIn ? '40rem' : '30rem' }"
  >
    <div v-if="isLoggedIn && isAdmin" class="">
      <div class="grid grid-cols-2 mb-2 content-center items-center gap-2">
        <label class="font-semibold">{{ $t('user_dialog.role') }}:</label>
        <span>Admin</span>
        <!-- <label class="font-semibold">{{ $t('user_dialog.form_status') }}:</label>
        <span>{{ user.status }}</span> -->
      </div>
      <div class="flex flex-row-reverse gap-2">
        <Button type="button" label="Logout" @click="onLogoutClicked"></Button>
      </div>
    </div>
    <div v-else-if="isLoggedIn">
      <div class="grid grid-cols-2 mb-2 content-center items-center gap-2">
        <label class="font-semibold">{{ $t('user_dialog.session') }}:</label>
        <span>{{ session.name }}</span>
        <label class="font-semibold">{{ $t('user_dialog.role') }}:</label>
        <span>{{ user.role }}</span>
        <label class="font-semibold">{{ $t('user_dialog.form_status') }}:</label>
        <span>{{ user.status }}</span>
      </div>
      <div class="flex flex-row-reverse gap-2">
        <Button type="button" label="Logout" @click="onLogoutClicked"></Button>
      </div>
    </div>
    <div v-else>
      <form @submit.prevent="onLoginClicked">
        <div class="grid grid-cols-3 grid-rows-2 mb-2 content-center items-center gap-y-2">
          <label for="username" class="font-semibold">{{ $t('user_dialog.username') }}</label>
          <InputText
            id="username"
            name="username"
            v-model="username"
            class="col-span-2"
            autocomplete="off"
          />
          <label for="password" class="font-semibold">{{ $t('user_dialog.password') }}</label>
          <Password
            id="password"
            name="password"
            v-model="password"
            :feedback="false"
            autocomplete="off"
            class="col-span-2"
            :fluid="true"
          />
        </div>
        <div class="flex flex-row-reverse gap-2">
          <Button type="submit" label="Login"></Button>
        </div>
        <span class="text-sm font-light"
          >In case of forgotten password, please contact the admin to reset the password</span
        >
      </form>
    </div>
  </Dialog>
</template>

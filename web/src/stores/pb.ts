import { computed, onMounted, ref } from 'vue'
import { defineStore } from 'pinia'
import PocketBase, {
  type ListResult,
  ClientResponseError,
  type RecordModel,
  type RecordListOptions
} from 'pocketbase'
import { config } from '@/config'
import { useToast } from 'primevue/usetoast'
import type { ToastServiceMethods } from 'primevue/toastservice'
import { id } from '@formkit/i18n'
export const usePocketbaseStore = defineStore('pocketbase', () => {
  const api = new PocketBase(config.apiBaseUrl)
  const isBusy = ref(false)

  const user = ref<RecordModel | undefined>(undefined)
  const session = ref<RecordModel | undefined>(undefined)
  let toast: ToastServiceMethods | undefined

  const isLoggedIn = computed(() => {
    return user.value !== undefined
  })

  onMounted(async () => {
    toast = useToast()
    const cookie = localStorage.getItem('pocketbase_auth')
    if (cookie) {
      const storageObject = JSON.parse(cookie)
      api.authStore.save(storageObject.token, storageObject.record)
      if (api.authStore.isValid) {
        user.value = storageObject.record
      } else {
        try {
          let userRes
          if (api.authStore.isSuperuser) {
            userRes = await api.collection('_superusers').authRefresh()
            user.value = userRes.record
          } else {
            userRes = await api.collection('app_users').authRefresh()
            user.value = userRes.record
          }
        } catch (e) {
          await logout()
        }
      }
      if (!api.authStore.isSuperuser) {
        try {
          const sessionRes = await api.collection('sessions').getOne(user.value?.session)
          session.value = sessionRes
        } catch (e) {
          await logout()
        }
      }
    } else {
      await logout()
    }

    isBusy.value = false
  })

  async function fetchSession() {
    if (!api.authStore.isSuperuser) {
      const sessionRes = await api.collection('sessions').getOne(user.value?.session)
      session.value = sessionRes
      return sessionRes
    } else {
      console.error('fetch session error, user is not authenticated')
      return null
    }
  }

  async function signup(data: any) {
    isBusy.value = true
    try {
      await api.collection('app_users').create(data)
      return true
    } catch (error) {
      const err = new ClientResponseError(error)
      toast!.add({
        severity: 'error',
        summary: 'Error',
        detail: err.data.message,
        group: 'br',
        life: 2000
      })
    } finally {
      isBusy.value = false
    }
  }

  async function login(identifier: string, password: string) {
    isBusy.value = true
    try {
      const userRes = await api.collection('app_users').authWithPassword(identifier, password)
      const sessionRes = await api.collection('sessions').getOne(userRes.record.session)
      user.value = userRes.record
      session.value = sessionRes
    } catch (error) {
      const err = new ClientResponseError(error)
      toast!.add({
        severity: 'error',
        summary: 'Error',
        detail: err.data.message,
        group: 'br',
        life: 2000
      })
    } finally {
      isBusy.value = false
    }
  }

  async function adminLogin(identifier: string, password: string) {
    isBusy.value = true
    try {
      const userRes = await api.collection('_superusers').authWithPassword(identifier, password)
      user.value = userRes.record
    } catch (error) {
      const err = new ClientResponseError(error)
      toast!.add({
        severity: 'error',
        summary: 'Error',
        detail: err.data.message,
        group: 'br',
        life: 2000
      })
    } finally {
      isBusy.value = false
    }
  }

  async function logout() {
    api.authStore.clear()

    user.value = undefined
    session.value = undefined
    Object.keys(localStorage)
      .filter((key) => key.startsWith(`formkit-u_`) || key.startsWith(`formkit-n-u_`))
      .forEach((key) => localStorage.removeItem(key))
  }

  async function getCollection<T>(
    collection: string,
    params?: {
      from?: number
      to?: number
      query?: RecordListOptions
    }
  ): Promise<ListResult<T> | undefined> {
    isBusy.value = true
    try {
      const res = await api
        .collection(collection)
        .getList<T>(params?.from || 1, params?.to || 10, params?.query || {})
      return res
    } catch (error) {
      const err = new ClientResponseError(error)
      toast!.add({
        severity: 'error',
        summary: 'Error',
        detail: err.data.message,
        group: 'br',
        life: 2000
      })
    } finally {
      isBusy.value = false
    }
  }

  return {
    api,
    isBusy,
    isLoggedIn,
    user,
    session,
    login,
    adminLogin,
    signup,
    logout,
    getCollection,
    fetchSession
  }
})

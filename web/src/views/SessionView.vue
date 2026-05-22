<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import Divider from 'primevue/divider'
import Card from 'primevue/card'
import Breadcrumb from 'primevue/breadcrumb'
import Panel from 'primevue/panel'
import { ClientResponseError } from 'pocketbase'
import { useToast } from 'primevue/usetoast'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const toast = useToast()
const router = useRouter()

import { usePocketbaseStore } from '@/stores/pb'
const { api } = usePocketbaseStore()
const props = defineProps({
  sessionName: String
})

const state = reactive<any>({ sessionInfo: null })

const fetchData = async () => {
  state.sessionInfo = await api
    .collection('sessions')
    .getFirstListItem(`name = "${props.sessionName}"`, {
      expand: 'users_token_via_session'
    })
}

const superUser = computed(() => {
  if (!state.sessionInfo) {
    return null
  }
  return state.sessionInfo.expand.users_token_via_session.filter((e: any) => e.role === 'super')[0]
})
const normalUsers = computed(() => {
  if (!state.sessionInfo) {
    return null
  }
  return state.sessionInfo.expand.users_token_via_session.filter((e: any) => e.role === 'normal')
})

const hasFinished = computed(() => {
  if (normalUsers.value.find((e: any) => e.status === 'done')) {
    return true
  }
  return false
})

const toPhase3 = async function () {
  try {
    const res = await api.collection('sessions').update(state.sessionInfo.id, { current_step: '3' })
    router.go(0)
  } catch (e) {
    const err = new ClientResponseError(e)
    if (err.status !== 404) {
      toast.add({
        severity: 'error',
        summary: 'error',
        detail: err.data.message + ' , ' + t('toast_contact_admin'),
        group: 'br',
        life: 2000
      })
    }
  }
}

const toPhase2 = async function () {
  try {
    const res = await api.collection('sessions').update(state.sessionInfo.id, { current_step: '2' })
    router.go(0)
  } catch (e) {
    const err = new ClientResponseError(e)
    if (err.status !== 404) {
      toast.add({
        severity: 'error',
        summary: 'error',
        detail: err.data.message + ' , ' + t('toast_contact_admin'),
        group: 'br',
        life: 2000
      })
    }
  }
}

await fetchData()
const breadcrumbHome = ref({
  icon: 'pi pi-home',
  route: '/session/' + props.sessionName
})
const breadcrumbItems = ref([])

const phase2Users = computed(() => {
  if (state.sessionInfo.current_step === 2 || state.sessionInfo.current_step === '2') {
    return normalUsers.value
  } else {
    return normalUsers.value.filter((e: any) => e.status === 'done')
  }
})
</script>

<template>
  <Breadcrumb :home="breadcrumbHome" :model="breadcrumbItems">
    <template #item="{ item, props }">
      <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
        <a :href="href" v-bind="props.action" @click="navigate">
          <span :class="[item.icon, 'text-color']" />
          <span class="text-primary font-semibold">{{ item.label }}</span>
        </a>
      </router-link>
      <a v-else :href="item.url" :target="item.target" v-bind="props.action">
        <span class="text-surface-700 dark:text-surface-0/80">{{ item.label }}</span>
      </a>
    </template></Breadcrumb
  >
  <Card class="mt-1">
    <template #title
      ><h1 class="text-xl md:text-2xl break-words">
        {{ state.sessionInfo.name }}
      </h1></template
    >
    <template #content>
      <PhaseDisplay
        class="text-center mb-4"
        :currentPhase="state.sessionInfo.current_step"
      ></PhaseDisplay>
      <div class="space-y-2 text-sm md:text-base">
        <p class="break-words">
          {{ $t('session_view.description') }}: {{ state.sessionInfo.description }}
        </p>
        <p class="break-words">
          {{ $t('session_view.email') }}: {{ state.sessionInfo.contact_email }}
        </p>
        <p class="break-words">{{ $t('session_view.tel') }}: {{ state.sessionInfo.contact_tel }}</p>
      </div>
    </template></Card
  >

  <Divider />

  <Card class="mt-1">
    <template #title
      ><h1 class="text-xl md:text-2xl">{{ $t('session_view.survey_status_header') }}</h1></template
    >
    <template #content>
      <div class="overflow-x-auto -mx-4 md:mx-0">
        <SessionUsersOrganizationChart
          :normalUsers="normalUsers"
          :superUser="superUser"
          :sessionName="state.sessionInfo.name"
        ></SessionUsersOrganizationChart>
      </div> </template
  ></Card>
  <Divider />

  <Panel
    v-if="state.sessionInfo.current_step === '3' || state.sessionInfo.current_step === 'completed'"
    toggleable
  >
    <template #header>
      <div class="flex items-center gap-2">
        <i v-if="state.sessionInfo.current_step === '3'" class="pi pi-arrow-down-right"></i>
        <i
          v-else-if="state.sessionInfo.current_step === 'completed'"
          class="pi pi-check"
          style="color: green"
        ></i>
        <p class="text-1xl font-bold">{{ $t('session_view.phase', { current: '3' }) }}</p>
      </div>

      <RouterLink
        :to="{ name: 'session-analysis3', params: { sessionName: state.sessionInfo.name } }"
        class="text-blue-600 visited:text-purple-600 ml-auto mr-1"
        >{{ $t('session_view.summary') }}</RouterLink
      >
    </template>

    <Message
      v-if="state.sessionInfo.current_step === '3'"
      severity="info"
      closable
      :pt="{
        text: {
          class: ['text-sm md:text-base', 'leading-none', 'font-medium', 'w-full']
        }
      }"
      :ptOptions="{ mergeSections: true, mergeProps: true }"
    >
      <div
        class="w-full flex flex-col md:flex-row items-start md:items-center gap-2 md:gap-0 justify-between"
      >
        <span class="text-sm md:text-base">{{ $t('session_view.to_phase_2_desc') }}</span>
        <Button
          :label="$t('session_view.to_phase_2_btn')"
          severity="secondary"
          size="small"
          icon="pi pi-backward"
          iconPos="left"
          class="w-full md:w-auto text-xs md:text-sm"
          @click="toPhase2"
        /></div
    ></Message>

    <SessionUserList :data="[superUser]" :sessionName="state.sessionInfo.name"></SessionUserList>
  </Panel>

  <Panel
    toggleable
    :collapsed="
      state.sessionInfo.current_step === '3' || state.sessionInfo.current_step === 'completed'
    "
  >
    <template #header>
      <div class="flex items-center gap-2">
        <i v-if="state.sessionInfo.current_step === '2'" class="pi pi-arrow-down-right"></i>
        <i v-else class="pi pi-check" style="color: green"></i>
        <p class="text-1xl font-bold">{{ $t('session_view.phase', { current: '2' }) }}</p>
      </div>

      <RouterLink
        :to="{ name: 'session-analysis2', params: { sessionName: state.sessionInfo.name } }"
        class="text-blue-600 visited:text-purple-600 ml-auto mr-1"
        >{{ $t('session_view.summary') }}</RouterLink
      >
    </template>
    <Message
      v-if="state.sessionInfo.current_step === '2' && hasFinished"
      severity="info"
      closable
      :pt="{
        text: {
          class: ['text-sm md:text-base', 'leading-none', 'font-medium', 'w-full']
        }
      }"
      :ptOptions="{ mergeSections: true, mergeProps: true }"
    >
      <div
        class="w-full flex flex-col md:flex-row items-start md:items-center gap-2 md:gap-0 justify-between"
      >
        <span class="text-sm md:text-base">{{ $t('session_view.to_phase_3_desc') }}</span>
        <Button
          :label="$t('session_view.to_phase_3_btn')"
          size="small"
          icon="pi pi-check"
          iconPos="left"
          class="w-full md:w-auto text-xs md:text-sm"
          @click="toPhase3"
        /></div
    ></Message>

    <SessionUserList :data="phase2Users" :sessionName="state.sessionInfo.name"></SessionUserList>
  </Panel>

  <p class="text-sm">
    {{
      $t('session_view.last_updated', {
        date: new Date(state.sessionInfo.updated).toLocaleString()
      })
    }}
  </p>
  <!-- <Inplace :closable="true" class="my-6">
    <template #display> Debugging Data </template>
    <template #content>
      <pre>{{ state.sessionInfo }}</pre>
    </template>
  </Inplace> -->
</template>

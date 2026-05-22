<script setup>
import { computed, ref } from 'vue'

import { useConfirm } from 'primevue/useconfirm'
import { useRouter } from 'vue-router'
import { usePocketbaseStore } from '@/stores/pb'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const toast = useToast()
const { api, session } = usePocketbaseStore()
const confirm = useConfirm()
const router = useRouter()
const props = defineProps({
  user: Object,
  sessionName: String
})

const token = ref(props.user.token)
const token_valid_until = ref(props.user.token_valid_until)

const is_token_valid = computed(() => {
  return new Date(token_valid_until.value) > new Date()
})
const isLoading = ref(false)

const generateNewToken = async (user) => {
  isLoading.value = true
  const res = await api.send('/token', { method: 'POST', query: { u: user } })
  token.value = res.token
  token_valid_until.value = res.valid_until
  isLoading.value = false
}

const copyCredentialToClipBoard = (event, link) => {
  const clipboardData =
    event.clipboardData ||
    window.clipboardData ||
    event.originalEvent?.clipboardData ||
    navigator.clipboard

  clipboardData.writeText(`${link}`)
  toast.add({
    severity: 'success',
    summary: t('copy_user_password.summary'),
    detail: t('copy_user_password.detail', { username: props.user.username }),
    group: 'br',
    life: 2000
  })
}

const fullPathURL = computed(() => {
  return `${window.location.origin}/session/${props.sessionName}/survey?u=${props.user.username}&t=${token.value}`
})
const fullPathSummaryURL = computed(() => {
  return `${window.location.origin}/session/${props.sessionName}/survey/${props.user.username}/view?u=${props.user.username}&t=${token.value}`
})

const confirm_to_normal = (routerParams) => {
  return function () {
    confirm.require({
      message: 'Are you sure you want to proceed to normal user account?',
      header: 'Confirmation',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptProps: {
        label: 'Proceed'
      },
      accept: () => {
        router.push(routerParams)
      },
      reject: () => {}
    })
  }
}
</script>

<template>
  <div class="flex flex-col gap-2 text-sm md:text-base">
    <div class="flex flex-col md:flex-row gap-2 md:gap-4 md:items-center">
      <span class="font-bold">{{ $t('user_dialog.username') }}:</span>
      <span class="break-all"> {{ props.user.username }}</span>
    </div>

    <div class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4">
      <span class="font-bold">{{ $t('session_user.survey_status_label') }}:</span>

      <div class="flex flex-row items-center gap-2">
        <FormFillStatusBadge :status="props.user.status"></FormFillStatusBadge>
        <RouterLink
          v-if="props.user.status === 'done'"
          :to="{
            name: 'session-form-view',
            params: { sessionName: props.sessionName, username: props.user.username }
          }"
          class="text-blue-600 visited:text-purple-600 text-sm"
        >
          {{ $t('session_user.summary_label') }}
        </RouterLink>
      </div>
    </div>

    <div
      v-if="props.user.role === 'normal'"
      class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4"
    >
      <span class="font-bold">Token Status:</span>

      <div class="flex flex-col md:flex-row md:items-center gap-2">
        <div>
          <TokenStatusBadge :is_valid="is_token_valid"></TokenStatusBadge>
        </div>
        <span v-if="is_token_valid" class="text-xs italic"
          >Until:{{ new Date(token_valid_until).toLocaleString('en-GB') }}</span
        >
        <Button
          type="button"
          :label="'Reload new Token'"
          icon="pi pi-refresh"
          size="small"
          class="w-full md:w-auto text-xs"
          :disabled="isLoading"
          @click="generateNewToken(props.user.id)"
        />
      </div>
    </div>

    <div v-if="props.user.role === 'super'">
      <div
        v-if="props.user.status !== 'done' && session.current_step === '3'"
        class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4"
      >
        <span class="font-bold">Link:</span>

        <RouterLink
          :to="{
            name: 'session-form-welcome',
            params: { sessionName: props.sessionName }
          }"
          class="text-blue-600 visited:text-purple-600 break-all text-sm"
        >
          To Phase 3 survey
        </RouterLink>
      </div>
    </div>
    <div v-else-if="props.user.role === 'normal'">
      <div v-if="props.user.status !== 'done'" class="flex flex-col gap-2">
        <div class="flex flex-col gap-2">
          <span class="font-bold">Link:</span>

          <span
            @click="
              confirm_to_normal({
                name: 'session-form-welcome',
                params: { sessionName: props.sessionName },
                query: { u: props.user.username, t: token }
              })
            "
            class="text-blue-600 visited:text-purple-600 cursor-pointer break-all text-xs md:text-sm"
          >
            {{ fullPathURL }}
          </span>

          <Button
            type="button"
            :label="$t('copy_user_password.button_label')"
            icon="pi pi-clipboard"
            severity="success"
            size="small"
            class="w-full md:w-auto"
            @click="(e) => copyCredentialToClipBoard(e, fullPathURL)"
          />
        </div>
      </div>
      <div v-else-if="props.user.status === 'done'" class="flex flex-col gap-2">
        <div class="flex flex-col gap-2">
          <span class="font-bold">Link to summary with token (limited time):</span>

          <span
            @click="
              confirm_to_normal({
                name: 'session-form-view',
                params: { sessionName: props.sessionName, username: props.user.username },
                query: { u: props.user.username, t: token }
              })
            "
            class="text-blue-600 visited:text-purple-600 cursor-pointer break-all text-xs md:text-sm"
          >
            {{ fullPathSummaryURL }}
          </span>

          <Button
            type="button"
            :label="$t('copy_user_password.button_label')"
            icon="pi pi-clipboard"
            severity="success"
            size="small"
            class="w-full md:w-auto"
            @click="(e) => copyCredentialToClipBoard(e, fullPathSummaryURL)"
          />
        </div>
      </div>
    </div>
    <div v-else></div>
  </div>
</template>

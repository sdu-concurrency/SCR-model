<script setup>
import Inplace from 'primevue/inplace'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const props = defineProps({
  username: String,
  password: String,
  link: String,
  status: String
})

const toast = useToast()
const fullPathURL = `${window.location.origin}${props.link}`
const copyCredentialToClipBoard = (event, username, password) => {
  const clipboardData =
    event.clipboardData ||
    window.clipboardData ||
    event.originalEvent?.clipboardData ||
    navigator.clipboard
  clipboardData.writeText(`Link: ${fullPathURL}
${t('copy_user_password.username')}: ${username}
${t('copy_user_password.password')}: ${password}`)
  toast.add({
    severity: 'success',
    summary: t('copy_user_password.summary'),
    detail: t('copy_user_password.detail', { username }),
    group: 'br',
    life: 2000
  })
}
</script>

<template>
  <div class="flex gap-6">
    <div class="flex flex-col">
      <div class="flex flex-row gap-6 items-center">
        <span class="font-bold">{{ $t('copy_user_password.username') }}:</span>
        <span>{{ props.username }}</span>
      </div>
      <div class="flex flex-row gap-6 items-center h-10">
        <span class="font-bold">{{ $t('copy_user_password.password') }}:</span>
        <span>
          <Inplace
            :closable="true"
            :pt="{
              root: '!p-0',
              content: '!p-0',
              display: '!p-0',
              closeButton: '!p-0'
            }"
            :ptOptions="{ mergeSections: true, mergeProps: true }"
            :closeButtonProps="{ text: true, size: 'small', style: '!p-0' }"
          >
            <template #display>
              <span style="vertical-align: middle">*****</span>
              <span class="pi pi-eye" style="margin-left: 0.5rem; vertical-align: middle"></span>
            </template>

            <template #content>
              <span style="vertical-align: middle">{{ props.password }}</span>
            </template>
            <template #closeicon><i class="pi pi-eye-slash" style=""></i> </template> </Inplace
        ></span>
      </div>
    </div>

    <div v-if="props.status !== 'done'" class="">
      <Button
        type="button"
        :label="$t('copy_user_password.button_label')"
        icon="pi pi-clipboard"
        @click="(e) => copyCredentialToClipBoard(e, props.username, props.password)"
      />
    </div>
  </div>
</template>

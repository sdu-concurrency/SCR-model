<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Select from 'primevue/select'
import Button from 'primevue/button'
import Drawer from 'primevue/drawer'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

const { t } = useI18n()
const breakpoints = useBreakpoints(breakpointsTailwind)
const isMobile = breakpoints.smallerOrEqual('md')

const mobileMenuVisible = ref(false)

const getFlagEmoji = (countryCode) => {
  return String.fromCodePoint(
    ...[...countryCode.toUpperCase()].map((x) => 0x1f1a5 + x.charCodeAt())
  )
}

const route = useRoute()

const adminPageName = 'admin-page'
const isAdminLogin = computed(() => {
  return route.name === adminPageName
})
const langs = ref([
  {
    name: 'da',
    flag: getFlagEmoji('DK')
  },
  {
    name: 'en',
    flag: getFlagEmoji('GB')
  }
])
</script>

<template>
  <Menubar
    class="!bg-medium-gray"
    :pt="{
      root: '!rounded-none !border-0 !border-b-8 !border-project-red'
    }"
  >
    <template #start>
      <div class="flex align-items-center gap-2 h-14">
        <img class="object-contain h-full md:h-auto max-h-14" src="/scr smv logo.jpg" />
      </div>
    </template>

    <template #end>
      <!-- Desktop menu -->
      <div v-if="!isMobile" class="flex gap-4 mr-4">
        <UserIcon v-if="!isAdminLogin"></UserIcon>
        <Select
          v-model="$i18n.locale"
          :options="langs"
          optionLabel="name"
          optionValue="name"
          :pt="{
            label: '!p-2'
          }"
          :ptOptions="{ mergeSections: true, mergeProps: true }"
        >
          <template #value="slotProps">
            <div v-if="slotProps.value" class="flex items-center">
              {{ langs.find((e) => e.name === slotProps.value).flag }}
              {{ langs.find((e) => e.name === slotProps.value).name }}
            </div>
          </template>
          <template #option="slotProps">
            <div class="flex items-center">
              {{ slotProps.option.flag }} {{ slotProps.option.name }}
            </div>
          </template>
        </Select>
      </div>

      <!-- Mobile menu button -->
      <div v-else class="flex gap-2 mr-2">
        <UserIcon v-if="!isAdminLogin"></UserIcon>
        <Button
          icon="pi pi-bars"
          @click="mobileMenuVisible = true"
          text
          severity="secondary"
          class="!text-white"
        />
      </div>
    </template>
  </Menubar>

  <!-- Mobile drawer -->
  <Drawer v-model:visible="mobileMenuVisible" position="right" class="w-64">
    <template #header>
      <h3 class="font-bold">{{ $t('header.menu') || 'Menu' }}</h3>
    </template>
    <div class="flex flex-col gap-4 p-4">
      <div>
        <label class="block mb-2 font-semibold">{{ $t('header.language') || 'Language' }}</label>
        <Select
          v-model="$i18n.locale"
          :options="langs"
          optionLabel="name"
          optionValue="name"
          class="w-full"
          :pt="{
            label: '!p-2'
          }"
          :ptOptions="{ mergeSections: true, mergeProps: true }"
        >
          <template #value="slotProps">
            <div v-if="slotProps.value" class="flex items-center gap-2">
              {{ langs.find((e) => e.name === slotProps.value).flag }}
              {{ langs.find((e) => e.name === slotProps.value).name }}
            </div>
          </template>
          <template #option="slotProps">
            <div class="flex items-center gap-2">
              {{ slotProps.option.flag }} {{ slotProps.option.name }}
            </div>
          </template>
        </Select>
      </div>
    </div>
  </Drawer>
</template>

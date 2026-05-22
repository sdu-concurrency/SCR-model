<script setup>
import Timeline from 'primevue/timeline'

import { useI18n } from 'vue-i18n'
import { ref, computed } from 'vue'
const { t } = useI18n()

const props = defineProps(['currentPhase'])

const phases = ref([
  { phase: 1, description_key: 'phase_display.phase_1' },
  { phase: 2, description_key: 'phase_display.phase_2' },
  { phase: 3, description_key: 'phase_display.phase_3' },
  { phase: 4, description_key: 'phase_display.phase_4' }
])

const highlightPhase = computed(() => {
  if (props.currentPhase === '2') {
    return 2
  } else if (props.currentPhase === '3') {
    return 3
  }
  return 4
})
</script>
<template>
  <Timeline
    :value="phases"
    layout="horizontal"
    align="top"
    :pt="{
      eventSeparator: ({ props }) => ({
        class: [
          'items-center justify-center'
          //   'flex items-center flex-initial',
          //   {
          //     'flex-col': props.layout === 'vertical',
          //     'flex-row': props.layout === 'horizontal'
          //   }
        ]
      }),
      event: ({ props }) => ({
        class: '!flex-1'
      })
    }"
  >
    <template #marker="slotProps">
      <!-- <span
        class="flex w-8 h-8 items-center justify-center text-white rounded-full z-10 shadow-sm"
        :style="{ backgroundColor: slotProps.item.phase === highlightPhase ? '#8e1315' : 'gray' }"
      > -->
      <!-- </span> -->
      <div
        class="rounded-full z-10 shadow-sm text-light-gray self-stretch"
        :style="{ backgroundColor: slotProps.item.phase === highlightPhase ? '#8e1315' : 'gray' }"
      >
        <i class="pi pi-play p-3 pl-4" style="font-size: 2.5rem"></i>
      </div>
    </template>
    <template #connector=""> <span></span></template>
    <template #opposite="slotProps">
      {{ $t('session_view.phase', { current: slotProps.item.phase }) }}
    </template>
    <template #content="slotProps">
      <small class="text-surface-500 dark:text-surface-400">{{
        $t(slotProps.item.description_key)
      }}</small>
    </template>
  </Timeline>
</template>

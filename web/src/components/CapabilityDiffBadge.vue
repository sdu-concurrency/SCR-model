<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const props = defineProps(['diff'])

const severity = computed(() => {
  if (props.diff >= 3) {
    return 'danger'
  } else if (props.diff >= 1) {
    return 'warn'
  } else {
    return 'success'
  }
})
const tooltip = computed(() => {
  const base = t('analysis_3.table_columns.difference')
  if (props.diff >= 3) {
    return `Red: ${base} greater than 2`
  } else if (props.diff >= 1) {
    return `Orange: ${base} between 2 and 1`
  } else {
    return `Green: ${base} is 0`
  }
})
</script>

<template>
  <Badge
    :value="props.diff + ''"
    v-tooltip="{
      value: tooltip,
      showDelay: 1000,
      hideDelay: 300
    }"
    :severity="severity"
  ></Badge>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const props = defineProps(['risk_level'])

const severity = computed(() => {
  if (props.risk_level > 16) {
    return 'danger'
  } else if (props.risk_level > 8) {
    return 'warn'
  } else {
    return 'success'
  }
})

const tooltip = computed(() => {
  const base = t('analysis_3.table_columns.risk_level')
  if (props.risk_level > 16) {
    return `Red: ${base} greater than 16`
  } else if (props.risk_level > 8) {
    return `Orange: ${base} between 16 and 8`
  } else {
    return `Green: ${base} less than 8`
  }
})
</script>

<template>
  <Badge
    :value="props.risk_level + ''"
    :severity="severity"
    v-tooltip="{
      value: tooltip,
      showDelay: 1000,
      hideDelay: 300
    }"
  ></Badge>
</template>

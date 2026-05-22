<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: String
})

const severity = computed(() => {
  switch (props.status) {
    case 'ready':
      return 'info'
    case 'unfinished':
      return 'warn'
    default:
      return 'success'
  }
})

const status_key = computed(() => {
  return 'form_status.' + props.status
})
</script>

<template>
  <Badge :value="$t(status_key)" :severity="severity"></Badge>
</template>

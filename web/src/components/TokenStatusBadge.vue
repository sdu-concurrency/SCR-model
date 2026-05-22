<script setup>
import { computed } from 'vue'

const props = defineProps({
  is_valid: {
    type: Boolean,
    default: function () {
      return false
    }
  }
})

const severity = computed(() => {
  if (props.is_valid) {
    return 'success'
  } else {
    return 'warn'
  }
})

const status_key = computed(() => {
  if (props.is_valid) {
    return 'Valid'
  } else {
    return 'Expired'
  }
})
</script>

<template>
  <Badge :value="status_key" :severity="severity"></Badge>
</template>

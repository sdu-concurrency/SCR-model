<script setup lang="ts">
import Analysis2ViewNew from '@/views/Analysis2ViewNew.vue'

const props = defineProps({
  sessionName: String
})
</script>

<template>
  <Analysis2ViewNew :sessionName="props.sessionName"></Analysis2ViewNew>
</template>

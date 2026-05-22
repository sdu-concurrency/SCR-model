<script setup lang="ts">
import Analysis3ViewNew from '@/views/Analysis3ViewNew.vue'

const props = defineProps({
  sessionName: String
})
</script>

<template>
  <Analysis3ViewNew :sessionName="props.sessionName"></Analysis3ViewNew>
</template>

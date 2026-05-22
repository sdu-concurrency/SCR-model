<script setup lang="ts">
import { computed, ref, onMounted, inject } from 'vue'
import { parentSymbol } from '@formkit/vue'
import FormResultSummaryNew from './FormResultSummaryNew.vue'

const l = ref('en')
const value = ref({})
let formConfig: any
onMounted(() => {
  formConfig = inject(parentSymbol, null)
  if (formConfig) {
    l.value = formConfig?.config.locale
    if (formConfig?.parent!.value) {
      value.value = { form: formConfig?.parent!.value }
    } else {
      value.value = props['data-prop'] as Object
    }
  }
})

const isActive = computed(() => {
  if (!formConfig) {
    return true
  } else {
    return formConfig?.props.activeStep == 'summary'
  }
})

const props = defineProps({
  'data-prop': Object,
  'hide-job-function': Boolean,
  'is-show-summary': Boolean
})

const beforePrint = () => summary.value?.beforePrint()
const summary = ref<{ beforePrint: () => void } | null>(null)
defineExpose({
  beforePrint
})
</script>
<template>
  <FormResultSummaryNew
    ref="summary"
    v-if="isActive"
    :data="value"
    :hide-job-function="props['hide-job-function']"
    :is-show-summary="props['is-show-summary']"
  ></FormResultSummaryNew>
</template>

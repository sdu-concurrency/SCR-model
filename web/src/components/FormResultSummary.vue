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
      value.value = props.dataProp as Object
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
  dataProp: Object,
  hideJobFunction: Boolean,
  isShowSummary: Boolean,
  jobFunctionSchema: Array,
  vulnerabilitySchema: Array,
  capabilitySchema: Array
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
    :hide-job-function="props.hideJobFunction"
    :is-show-summary="props.isShowSummary"
    :job-function-schema="props.jobFunctionSchema"
    :vulnerability-schema="props.vulnerabilitySchema"
    :capability-schema="props.capabilitySchema"
  ></FormResultSummaryNew>
</template>

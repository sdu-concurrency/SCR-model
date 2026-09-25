<script setup>
import { computed, ref, inject, onMounted } from 'vue'
import { parentSymbol, useFormKitNodeById } from '@formkit/vue'
import { useI18n } from 'vue-i18n'
import { normalizeQuestionSchema, resolveLabel } from '@/composables/useQuestionLabels'
import { config } from '@/config'
const { t } = useI18n()
const now = ref(new Date().toISOString().split('T')[0])
const isOther = ref(false)
const l = ref('en')
onMounted(() => {
  const formConfig = inject(parentSymbol)
  l.value = formConfig.config.locale

  useFormKitNodeById('job_function', async (node) => {})
})

const props = defineProps({
  jobFunctionSchema: {
    type: Array,
    default: null
  }
})

const isOtherValue = (value) => value === 'Andet' || value === 'Other'

// Job function options resolved from the question schema (dynamic).
// Falls back to the static i18n options for legacy question sets without a schema.
const jobFunctionOptions = computed(() => {
  const options = {}
  for (const category of normalizeQuestionSchema(props.jobFunctionSchema)) {
    for (const item of category.items ?? []) {
      options[item.value] = resolveLabel(item.label, l.value || 'en', item.value)
    }
  }
  if (Object.keys(options).length > 0) {
    return options
  }
  return {
    Indkøb: t('form.new.step_1.select_jobfunction_options.Indkøb', l),
    Produktion: t('form.new.step_1.select_jobfunction_options.Produktion', l),
    Salg: t('form.new.step_1.select_jobfunction_options.Salg', l),
    Produktudvikling: t('form.new.step_1.select_jobfunction_options.Produktudvikling', l),
    'Økonomi / IT': t('form.new.step_1.select_jobfunction_options.Økonomi / IT', l),
    Andet: t('form.new.step_1.select_jobfunction_options.Andet', l)
  }
})
</script>

<template>
  <div class="space-y-4 sm:space-y-5">
    <FormKit
      type="select"
      name="job_function"
      id="job_function"
      :label="t('form.new.step_1.select_jobfunction_label', l)"
      :placeholder="t('form.new.step_1.select_jobfunction_placeholder', l)"
      validation="required"
      :options="jobFunctionOptions"
      @node="
        (node) => {
          node.on('commit', ({ payload }) => {
            isOther = isOtherValue(payload) ? true : false
          })
        }
      "
    ></FormKit>
    <FormKit
      v-if="isOther"
      type="text"
      id="job_function_other"
      name="job_function_other"
      :label="t('form.new.step_1.select_jobfunctionother_label', l)"
      :placeholder="t('form.new.step_1.select_other_placeholder', l)"
      validation="required"
    ></FormKit>
    <FormKit
      v-if="config.companyName"
      type="hidden"
      id="company"
      name="company"
      :value="config.companyName"
    ></FormKit>
    <FormKit
      v-else
      type="text"
      id="company"
      name="company"
      :label="t('form.new.step_1.select_company_label', l)"
      :placeholder="t('form.new.step_1.select_company_placeholder', l)"
      validation="required"
      autocomplete="false"
    ></FormKit>
    <FormKit type="hidden" id="dato" name="dato" :value="now"></FormKit>
  </div>
</template>

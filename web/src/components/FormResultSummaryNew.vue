<script setup>
import { computed, toRef, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useThrottleFn, useMediaQuery } from '@vueuse/core'
const print = useMediaQuery('print')

const { t, locale } = useI18n()

function resolveLabel(labelVal) {
  if (!labelVal) return ''
  if (typeof labelVal === 'object') return labelVal[locale.value] ?? labelVal['en'] ?? ''
  return labelVal
}

import { useResponseMapper } from '../composables/responseMapper'

const props = defineProps({
  data: Object,
  'hide-job-function': Boolean,
  'is-show-summary': Boolean
})
const date = computed(() => {
  if (props.isShowSummary === false) {
    return ''
  }
  if (!props.data || !props.data.form || !props.data.form.step_1) {
    return ''
  }
  return props.data.form.step_1.dato ? props.data.form.step_1.dato : ''
})
const job_function = computed(() => {
  if (
    !props.data ||
    !props.data.form ||
    !props.data.form.step_1 ||
    !props.data.form.step_1.job_function
  ) {
    return ''
  }
  if (props.data.form.step_1.job_function === 'Andet') {
    return props.data.form.step_1.job_function_other
  } else {
    return useResponseMapper(t, 'V1', props.data.form.step_1.job_function)
  }
})

const virksomhed = computed(() => {
  if (props.isShowSummary === false) {
    return ''
  }
  if (!props.data || !props.data.form || !props.data.form.step_1) {
    return ''
  }
  if (props.data.form.step_1.company) {
    return props.data.form.step_1.company
  } else if (props.data.form.step_1.virksomhed) {
    return props.data.form.step_1.virksomhed
  } else {
    return ''
  }
})

const riskMapper = (risk_management) => {}

const items = computed(() => {
  if (props.isShowSummary === false) {
    return ''
  }
  const res = []
  if (!props.data) {
    return []
  }
  if (!props.data.form) {
    return []
  }
  if (props.data.form.step_1.job_function == '') {
    return []
  }
  if (!props.data.form.step_2.vulnerability) {
    return []
  }
  let key = 'step_'
  let index = 2
  while (props.data.form[key + index]) {
    const vulnerability = props.data.form[key + index]
    const vulnerability_label = `${vulnerability.vulnerability.value} - ${resolveLabel(vulnerability.vulnerability.label)}`
    const note_vulnerability = vulnerability.vulnerability.note
    index++
    const impact = vulnerability.risk_management.horizontal.index + 1
    const likelihood = 5 - vulnerability.risk_management.vertical.index
    const risk_level = impact * likelihood
    const capability_list = vulnerability.capability_list
    const capability_list_with_diff = capability_list.map((e) => {
      return {
        diff: Math.abs(e.current_ability - e.importance),
        ...e
      }
    })
    res.push({
      vulnerability: vulnerability_label,
      note_vulnerability,
      impact,
      likelihood,
      risk_level,
      capability_list_with_diff: capability_list_with_diff
    })
  }
  return res
})

const expandedRows = ref([])

const expandAll = () => {
  expandedRows.value = items.value.reduce((acc, p) => (acc[p.vulnerability] = true) && acc, {})
}
const collapseAll = () => {
  expandedRows.value = null
}

let hist = ref({})
const beforePrint = useThrottleFn(() => {
  hist.value = expandedRows.value
  expandAll()
}, 500)

defineExpose({
  beforePrint
})
window.addEventListener('beforeprint', beforePrint)
window.addEventListener('afterprint', () => {
  expandedRows.value = hist.value
})
</script>
<template>
  <div v-if="props.isShowSummary !== false" class="mb-4 pb-4 card !print:text-xs">
    <div class="flex">
      <div class="grow">
        <h3 v-show="!props.hideJobFunction" class="text-lg">
          <span class="font-bold">{{ $t('survey_summary.jobfunction') }}:</span>
          {{ job_function }}
        </h3>
        <h3 class="text-lg">
          <span class="font-bold">{{ $t('survey_summary.virksomhed') }}:</span>
          {{ virksomhed }}
        </h3>
      </div>
      <div>
        <time class="place-self-end text-lg"
          ><span class="font-bold">{{ $t('survey_summary.dato') }}:</span>
          {{ date ? $d(new Date(date)) : '' }}</time
        >
      </div>
    </div>

    <DataTable
      v-model:expandedRows="expandedRows"
      :value="items"
      dataKey="vulnerability"
      class="print:text-xs print:mt-0 print:overflow-visible"
      :pt="{
        tableContainer: 'print:!max-h-none',
        column: {
          headerCell: 'print:!p-0',
          bodycell: 'print:!p-0'
        }
      }"
      :ptOptions="{ mergeSections: true, mergeProps: true }"
    >
      <template #header v-if="!print">
        <div class="flex flex-wrap justify-end gap-2 print:hidden">
          <Button
            text
            icon="pi pi-plus"
            :label="$t('survey_summary.expand_all')"
            @click="expandAll"
          />
          <Button
            text
            icon="pi pi-minus"
            :label="$t('survey_summary.collapse_all')"
            @click="collapseAll"
          />
        </div>
      </template>
      <Column expander style="width: 5rem" />
      <Column
        field="vulnerability_label"
        :header="$t('survey_summary.table_columns.vulnerability')"
      >
        <template #body="slotProps">
          <div class="flex flex-col">
            <span>{{ slotProps.data.vulnerability }}</span>
            <span class="text-xs font-light">{{ slotProps.data.note_vulnerability }}</span>
          </div>
        </template>
      </Column>
      <Column field="impact" :header="$t('survey_summary.table_columns.impact')"></Column>
      <Column field="likelihood" :header="$t('survey_summary.table_columns.likelihood')"></Column>
      <Column field="risk_level" :header="$t('survey_summary.table_columns.risk_level')">
        <template #body="slotProps">
          <RiskLevelBadge :risk_level="slotProps.data.risk_level"></RiskLevelBadge> </template
      ></Column>
      <template #expansion="slotProps">
        <div class="bg-slate-300 print:m-0 -m-4">
          <DataTable
            v-if="slotProps.data.capability_list_with_diff.length > 0"
            :value="slotProps.data.capability_list_with_diff"
            :pt="{
              tableContainer: 'print:!max-h-none',
              column: {
                headerCell: 'print:!p-1 bg-slate-100',
                bodycell: 'print:!p-1'
              }
            }"
          >
            <Column field="capability" :header="$t('survey_summary.table_columns.capability')">
              <template #body="slotProps">
                <div class="flex flex-col">
                  <span>{{ slotProps.data.value }} - {{ resolveLabel(slotProps.data.label) }}</span>
                  <span class="text-xs font-light">{{ slotProps.data.note }}</span>
                </div>
              </template></Column
            >
            <Column
              field="importance"
              :header="$t('survey_summary.table_columns.importance')"
            ></Column>
            <Column
              field="current_ability"
              :header="$t('survey_summary.table_columns.current_ability')"
            ></Column>
            <Column field="diff" :header="$t('survey_summary.table_columns.difference')"
              ><template #body="slotProps">
                <CapabilityDiffBadge
                  v-if="Number(slotProps.data.diff) !== NaN"
                  :diff="slotProps.data.diff"
                ></CapabilityDiffBadge> </template
            ></Column>
          </DataTable>
        </div>
      </template>
    </DataTable>
  </div>
</template>

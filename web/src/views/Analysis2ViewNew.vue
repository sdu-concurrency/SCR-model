<script setup>
import { onMounted, reactive, ref, computed, nextTick } from 'vue'
import resolveConfig from 'tailwindcss/resolveConfig'
import tailwindConfig from '../../tailwind.config.ts'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { usePocketbaseStore } from '@/stores/pb'
import Chart from 'primevue/chart'
import Breadcrumb from 'primevue/breadcrumb'
import Badge from 'primevue/badge'
import Divider from 'primevue/divider'

import { Colors, Chart as ChartJS } from 'chart.js'
import ChartDataLabels from 'chartjs-plugin-datalabels'
ChartJS.register(ChartDataLabels)

import Accordion from 'primevue/accordion'
import AccordionPanel from 'primevue/accordionpanel'
import AccordionHeader from 'primevue/accordionheader'
import AccordionContent from 'primevue/accordioncontent'
import { useResponseMapper } from '../composables/responseMapper.ts'
import { useI18n } from 'vue-i18n'
import CapabilityStackedChart from '../components/CapabilityStackedChart.vue'
const { t } = useI18n()
const twConfig = resolveConfig(tailwindConfig)
const { api } = usePocketbaseStore()
const props = defineProps({
  sessionName: String
})
const state = reactive({ surveys: [] })

// Helper function to deduplicate job_function values
const getUniqueJobFunction = (jobFunction, occurrences) => {
  if (!occurrences[jobFunction]) {
    occurrences[jobFunction] = 0
  }
  occurrences[jobFunction]++

  if (occurrences[jobFunction] === 1) {
    return jobFunction
  } else {
    return `${jobFunction}_${occurrences[jobFunction] - 1}`
  }
}

state.surveys = await api.collection('surveys').getFullList({
  filter: `user.role = "normal" && session.name = "${props.sessionName}"`,
  expand: 'user,session'
})

const common_vulnerabilities = computed(() => {
  const res = []
  for (const survey of state.surveys) {
    const { job_function } = survey.response.form.step_1
    const answers = answersFromForm(survey.response)
    for (const { capability_list, vulnerability, risk_management } of answers) {
      const risk_level =
        (risk_management.horizontal.index + 1) * (risk_management.vertical.index + 1)

      const entry = res.find((e) => e.vulnerability === vulnerability.value)
      if (!entry) {
        res.push({
          count: 1,
          risk_level_avg: risk_level,
          impact_avg: Number(risk_management.horizontal.index + 1),
          likelihood_avg: Number(risk_management.vertical.index + 1),
          vulnerability: vulnerability.value,
          note_vulnerability: [vulnerability.note],
          risk_level_sum: risk_level,
          impact_sum: Number(risk_management.horizontal.index + 1),
          likelihood_sum: Number(risk_management.vertical.index + 1)
        })
      } else {
        entry.count++
        entry.note_vulnerability.push(vulnerability.note)
        entry.risk_level_sum += risk_level
        entry.impact_sum += Number(risk_management.horizontal.index + 1)
        entry.likelihood_sum += Number(risk_management.vertical.index + 1)
        entry.risk_level_avg = entry.risk_level_sum / entry.count
        entry.impact_avg = entry.impact_sum / entry.count
        entry.likelihood_avg = entry.likelihood_sum / entry.count
      }
    }
  }
  return res.sort((a, b) =>
    b.count - a.count == 0 ? b.risk_level_avg - a.risk_level_avg : b.count - a.count
  )
})

const common_vulnerabilities_count_chart_data = computed(() => {
  return {
    labels: common_vulnerabilities.value.map((e) => useResponseMapper(t, 'V1', e.vulnerability)),
    datasets: [
      {
        label: t('analysis_2.common_vulnerabilities_count_label'),
        backgroundColor: twConfig.theme.colors.cyan[500],
        borderColor: twConfig.theme.colors.cyan[500],
        data: common_vulnerabilities.value.map((e) => e.count),
        borderWidth: 1,
        datalabels: {
          anchor: 'center',
          align: 'center',
          color: '#fff',
          font: { weight: 'bold', size: 11 },
          formatter: (v) => v,
          display: (ctx) => ctx.dataset.data[ctx.dataIndex] !== 0
        }
      },
      {
        label: t('analysis_2.common_vulnerabilities_risk_avg_label'),
        backgroundColor: twConfig.theme.colors.red[500],
        borderColor: twConfig.theme.colors.red[500],
        data: common_vulnerabilities.value.map((e) => e.risk_level_avg.toFixed(1)),
        borderWidth: 1,
        datalabels: {
          anchor: 'center',
          align: 'center',
          color: '#fff',
          font: { weight: 'bold', size: 11 },
          formatter: (v) => v,
          display: (ctx) => ctx.dataset.data[ctx.dataIndex] !== 0
        }
      },
      {
        label: t('analysis_2.common_vulnerabilities_risk_sum_label'),
        backgroundColor: twConfig.theme.colors.yellow[500],
        borderColor: twConfig.theme.colors.yellow[500],
        data: common_vulnerabilities.value.map((e) => e.risk_level_sum),
        borderWidth: 1,
        datalabels: {
          anchor: 'center',
          align: 'center',
          color: '#fff',
          font: { weight: 'bold', size: 11 },
          formatter: (v) => v,
          display: (ctx) => ctx.dataset.data[ctx.dataIndex] !== 0
        }
      }
    ]
  }
})

const functional_vulnerabilities = computed(() => {
  const res = []
  const jobFunctionOccurrences = {}
  for (const survey of state.surveys) {
    let { job_function, job_function_other } = survey.response.form.step_1
    job_function = job_function === 'Andet' ? job_function_other : job_function

    // Apply deduplication to job_function
    const unique_job_function = getUniqueJobFunction(job_function, jobFunctionOccurrences)

    const answers = answersFromForm(survey.response)
    let priority = 1

    for (const { capability_list, vulnerability, risk_management } of answers) {
      const risk_level =
        (risk_management.horizontal.index + 1) * (risk_management.vertical.index + 1)

      res.push({
        job_function: unique_job_function,
        priority,
        vulnerability: vulnerability.value,
        note_vulnerability: vulnerability.note,
        risk_level: risk_level,
        capability_list: capability_list.map((e) => {
          return { diff: Math.abs(e.importance - e.current_ability), ...e }
        })
      })
      expand['accordion_' + unique_job_function + '_' + priority + '_' + vulnerability.value] = ''
      priority++
    }
  }
  return res
})

const answersFromForm = ({ form }) => {
  const res = []
  let index = 2
  while (form[`step_${index}`]) {
    res.push(form[`step_${index}`])
    index++
  }
  return res
}

const summary = computed(() => {
  const vulnerabilities = {}
  const jobFunctionOccurrences = {}
  for (const survey of state.surveys) {
    let { job_function, job_function_other } = survey.response.form.step_1
    job_function = job_function === 'Andet' ? job_function_other : job_function

    // Apply deduplication to job_function
    const unique_job_function = getUniqueJobFunction(job_function, jobFunctionOccurrences)

    const answers = answersFromForm(survey.response)
    for (const { capability_list, vulnerability, risk_management } of answers) {
      // for (const { current_ability, importance, value } of capability_list) {
      //   console.log(capability_list, vulnerability, risk_management, {
      //     current_ability,
      //     importance,
      //     value
      //   })
      // }
      const risk_level =
        (risk_management.horizontal.index + 1) * (risk_management.vertical.index + 1)

      if (vulnerabilities[vulnerability.value]) {
        vulnerabilities[vulnerability.value].push({
          job_function: unique_job_function,
          risk_level,
          note_vulnerability: vulnerability.note,
          capability_list: capability_list.map((e) => {
            return {
              ...e,
              diff: Math.abs(e.current_ability - e.importance)
            }
          })
        })
      } else {
        vulnerabilities[vulnerability.value] = [
          {
            job_function: unique_job_function,
            risk_level,
            note_vulnerability: vulnerability.note,
            capability_list: capability_list.map((e) => {
              return {
                ...e,
                diff: Math.abs(e.current_ability - e.importance)
              }
            })
          }
        ]
      }
    }
  }

  const res = []
  for (const [key, value] of Object.entries(vulnerabilities)) {
    const vulnerability = {
      vulnerability: key,
      risk_level_avg: null,
      job_function_list: value,
      capability_list: []
    }

    vulnerability.risk_level_avg = value.reduce((acc, c) => acc + c.risk_level, 0) / value.length

    // map capability
    for (let job_function of value) {
      for (let capability of job_function.capability_list) {
        const existed = vulnerability.capability_list.find((e) => e.capability === capability.value)
        if (existed) {
          existed.current_ability_sum += Number(capability.current_ability)
          existed.importance_sum += Number(capability.importance)
          existed.diff_sum += Number(capability.diff)
          existed.job_functions.push({
            job_function: job_function.job_function,
            current_ability: capability.current_ability,
            importance: capability.importance,
            note_capability: capability.note_capability,
            diff: capability.diff
          })
          existed.current_ability_avg = existed.current_ability_sum / existed.job_functions.length
          existed.importance_avg = existed.importance_sum / existed.job_functions.length
          existed.diff_avg = existed.diff_sum / existed.job_functions.length
        } else {
          vulnerability.capability_list.push({
            capability: capability.value,
            current_ability_avg: Number(capability.current_ability),
            importance_avg: Number(capability.importance),
            diff_avg: capability.diff,
            current_ability_sum: Number(capability.current_ability),
            importance_sum: Number(capability.importance),
            diff_sum: capability.diff,
            job_functions: [
              {
                job_function: job_function.job_function,
                current_ability: capability.current_ability,
                importance: capability.importance,
                note_capability: capability.note_capability,
                diff: capability.diff
              }
            ]
          })
        }
      }
    }
    vulnerability.capability_list = vulnerability.capability_list.sort(
      (a, b) => b.diff_avg - a.diff_avg
    )
    res.push(vulnerability)
  }

  return res.sort((a, b) => b.job_function_list.length - a.job_function_list.length)
})

const capability_chart_data = computed(() => {
  const cap = {}
  if (summary.value) {
    for (const vulnerability of summary.value) {
      for (const capability of vulnerability.capability_list) {
        if (cap[capability.capability]) {
          cap[capability.capability].push(...capability.job_functions)
        } else {
          cap[capability.capability] = [...capability.job_functions]
        }
      }
    }
  }

  // Build sorted entries so labels and data stay in sync
  const entries = Object.entries(cap)
    .map(([key, funcs]) => ({
      label: useResponseMapper(t, 'V1', key),
      count: funcs.length,
      avgDiff: funcs.reduce((acc, c) => acc + c.diff, 0) / funcs.length
    }))
    .sort((a, b) => b.count - a.count)

  return {
    labels: entries.map((e) => e.label),
    datasets: [
      {
        label: t('analysis_2.capability_chart_count'),
        backgroundColor: twConfig.theme.colors.cyan[500] + 'CC',
        borderColor: twConfig.theme.colors.cyan[500],
        data: entries.map((e) => e.count),
        borderWidth: 1,
        datalabels: {
          anchor: 'center',
          align: 'center',
          color: '#fff',
          font: { weight: 'bold', size: 11 },
          formatter: (v) => v
        }
      },
      {
        label: t('analysis_2.capability_chart_sum'),
        backgroundColor: twConfig.theme.colors.red[500] + 'CC',
        borderColor: twConfig.theme.colors.red[500],
        data: entries.map((e) => Number(e.avgDiff.toFixed(1))),
        borderWidth: 1,
        datalabels: {
          anchor: 'center',
          align: 'center',
          color: '#fff',
          font: { weight: 'bold', size: 11 },
          formatter: (v) => v.toFixed(1),
          display: (ctx) => ctx.dataset.data[ctx.dataIndex] !== 0
        }
      }
    ]
  }
})

const vulnerabilities_chart_options = computed(() => ({
  indexAxis: 'y',
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      labels: {
        padding: 12,
        usePointStyle: true,
        pointStyle: 'rectRounded',
        font: { size: 11 }
      }
    },
    tooltip: {
      callbacks: {
        title: (items) => {
          // Show full label in tooltip even if it was truncated on the axis
          return items[0]?.label || ''
        }
      }
    }
  },
  scales: {
    y: {
      stacked: true,
      ticks: {
        autoSkip: false,
        font: { size: 10 },
        callback: function (value, index) {
          const label = this.getLabelForValue(index)
          const maxLen = 28
          return label && label.length > maxLen ? label.slice(0, maxLen) + '…' : label
        }
      },
      grid: { display: false }
    },
    x: {
      stacked: true,
      ticks: { precision: 0, font: { size: 10 } },
      beginAtZero: true,
      grid: { color: 'rgba(0,0,0,0.05)' }
    }
  }
}))

const capability_chart_options = computed(() => ({
  indexAxis: 'y',
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      labels: {
        padding: 16,
        usePointStyle: true,
        pointStyle: 'rectRounded',
        font: { size: 12 }
      }
    },
    tooltip: {
      callbacks: {
        label: (ctx) => {
          const ds = ctx.dataset
          return `${ds.label}: ${ctx.formattedValue}`
        }
      }
    }
  },
  scales: {
    y: {
      stacked: true,
      ticks: {
        autoSkip: false,
        font: { size: 11 }
      },
      grid: {
        display: false
      }
    },
    x: {
      stacked: true,
      ticks: {
        precision: 0,
        font: { size: 10 }
      },
      beginAtZero: true,
      grid: {
        color: 'rgba(0,0,0,0.05)'
      }
    }
  }
}))

const breadcrumbHome = ref({
  icon: 'pi pi-home',
  route: '/session/' + props.sessionName
})
const breadcrumbItems = ref([{ label: t('analysis_2.header') }])

const expand = reactive({})
const expand_mem = reactive({})
const expandedRows = reactive({})
const expandedRows_mem = reactive({})

const chart = ref()
const chart2 = ref()
async function handleBeforePrint() {
  // Dynamically size charts based on number of labels so none are skipped
  // Use wider width for landscape print layout
  const printWidth = 1100
  const chart1LabelCount = common_vulnerabilities_count_chart_data.value?.labels?.length || 10
  const chart1Height = Math.max(600, chart1LabelCount * 40)
  if (chart.value && chart.value.getChart()) {
    chart.value.getChart().resize(printWidth, chart1Height)
  }

  const chart2LabelCount = capability_chart_data.value?.labels?.length || 10
  // Cap height so the chart fits on one landscape page (~650px usable after title/description)
  const maxPrintChartHeight = 650
  const chart2Height = Math.min(maxPrintChartHeight, Math.max(400, chart2LabelCount * 40))
  if (chart2.value && chart2.value.getChart()) {
    chart2.value.getChart().resize(printWidth, chart2Height)
  }

  // Expand all accordions
  Object.keys(expand).forEach(function (key) {
    expand_mem[key] = expand[key]
    expand[key] = '0'
  })

  // Expand all rows in the main DataTable
  summary.value.forEach(function (item) {
    expandedRows_mem[item.vulnerability] = expandedRows[item.vulnerability]
    expandedRows[item.vulnerability] = true
  })

  await nextTick()
}
window.addEventListener('beforeprint', handleBeforePrint)

function handleAfterPrint() {
  if (chart.value && chart.value.getChart()) {
    chart.value.getChart().resize()
  }
  if (chart2.value && chart2.value.getChart()) {
    chart2.value.getChart().resize()
  }

  // Restore accordion states
  Object.keys(expand).forEach(function (key) {
    expand[key] = expand_mem[key]
    expand_mem[key] = undefined
  })

  // Restore DataTable expansion states
  Object.keys(expandedRows).forEach(function (key) {
    expandedRows[key] = expandedRows_mem[key]
    expandedRows_mem[key] = undefined
  })
}
window.addEventListener('afterprint', handleAfterPrint)

const printPage = async function () {
  await handleBeforePrint()
  window.print()
  handleAfterPrint()
}
</script>

<template>
  <Breadcrumb :home="breadcrumbHome" :model="breadcrumbItems">
    <template #item="{ item, props }">
      <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
        <a :href="href" v-bind="props.action" @click="navigate">
          <span :class="[item.icon, 'text-color']" />
          <span class="text-primary font-semibold">{{ item.label }}</span>
        </a>
      </router-link>
      <a v-else :href="item.url" :target="item.target" v-bind="props.action">
        <span class="text-surface-700 dark:text-surface-0/80">{{ item.label }}</span>
      </a>
    </template></Breadcrumb
  >
  <div class="">
    <h1 class="text-2xl font-bold mt-6">{{ $t('analysis_2.header') }}</h1>

    <h2 class="text-xl font-bold mt-4">
      {{ $t('analysis_2.session') }}: {{ $route.params.sessionName }}
    </h2>
  </div>

  <div class="print:hidden flex flex-row-reverse">
    <Button icon="pi pi-print" @click="printPage" severity="secondary" aria-label="Print" />
  </div>

  <div class="pagebreak"></div>

  <Card class="py-1 my-1 chart-page">
    <template #title>
      <h3 class="text-lg">{{ $t('analysis_2.common_vulnerabilities_label') }}</h3></template
    >
    <template #content>
      <div class="flex flex-col">
        <div class="vuln-chart-wrapper">
          <Chart
            ref="chart"
            type="bar"
            :data="common_vulnerabilities_count_chart_data"
            :options="vulnerabilities_chart_options"
            class="vuln-chart"
          />
        </div>

        <i18n-t v-if="chart" keypath="analysis_2.vulnerabilities_graph_description.content" tag="p">
          <template v-slot:first>
            <span :style="{ color: chart.data.datasets[0].borderColor }">{{
              $t('analysis_2.vulnerabilities_graph_description.first')
            }}</span>
          </template>
          <template v-slot:second>
            <span :style="{ color: chart.data.datasets[1].borderColor }">{{
              $t('analysis_2.vulnerabilities_graph_description.second')
            }}</span>
          </template>
          <template v-slot:third>
            <span :style="{ color: chart.data.datasets[2].borderColor }">{{
              $t('analysis_2.vulnerabilities_graph_description.third')
            }}</span>
          </template>
        </i18n-t>
      </div>
    </template>
  </Card>

  <Card class="py-1 my-1">
    <template #title> </template>
    <template #content>
      <DataTable
        class="mt-2 print:text-xs print:mt-0 print:overflow-visible"
        :value="common_vulnerabilities"
        sortField="risk_level_avg"
        :sortOrder="-1"
        scrollHeight="500px"
        :pt="{
          tableContainer: 'overflow-x-auto print:!max-h-none print:overflow-visible',
          column: {
            headerCell: 'print:!p-1',
            bodycell: 'print:!p-1'
          }
        }"
        :ptOptions="{ mergeSections: true, mergeProps: true }"
      >
        <Column field="vulnerability" :header="$t('analysis_2.vulnerability')">
          <template #body="slotProps">
            <div class="flex flex-col">
              <span>{{ useResponseMapper($t, 'V1', slotProps.data.vulnerability) }}</span>
              <span
                v-for="note in slotProps.data.note_vulnerability"
                :key="note"
                class="text-xs font-light"
                >{{ note }}</span
              >
            </div>
          </template></Column
        >
        <Column field="risk_level_sum" :header="$t('analysis_2.sum')" sortable> </Column>

        <Column field="count" :header="$t('analysis_2.count')" sortable>
          <template #body="slotProps">
            {{ slotProps.data.count }}
          </template>
        </Column>
        <Column field="risk_level_avg" :header="$t('analysis_2.risk_level_avg')" sortable>
          <template #body="slotProps">
            <RiskLevelBadge
              class="ml-1"
              :risk_level="slotProps.data.risk_level_avg.toFixed(1)"
            ></RiskLevelBadge>
          </template>
        </Column>
      </DataTable>
    </template>
  </Card>
  <div class="pagebreak"></div>

  <Card class="py-1 my-1 mx-auto">
    <template #title>
      <h3 class="text-lg font-bold">{{ $t('analysis_3.summary_header') }}</h3>
    </template>
    <template #content>
      <h5 class="text-lg mb-2">{{ $t('analysis_2.functional_vulnerabilities') }}</h5>
      <DataTable
        :value="functional_vulnerabilities"
        class="print:text-xs"
        rowGroupMode="rowspan"
        groupRowsBy="job_function"
        scrollHeight="1000px"
        tableStyle="min-width: 50rem"
        sortMode="single"
        size="small"
        :pt="{
          tableContainer: 'overflow-x-auto print:!max-h-none print:overflow-visible',
          header: {
            bodycell: 'print:!p-1'
          },
          column: {
            headerCell: '!p-1',
            bodycell: '!p-1'
          }
        }"
        :ptOptions="{ mergeSections: true, mergeProps: true }"
      >
        <Column field="job_function" sortable :header="$t('analysis_2.job_function')">
          <template #body="functionalSlotProps">
            {{ useResponseMapper($t, 'V1', functionalSlotProps.data.job_function) }}
          </template>
        </Column>
        <Column field="priority" sortable :header="$t('analysis_2.priority')"> </Column>
        <Column style="width: 70%" :header="$t('analysis_2.vulnerability')">
          <template #body="slotProps">
            <Accordion
              v-model:value="
                expand[
                  'accordion_' +
                    slotProps.data.job_function +
                    '_' +
                    slotProps.data.priority +
                    '_' +
                    slotProps.data.vulnerability
                ]
              "
              :pt="{
                root: 'print:w-4/5'
              }"
              :ptOptions="{ mergeProps: false }"
            >
              <AccordionPanel
                :pt="{
                  root: 'flex flex-col border-none [&>[data-pc-name=accordionheader]]:text-surface-600 dark:[&>[data-pc-name=accordionheader]]:text-surface-0 [&:last-child>[data-pc-name=accordioncontent]>[data-pc-section=content]]:rounded-b-md [&:nth-child(n+2)>[data-pc-name=accordionheader]]:border-t-0 [&:first-child>[data-pc-name=accordionheader]]:rounded-t-md'
                }"
                :ptOptions="{ mergeProps: false }"
                value="0"
              >
                <AccordionHeader
                  :pt="{
                    root: 'flex !p-1 bg-transparent border-0 items-start text-left align-center print:!p-0'
                  }"
                  :ptOptions="{ mergeProps: true }"
                >
                  <div class="flex flex-col">
                    <span>
                      {{ useResponseMapper($t, 'V1', slotProps.data.vulnerability) }}
                      <RiskLevelBadge
                        class="ml-1"
                        :risk_level="slotProps.data.risk_level"
                      ></RiskLevelBadge>
                    </span>
                    <span class="text-xs font-light">{{ slotProps.data.note_vulnerability }}</span>
                  </div>
                </AccordionHeader>
                <AccordionContent
                  :pt="{
                    content: '!p-0'
                  }"
                  :ptOptions="{ mergeProps: true }"
                >
                  <div class="overflow-x-auto">
                    <DataTable
                      rowGroupMode="subheader"
                      :value="slotProps.data.capability_list"
                      size="small"
                    >
                      <Column field="value" :header="$t('analysis_2.capability')">
                        <template #body="slotProps">
                          {{ useResponseMapper($t, 'V1', slotProps.data.value) }}
                        </template>
                      </Column>
                      <Column field="importance" :header="$t('analysis_2.importance')"></Column>
                      <Column
                        field="current_ability"
                        :header="$t('analysis_2.current_ability')"
                      ></Column>
                      <Column field="diff" :header="$t('analysis_3.table_columns.difference')"
                        ><template #body="slotProps">
                          <CapabilityDiffBadge
                            :diff="slotProps.data.diff"
                          ></CapabilityDiffBadge> </template
                      ></Column>
                    </DataTable>
                  </div>
                </AccordionContent>
              </AccordionPanel>
            </Accordion>
          </template>
        </Column>
        <template #groupheader="slotProps">
          <p class="text-lg">
            <span class="font-bold">{{ $t('analysis_2.job_function') }}</span
            >: {{ useResponseMapper($t, 'V1', slotProps.data.job_function) }}
          </p>
        </template>
      </DataTable>
      <Divider class="print:hidden" type="dashed" />
      <div class="pagebreak"></div>
      <h5 class="text-lg mb-2">{{ $t('analysis_2.vulnerabilities_summary_header') }}</h5>
      <DataTable
        v-model:expandedRows="expandedRows"
        :value="summary"
        class="print:text-xs"
        dataKey="vulnerability"
        scrollable
        scrollHeight="1000px"
        tableStyle="min-width: 50rem"
        size="small"
        :pt="{
          tableContainer: 'overflow-x-auto print:!max-h-none print:overflow-visible',
          column: {
            headerCell: 'print:!p-1',
            bodycell: 'print:!p-1'
          }
        }"
        :ptOptions="{ mergeSections: true, mergeProps: true }"
      >
        <Column expander style="width: 5rem" />

        <Column field="vulnerability" :header="$t('analysis_3.table_columns.vulnerability')">
          <template #body="slotProps">
            {{ useResponseMapper($t, 'V1', slotProps.data.vulnerability) }}
          </template>
        </Column>

        <Column
          :sortField="
            (e) => {
              return e.job_function_list.length
            }
          "
          :header="'Count'"
          sortable
        >
          <template #body="slotProps">
            {{ slotProps.data.job_function_list.length }}
          </template>
        </Column>

        <Column field="risk_level_avg" :header="'Average Risk Level'" sortable>
          <template #body="slotProps">
            <RiskLevelBadge
              class="ml-1"
              :risk_level="slotProps.data.risk_level_avg.toFixed(1)"
            ></RiskLevelBadge>
          </template>
        </Column>

        <template #expansion="slotProps">
          <div
            v-if="slotProps.data.job_function_list.filter((e) => e.note_vulnerability).length > 0"
            class="p-2"
          >
            <h5 class="font-bold">{{ $t('analysis_2.notes') }}:</h5>
            <div class="flex flex-col gap-1">
              <span
                v-for="job_function in slotProps.data.job_function_list.filter(
                  (e) => e.note_vulnerability
                )"
                :key="job_function"
                class="text-xs font-light"
              >
                <Badge severity="secondary" :value="job_function.job_function" rounded> </Badge>
                : {{ job_function.note_vulnerability }}</span
              >
            </div>
          </div>

          <div class="p-2">
            <h5 class="font-bold">{{ $t('analysis_2.capability_summary') }}:</h5>
            <div class="overflow-x-auto">
              <DataTable :value="slotProps.data.capability_list" size="small" class="print:text-xs">
                <Column field="capability" :header="$t('analysis_2.capability')">
                  <template #body="slotProps">
                    {{ useResponseMapper($t, 'V1', slotProps.data.capability) }}
                    <div class="flex flex-col gap-1">
                      <div
                        v-for="(job_function, index) in slotProps.data.job_functions"
                        :key="index"
                        class="mr-1"
                      >
                        <Badge severity="secondary" :value="job_function.job_function" rounded>
                        </Badge>

                        <span v-if="job_function.note_capability" class="ml-1 text-xs font-light"
                          >: {{ job_function.note_capability }}</span
                        >
                      </div>
                    </div>
                  </template>
                </Column>
                <Column
                  field="importance_avg"
                  :header="$t('survey_summary.table_columns.importance_avg')"
                ></Column>
                <Column
                  field="current_ability_avg"
                  :header="$t('survey_summary.table_columns.current_ability_avg')"
                ></Column>
                <Column field="diff_avg" :header="$t('survey_summary.table_columns.difference_avg')"
                  ><template #body="slotProps">
                    <CapabilityDiffBadge
                      :diff="slotProps.data.diff_avg.toFixed(1)"
                    ></CapabilityDiffBadge> </template
                ></Column>
              </DataTable>
            </div>
          </div>
        </template>
      </DataTable>
    </template>
  </Card>
  <Card class="py-1 my-1 mx-auto">
    <template #title>
      <h3 class="text-lg font-bold">{{ $t('analysis_2.capability_chart') }}</h3>
    </template>
    <template #content>
      <div class="flex flex-col chart-page">
        <div class="capability-chart-wrapper">
          <Chart
            ref="chart2"
            type="bar"
            :data="capability_chart_data"
            :options="capability_chart_options"
            class="capability-chart"
          />
        </div>

        <i18n-t v-if="chart2" keypath="analysis_2.capabilities_graph_description.content" tag="p">
          <template v-slot:first>
            <span :style="{ color: chart2.data.datasets[0].borderColor }">{{
              $t('analysis_2.capabilities_graph_description.first')
            }}</span>
          </template>
          <template v-slot:second>
            <span :style="{ color: chart2.data.datasets[1].borderColor }">{{
              $t('analysis_2.capabilities_graph_description.second')
            }}</span>
          </template>
        </i18n-t>
      </div>
      <Card class="py-1 my-1 mx-auto">
        <template #title>
          <h3 class="text-lg font-bold">
            {{ $t('capability_stacked_chart.category_distribution_title') }}
          </h3>
        </template>
        <template #content>
          <div class="flex flex-col">
            <CapabilityStackedChart
              :capabilities-data="summary"
              :session-name="props.sessionName"
            />
          </div>
        </template>
      </Card>
    </template>
  </Card>
</template>

<style lang="css">
@page {
  margin: 5mm;
}
@page chart-landscape {
  size: landscape;
  margin: 5mm;
}
</style>

<style lang="css" scoped>
.vuln-chart-wrapper {
  position: relative;
  width: 100%;
  min-height: 400px;
  height: v-bind(
    "Math.max(400, (common_vulnerabilities_count_chart_data?.labels?.length || 10) * 38) + 'px'"
  );
}

.vuln-chart {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.capability-chart-wrapper {
  position: relative;
  width: 100%;
  min-height: 500px;
  height: v-bind("Math.max(500, (capability_chart_data?.labels?.length || 10) * 38) + 'px'");
}

.capability-chart {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

@media print {
  .pagebreak {
    clear: both;
    page-break-after: always;
  }

  .chart-page {
    page: chart-landscape;
  }

  .vuln-chart-wrapper {
    min-height: 400px;
    height: 600px !important;
  }

  .capability-chart-wrapper {
    min-height: 400px;
    max-height: 650px;
    height: 650px !important;
  }
}
</style>

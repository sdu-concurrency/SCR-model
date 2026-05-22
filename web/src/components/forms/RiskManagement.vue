<script setup lang="ts">
import { ref, onMounted, computed, inject, warn, watch } from 'vue'
import { parentSymbol } from '@formkit/vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const l = ref('en')
onMounted(() => {
  const formConfig = inject(parentSymbol)
  l.value = formConfig!.config.locale
})
const props = defineProps({
  context: Object
})

watch(
  () => props.context?._value,
  () => {
    risk_number.value =
      props.context?._value.horizontal.index + 1 + props.context?._value.vertical.index * 5
  },
  { deep: true }
)

const risk_number = ref(0)
const rows = 5
const cols = 5

function row(no: number) {
  return Math.floor((no - 1) / 5)
}
function col(no: number) {
  return (no - 1) % 5
}

function severity(no: number) {
  return (5 - (row(no) + 1) + 1) * (col(no) + 1)
}
function severityClass(no: number) {
  if (severity(no) > 16) {
    return 'border-red-300 bg-red-100'
  } else if (severity(no) > 8) {
    return 'border-yellow-300 bg-yellow-100'
  } else {
    return 'border-green-300 bg-green-100'
  }
}
const horizontal_label = 'Impact'
const vertical_label = 'Likelihood'
const col_headers = ['1 - Negligible', '2 - Minor', '3 - Moderate', '4 - Significant', '5 - Severe']
const col_headers_short = ['1', '2', '3', '4', '5']
const row_headers = ['1 - Very Unlikely', '2 - Unlikely', '3 - Possible', '4 - Likely', '5 - Very Likely'].reverse()
const row_headers_short = ['1', '2', '3', '4', '5'].reverse()

function handleInput(e: any) {
  props.context!.node.input({
    horizontal: {
      header: horizontal_label,
      value: col_headers[col(e.target.value)],
      index: col(e.target.value)
    },
    vertical: {
      header: vertical_label,
      value: row_headers[row(e.target.value)],
      index: row(e.target.value)
    }
  })
}

const zincRadio = ref({
  width: '1.75rem',
  height: '1.75rem',
  checkedBackground: '{zinc.700}',
  checkedHoverBackground: '{zinc.800}',
  disabledBackground: '{form.field.disabled.background}',
  filledBackground: '{form.field.filled.background}',
  borderColor: '{form.field.border.color}',
  hoverBorderColor: '{form.field.hover.border.color}',
  focusBorderColor: '{form.field.focus.border.color}',
  checkedBorderColor: '{zinc.700}',
  checkedHoverBorderColor: '{zinc.800}',
  checkedFocusBorderColor: '{zinc.700}',
  checkedDisabledBorderColor: '{form.field.border.color}',
  invalidBorderColor: '{form.field.invalid.border.color}',
  icon: {
    size: '0rem'
  }
})
</script>

<template>
  <div class="flex flex-col font-light text-xs sm:text-sm gap-4">
    <p>
      {{ t('form.new.step_2.risk_management_sub_label', l) }}
    </p>
    <div class="flex flex-col md:flex-row gap-4 md:gap-x-4">
      <div class="flex-1">
        <span class="text-neutral-700 text-sm font-bold mb-1 dark:text-neutral-300">{{
          t('form.new.step_2.risk_management_likelihood', l)
        }}</span
        >:
        <ul class="list-disc pl-4 space-y-1">
          <li>
            {{ t('form.new.step_2.risk_management_likelihood_help_list.1', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_likelihood_help_list.2', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_likelihood_help_list.3', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_likelihood_help_list.4', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_likelihood_help_list.5', l) }}
          </li>
        </ul>
      </div>
      <div class="flex-1">
        <span class="text-neutral-700 text-sm font-bold mb-1 dark:text-neutral-300">{{
          t('form.new.step_2.risk_management_impact', l)
        }}</span
        >:
        <ul class="list-disc pl-4 space-y-1">
          <li>
            {{ t('form.new.step_2.risk_management_impact_help_list.1', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_impact_help_list.2', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_impact_help_list.3', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_impact_help_list.4', l) }}
          </li>
          <li>
            {{ t('form.new.step_2.risk_management_impact_help_list.5', l) }}
          </li>
        </ul>
      </div>
    </div>

    <div
      class="grid grid-cols-7 place-content-center gap-1 sm:gap-2 md:gap-4 lg:gap-8 text-[0.625rem] sm:text-xs md:text-sm lg:text-base overflow-x-auto"
    >
      <div class="col-start-2 col-span-5 place-self-center font-semibold text-center">
        {{ horizontal_label }}
      </div>

      <div
        class="col-start-3 grid grid-cols-subgrid place-content-center col-span-5 gap-1 sm:gap-2 md:gap-4 lg:gap-8"
      >
        <div
          v-for="(header, i) in col_headers"
          :key="header"
          class="place-self-center text-center leading-tight break-words max-w-[3rem] sm:max-w-none"
        >
          <span class="hidden md:inline">{{ header }}</span>
          <span class="md:hidden">{{ col_headers_short[i] }}</span>
        </div>
      </div>
      <div
        class="row-start-3 row-span-5 place-self-center font-semibold [writing-mode:vertical-lr] lg:[writing-mode:horizontal-tb] rotate-180 lg:rotate-0 text-center"
      >
        {{ vertical_label }}
      </div>

      <div
        class="row-start-3 grid grid-row-subgrid row-span-5 gap-1 sm:gap-2 md:gap-4 lg:gap-8 text-center lg:text-left"
      >
        <div
          v-for="(header, i) in row_headers"
          :key="header"
          class="leading-tight break-words flex items-center justify-center lg:justify-start"
        >
          <span class="hidden md:inline">{{ header }}</span>
          <span class="md:hidden">{{ row_headers_short[i] }}</span>
        </div>
      </div>
      <div
        v-for="i in rows * cols"
        :key="i"
        class="flex items-center justify-center rounded-lg p-0.5 sm:p-1 md:p-2 border-2 sm:border-3 md:border-4 min-w-[1.75rem] min-h-[1.75rem] sm:min-w-[2.5rem] sm:min-h-[2.5rem] md:min-w-[3rem] md:min-h-[3rem]"
        :class="severityClass(i)"
      >
        <RadioButton
          v-model="risk_number"
          :inputId="'risk-' + i"
          name="risk"
          :value="i"
          @input="handleInput"
          :dt="zincRadio"
        ></RadioButton>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs sm:text-sm">
      <span class="font-semibold">{{ t('form.new.step_2.risk_level_formula', l) }}</span>
      <span class="flex items-center gap-1.5">
        <span class="inline-block w-3 h-3 rounded-sm border-2 border-green-300 bg-green-100"></span>
        {{ t('form.new.step_2.risk_level_low', l) }}
      </span>
      <span class="flex items-center gap-1.5">
        <span class="inline-block w-3 h-3 rounded-sm border-2 border-yellow-300 bg-yellow-100"></span>
        {{ t('form.new.step_2.risk_level_medium', l) }}
      </span>
      <span class="flex items-center gap-1.5">
        <span class="inline-block w-3 h-3 rounded-sm border-2 border-red-300 bg-red-100"></span>
        {{ t('form.new.step_2.risk_level_high', l) }}
      </span>
    </div>
  </div>
</template>

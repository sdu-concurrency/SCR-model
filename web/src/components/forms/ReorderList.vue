<script setup lang="ts">
// @ts-nocheck
import { ref, onMounted, inject, computed, reactive, type Ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { FormKitNode } from '@formkit/core'
import { parentSymbol, useFormKitNodeById } from '@formkit/vue'
import { useI18n } from 'vue-i18n'
import { object } from 'yup'
const { t } = useI18n()
const l = ref('en')
const toast = useToast()
const props = defineProps({
  context: Object
})
const form: Ref<FormKitNode<any> | undefined> = useFormKitNodeById('form')
let formConfig: any
onMounted(() => {
  formConfig = inject(parentSymbol)
  l.value = formConfig!.config.locale
})

const onSelectionChange = (e) => {
  if (e.length > 1) {
    e.splice(0, 1)
  }
}
</script>

<template>
  <OrderList
    v-if="formConfig?.props.isActiveStep"
    :modelValue="props.context._value"
    @update:selection="onSelectionChange"
    scrollHeight="30rem"
    @reorder="(e) => context.node.emit('reorder', e)"
  >
    <template #option="{ option, selected, index }">
      <div
        class="flex flex-col sm:flex-row flex-wrap p-2 sm:p-3 items-start sm:items-center gap-3 sm:gap-4 w-full"
      >
        <div class="flex-1 min-w-0 w-full sm:w-auto" v-if="option.vulnerability.value !== ''">
          <span class="font-medium text-sm sm:text-base block break-words"
            >{{ option.vulnerability.value }} : {{ option.vulnerability.label }}
          </span>
          <span
            :class="[
              'text-xs sm:text-sm block mt-1 break-words',
              { 'text-surface-500 dark:text-surface-400': !selected, 'text-inherit': selected }
            ]"
            >{{ option.vulnerability.note }}</span
          >
        </div>
        <div
          class="flex-1 min-w-0 w-full sm:w-auto space-y-1"
          v-if="option.vulnerability.value !== ''"
        >
          <span class="font-bold text-xs sm:text-sm block"
            >{{ option.risk_management.horizontal.header }}:
            {{ option.risk_management.horizontal.value }}</span
          >
          <span class="font-bold text-xs sm:text-sm block"
            >{{ option.risk_management.vertical.header }}:
            {{ option.risk_management.vertical.value }}</span
          >
        </div>
      </div>
    </template>
  </OrderList>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, inject } from 'vue'
import { useToast } from 'primevue/usetoast'

import { parentSymbol } from '@formkit/vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const l = ref('en')
const toast = useToast()
const props = defineProps({
  context: Object
})
onMounted(() => {
  const formConfig = inject(parentSymbol)
  l.value = formConfig!.config.locale
})

const rawCapabilities = props.context!.options

function resolveLabel(locale: string, labelVal: any): string {
  if (!labelVal) return ''
  if (typeof labelVal === 'object') return labelVal[locale] ?? labelVal['en'] ?? ''
  return labelVal
}

const localizedCapabilities = computed(() => {
  return rawCapabilities.map((group: any) => ({
    ...group,
    label: resolveLabel(l.value, group.label),
    items: group.items.map((item: any) => ({
      ...item,
      label: resolveLabel(l.value, item.label)
    }))
  }))
})

async function onListBoxChange(e: any) {
  if (e.value.length > 5) {
    toast.add({
      severity: 'error',
      summary: t('form.new.step_2.capability_toomany_error_summary', l.value),
      detail: t('form.new.step_2.capability_toomany_error_detail', l.value),
      group: 'br',
      life: 3000
    })
    await props.context!.node.input(props.context?._value.slice(0, 5))
    props.context!._value = props.context?._value.slice(0, 5)
    return
  }
  await props.context!.node.input(
    e.value.map((a: any) => {
      return {
        current_ability: 1,
        importance: 1,
        ...a
      }
    })
  )
}

async function removeIndex(index: number) {
  if (index > -1) {
    // only splice array when item is found
    await props.context!.node.input(props.context?._value.toSpliced(index, 1))
  }
}

function showCustom(capability_id: string) {
  return (
    rawCapabilities.find((e: any) =>
      e.items.find((y: any) => {
        if (y.value !== capability_id) return false
        const label = resolveLabel(l.value, y.label)
        return label === 'Andet' || label === 'Other'
      })
    ) != null
  )
}

const sliderDT = ref({
  track: {
    size: '20px',
    border: {
      radius: '20px'
    }
  },
  handle: {
    width: '32px',
    height: '32px'
  }
})
</script>

<template>
  <div class="space-y-4">
    <span class="font-light text-sm block">{{ t('form.new.step_2.capability_help', l) }}</span>
    <Listbox
      :modelValue="props.context?._value"
      :options="localizedCapabilities"
      optionLabel="label"
      optionGroupLabel="label"
      optionGroupChildren="items"
      class="w-full"
      listStyle="max-height:400px"
      multiple
      checkmark
      :id="props.context?.id"
      @change="onListBoxChange"
      data-key="value"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-if="props.context?._value" v-for="[i, c] of props.context?._value.entries()" :key="i">
        <Card class="h-full">
          <template #title>
            <div class="flex flex-col gap-3 sm:gap-4">
              <InputGroup class="flex-wrap sm:flex-nowrap">
                <InputGroupAddon class="text-xs break-all">{{ c.value }} </InputGroupAddon>
                <InputText :value="resolveLabel(l, c.label)" :disabled="true" class="min-w-0 flex-1" />
                <InputGroupAddon>
                  <Button
                    icon="pi pi-times"
                    severity="secondary"
                    size="small"
                    @click="removeIndex(i)"
                    aria-label="Remove capability"
                  />
                </InputGroupAddon>
              </InputGroup>
              <InputText
                v-model="c.note"
                :placeholder="t('form.new.step_1.select_other_placeholder', l)"
                class="w-full"
              />
              <Message
                :class="{
                  invisible: !(showCustom(c.value) && (!c.note || c.note == ''))
                }"
                severity="error"
                variant="simple"
                size="small"
                >{{ t('form.new.step_1.select_other_placeholder', l) }}</Message
              >
            </div>
          </template>
          <template #content>
            <div class="flex flex-col gap-4 sm:gap-6">
              <div class="space-y-2">
                <label class="block text-sm sm:text-base"
                  >{{ t('form.new.step_2.current_ability_label', l) }} ({{
                    c.current_ability
                  }})</label
                >
                <div class="flex flex-row items-center gap-2">
                  <Button
                    icon="pi pi-minus"
                    severity="danger"
                    variant="text"
                    rounded
                    size="small"
                    aria-label="reduce"
                    class="shrink-0"
                    @click="c.current_ability <= 1 ? null : c.current_ability--"
                  />
                  <Slider
                    :min="1"
                    :max="5"
                    v-model="c.current_ability"
                    :dt="sliderDT"
                    class="flex-1 min-w-0"
                  />
                  <Button
                    icon="pi pi-plus"
                    severity="danger"
                    variant="text"
                    rounded
                    size="small"
                    aria-label="increase"
                    class="shrink-0"
                    @click="c.current_ability >= 5 ? null : c.current_ability++"
                  />
                </div>
              </div>

              <div class="space-y-2">
                <label class="block text-sm sm:text-base"
                  >{{ t('form.new.step_2.importance_label', l) }} ({{ c.importance }})</label
                >
                <div class="flex flex-row items-center gap-2">
                  <Button
                    icon="pi pi-minus"
                    severity="danger"
                    variant="text"
                    rounded
                    size="small"
                    aria-label="reduce"
                    class="shrink-0"
                    @click="c.importance <= 1 ? null : c.importance--"
                  />
                  <Slider
                    :min="1"
                    :max="5"
                    v-model="c.importance"
                    :dt="sliderDT"
                    class="flex-1 min-w-0"
                  />
                  <Button
                    icon="pi pi-plus"
                    severity="danger"
                    variant="text"
                    rounded
                    size="small"
                    aria-label="increase"
                    class="shrink-0"
                    @click="c.importance >= 5 ? null : c.importance++"
                  />
                </div>
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

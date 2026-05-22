<script setup lang="ts">
import { ref, onMounted, computed, inject, type Ref } from 'vue'
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

const rawOptions = props.context!.options
const label: Ref<string> = ref(props.context?.label)

function resolveLabel(l: string, labelVal: any): string {
  if (!labelVal) return ''
  if (typeof labelVal === 'object') return labelVal[l] ?? labelVal['en'] ?? ''
  return labelVal
}

const localizedOptions = computed(() => {
  return rawOptions.map((group: any) => ({
    ...group,
    label: resolveLabel(l.value, group.label),
    items: group.items.map((item: any) => ({
      ...item,
      label: resolveLabel(l.value, item.label)
    }))
  }))
})

const displayLabel = computed(() => resolveLabel(l.value, props.context?._value?.label))

const showNote = computed(() => {
  return (
    props.context?._value &&
    (displayLabel.value.includes('Andet') || displayLabel.value.includes('Other'))
  )
})

async function handleInput(e: any) {
  await props.context?.node.input(e.value)
}

async function onNoteFilled(e: any) {
  await props.context?.node.input({ ...props.context?._value, note: e })
}
</script>

<template>
  <div class="space-y-3 sm:space-y-4">
    <Inplace>
      <template #display>
        <div class="w-full">
          <h2
            v-if="props.context?._value && props.context?._value.label"
            class="text-base sm:text-lg lg:text-xl font-medium"
          >
            {{ displayLabel }}
          </h2>
          <Button
            v-else
            :label="label"
            icon="pi pi-search"
            iconPos="right"
            severity="secondary"
            raised
            class="w-full sm:w-auto"
          />
        </div>
      </template>
      <template #content="{ closeCallback }">
        <div
          class="flex flex-col sm:flex-row place-items-stretch sm:place-items-center justify-center gap-2"
        >
          <Listbox
            v-model="props.context!._value"
            :options="localizedOptions"
            optionLabel="label"
            optionGroupLabel="label"
            optionGroupChildren="items"
            :id="props.context?.id"
            class="w-full sm:max-w-md"
            listStyle="max-height:400px"
            @change="
              (e: any) => {
                if (e.value == null) {
                  return false
                }
                handleInput(e)
                closeCallback()
              }
            "
          >
            <template #optiongroup="slotProps">
              <div class="flex items-center">
                <div>{{ slotProps.option.label }}</div>
              </div>
            </template>
          </Listbox>
          <Button
            icon="pi pi-times"
            text
            severity="danger"
            @click="closeCallback"
            class="w-full sm:w-auto"
            aria-label="Close"
          />
        </div>
      </template>
    </Inplace>

    <InputText
      label="Note"
      :placeholder="t('form.new.step_1.select_other_placeholder', l)"
      type="text"
      fluid
      :modelValue="props.context?._value?.note"
      @update:modelValue="onNoteFilled"
      class="w-full"
    />

    <Message
      v-if="props.context?._value?.note == '' && showNote"
      severity="error"
      variant="simple"
      size="small"
      >{{ t('form.new.step_1.select_other_placeholder', l) }}</Message
    >
  </div>
</template>

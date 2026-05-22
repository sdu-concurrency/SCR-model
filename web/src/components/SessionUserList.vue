<script setup>
import DataView from 'primevue/dataview'

const props = defineProps({
  data: Array,
  sessionName: String
})
</script>

<template>
  <DataView :value="props.data">
    <template #list="slotProps">
      <div class="flex flex-col">
        <div v-for="(item, index) in slotProps.items" :key="index" class="w-full">
          <div
            class="flex flex-col gap-4 p-4"
            :class="{ 'border-t border-surface-200 dark:border-surface-700': index !== 0 }"
          >
            <SessionUser :user="item" :sessionName="props.sessionName"></SessionUser>
          </div>
        </div>
      </div>
    </template>
  </DataView>
</template>

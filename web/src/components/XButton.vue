<script setup>
import { computed } from 'vue'
const props = defineProps({
  delete_target: String,
  delete_survey_index: Number,
  onDeleteQuestion: Function
})
const classAttr = computed(() => {
  const suffix = Number(props.delete_target.split('_').pop())
  return {
    hidden: suffix === 1
  }
})
</script>

<template>
  <button
    class="btn btn-error btn-circle float-right mt-10"
    :class="classAttr"
    @click.prevent="props.onDeleteQuestion(props.delete_survey_index, props.delete_target)"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  </button>
</template>

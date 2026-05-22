<script setup>
const props = defineProps({
  label: String,
  add_target: String,
  add_survey_step_index: Number,
  onAddQuestion: Function
})
</script>

<template>
  <Button
    type="button"
    class="py-3 px-6 text-sm font-semibold"
    :label="props.label"
    @click.prevent="props.onAddQuestion(props.add_survey_step_index, props.add_target)"
  ></Button>
</template>

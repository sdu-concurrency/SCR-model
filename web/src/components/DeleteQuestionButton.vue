<script setup>
const props = defineProps({
  label: String,
  delete_target: String,
  delete_survey_step_index: Number,
  onDeleteQuestion: Function
})
</script>

<template>
  <Button
    type="button"
    class="py-3 px-6 text-sm font-semibold ml-auto"
    severity="danger"
    :label="props.label"
    @click.prevent="props.onDeleteQuestion(props.delete_survey_step_index, props.delete_target)"
  ></Button>
</template>

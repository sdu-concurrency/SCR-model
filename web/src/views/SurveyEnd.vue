<script setup>
const props = defineProps({
  username: String,
  sessionName: String
})
</script>
<template>
  <main class="container mx-auto">
    <h1 class="text-xl md:text-2xl mb-4">{{ $t('survey_end.header') }}</h1>
    <p class="mb-4 text-sm md:text-base">{{ $t('survey_end.text') }}</p>
    <RouterLink
      :to="{
        name: 'session-form-view',
        params: { sessionName: props.sessionName, username: props.username }
      }"
      class="text-blue-600 visited:text-purple-600 text-sm md:text-base inline-block"
    >
      {{ $t('survey_end.link_to_result') }}
    </RouterLink>
  </main>
</template>

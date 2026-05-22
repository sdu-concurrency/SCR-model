<script setup lang="ts">
import { usePocketbaseStore } from '@/stores/pb'
import { storeToRefs } from 'pinia'
const { user } = storeToRefs(usePocketbaseStore())
const props = defineProps({
  sessionName: String
})
</script>
<template>
  <main class="container mx-auto">
    <h1 class="text-xl md:text-2xl mb-4">{{ $t('survey_start.header') }}</h1>
    <p class="mb-4 text-sm md:text-base">{{ $t('survey_start.text') }}</p>
    <RouterLink
      :to="{
        name: 'session-form',
        params: { sessionName: props.sessionName, username: user!.username }
      }"
      class="text-blue-600 visited:text-purple-600 text-sm md:text-base inline-block"
    >
      {{ $t('survey_start.link_to_survey') }}
    </RouterLink>
  </main>
</template>

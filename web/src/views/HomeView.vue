<script setup lang="ts">
import Button from 'primevue/button'
import Image from 'primevue/image'
import { onMounted, ref } from 'vue'
import { useAppStore } from '@/stores/app'
import { usePocketbaseStore } from '@/stores/pb'
import { storeToRefs } from 'pinia'
import { config } from '@/config'
const { isLoggedIn } = storeToRefs(usePocketbaseStore())
const { setUserDialogVisibility } = useAppStore()

// Company landing page (self-hosted builds): when VITE_LANDING_HTML points to
// an HTML fragment (e.g. /landing.html placed in web/public), it replaces the
// default landing content. Anchors with href="#open-login" open the login
// dialog; all other links (mailto, external) work natively.
const landingHtml = ref('')
onMounted(async () => {
  if (!config.landingHtmlPath) return
  try {
    const res = await fetch(config.landingHtmlPath)
    if (res.ok) {
      landingHtml.value = await res.text()
    }
  } catch {
    // fall back to the default landing page
  }
})

const onLandingClick = (e: MouseEvent) => {
  const anchor = (e.target as HTMLElement).closest?.('a')
  if (anchor && anchor.getAttribute('href') === '#open-login') {
    e.preventDefault()
    setUserDialogVisibility(true)
  }
}
</script>
<template>
  <!-- Company landing page (fetched HTML fragment) -->
  <main
    v-if="landingHtml"
    class="container mx-auto px-4 landing-html"
    @click="onLandingClick"
    v-html="landingHtml"
  ></main>

  <!-- Default landing page -->
  <main v-else class="container mx-auto px-4">
    <h1 class="text-xl md:text-2xl">Welcome to the Supply Chain Resilience Assessment Tool</h1>
    <p class="mt-3 text-sm md:text-base">
      {{ $t('home.desc') }}
    </p>
    <!-- <p class="mt-3">
      {{ $t('home.funding_desc') }}
    </p> -->
    <div class="flex flex-col md:flex-row gap-4 mt-5 items-center">
      <Image src="IndustriensFond_logo_BLACK_RGB.png" alt="Image" class="w-full max-w-[250px]" />
      <Image src="SDU_BLACK_RGB_png.png" alt="Image" class="w-full max-w-[250px]" />
    </div>

    <!-- <p class="text-sm">version 2</p> -->
  </main>

    <div class="mt-4 grid justify-items-center">
      <Button
        :label="isLoggedIn ? 'User' : 'Press here to login for session manager'"
        severity="secondary"
        @click.prevent="setUserDialogVisibility(true)"
        class="w-full md:w-auto text-sm md:text-base"
      />
    </div>
</template>

<style scoped>
/* Minimal prose styling for the injected company landing fragment (Tailwind
   classes cannot be used inside runtime-injected HTML). */
.landing-html :deep(h1) {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 1rem;
}
.landing-html :deep(h2) {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 1.5rem;
}
.landing-html :deep(p) {
  margin-top: 0.75rem;
  font-size: 0.9rem;
  line-height: 1.6;
}
.landing-html :deep(a) {
  color: rgb(37 99 235);
  text-decoration: underline;
}
</style>

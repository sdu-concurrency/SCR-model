<script setup lang="ts">
import Button from 'primevue/button'
import Image from 'primevue/image'
import { useAppStore } from '@/stores/app'
import { usePocketbaseStore } from '@/stores/pb'
import { storeToRefs } from 'pinia'
const { isLoggedIn } = storeToRefs(usePocketbaseStore())
const { setUserDialogVisibility } = useAppStore()
</script>
<template>
  <main class="container mx-auto px-4">
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

    <div class="mt-4 grid justify-items-center">
      <Button
        :label="isLoggedIn ? 'User' : 'Press here to login for session manager'"
        severity="secondary"
        @click.prevent="setUserDialogVisibility(true)"
        class="w-full md:w-auto text-sm md:text-base"
      />
    </div>

    <!-- <p class="text-sm">version 2</p> -->
  </main>
</template>

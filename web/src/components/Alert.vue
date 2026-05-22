<script setup></script>

<template>
  <span class="badge badge-flat-error">
    <slot />
  </span>
  <!-- <div role="alert" class="rounded border-s-4 border-green-500 bg-green-50 p-4">
    <strong v-if="header && header !== ''" class="block font-medium text-green-800">
      {{ header }}
    </strong>

    <p class="text-sm text-green-700">
      <slot />
    </p>
  </div> -->
</template>

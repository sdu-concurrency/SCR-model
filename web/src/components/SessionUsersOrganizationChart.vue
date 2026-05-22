<script setup>
// import OverlayPanel from 'primevue/overlaypanel'
import Popover from 'primevue/popover'
import OrganizationChart from 'primevue/organizationchart'
import { computed, reactive } from 'vue'

const props = defineProps({
  normalUsers: Array,
  superUser: Object,
  sessionName: String
})

const userTrees = computed(() => {
  return {
    ...props.superUser,
    styleClass: [
      ' !text-white rounded-xl',
      { '!bg-slate-500': props.superUser.status === 'ready' },
      { '!bg-slate-500': props.superUser.status === 'unfinished' },
      { '!bg-green-500': props.superUser.status === 'done' }
    ],
    children: props.normalUsers.map((e) => {
      return {
        ...e,
        styleClass: [
          ' !text-white rounded-xl',
          { '!bg-slate-500': e.status === 'ready' },
          { '!bg-slate-500': e.status === 'unfinished' },
          { '!bg-green-500': e.status === 'done' }
        ]
      }
    })
  }
})

const overlayRef = reactive({})

function setOverlayRef(key, el) {
  overlayRef[key] = el
}
</script>

<template>
  <div class="touch-pan-x overflow-x-auto">
    <OrganizationChart
      :value="userTrees"
      :pt="{
        node: '!p-2',
        table: 'min-w-full'
      }"
      :ptOptions="{ mergeSections: true, mergeProps: true }"
    >
      <template #default="slotProps">
        <p
          class="text-xs md:text-sm cursor-pointer"
          @click.prevent="
            (event) => {
              overlayRef[slotProps.node.username].toggle(event)
            }
          "
        >
          {{ slotProps.node.username }}
        </p>
        <Popover :ref="(el) => setOverlayRef(slotProps.node.username, el)">
          <SessionUser :user="slotProps.node" :sessionName="sessionName"></SessionUser>
        </Popover>
      </template>
    </OrganizationChart>
  </div>
</template>

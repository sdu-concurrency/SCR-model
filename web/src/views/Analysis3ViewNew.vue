<script setup>
import { reactive, onMounted, ref, nextTick } from 'vue'
import { usePocketbaseStore } from '@/stores/pb'
import Breadcrumb from 'primevue/breadcrumb'
import Card from 'primevue/card'
import FormResultSummaryNew from '@/components/FormResultSummaryNew.vue'

const { api } = usePocketbaseStore()
const props = defineProps({
  sessionName: String,
  username: String
})
const state = reactive({ survey: {} })

const fetchData = async () => {
  const data = await api.collection('surveys').getFullList({
    filter: `user.role = "super" && session.name = "${props.sessionName}"`,
    expand: 'user,session'
  })

  state.survey = data
}

const home = ref({
  icon: 'pi pi-home',
  route: '/session/' + props.sessionName
})
const items = ref([{ label: 'Analysis 3' }])

await fetchData()

const printPage = async function () {
  summary.value.beforePrint()
  await nextTick()
  window.print()
}
const summary = ref(null)
</script>

<template>
  <Breadcrumb :home="home" :model="items">
    <template #item="{ item, props }">
      <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
        <a :href="href" v-bind="props.action" @click="navigate">
          <span :class="[item.icon, 'text-color']" />
          <span class="text-primary font-semibold">{{ item.label }}</span>
        </a>
      </router-link>
      <span v-else class="text-surface-700 dark:text-surface-0/80">{{ item.label }}</span>
    </template>
  </Breadcrumb>
  <Card class="py-4">
    <template #title>
      <h1 class="text-2xl">{{ $t('survey_summary.header') }}</h1></template
    >
    <template #content>
      <div class="print:hidden flex flex-row-reverse">
        <Button icon="pi pi-print" @click="printPage" severity="secondary" aria-label="Print" />
      </div>
      <FormResultSummaryNew
        v-if="state.survey[0]"
        ref="summary"
        :data="state.survey[0].response"
        :hide-job-function="true"
        :is-show-summary="true"
      ></FormResultSummaryNew>
      <div v-else>
        <p>{{ $t('survey_summary.data_missing') }}</p>
      </div>
    </template>
  </Card>
</template>

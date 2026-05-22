<script setup>
import { getNode } from '@formkit/core'
import { getValidationMessages } from '@formkit/validation'

import { onMounted, ref, computed, reactive, watch, toRaw } from 'vue'
import { usePocketbaseStore } from '@/stores/pb'
import { useI18n } from 'vue-i18n'
import { changeLocale, useFormKitNodeById } from '@formkit/vue'
import FormResultSummary from '@/components/FormResultSummary.vue'
import { createLocalStoragePlugin } from '@formkit/addons'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import { ClientResponseError } from 'pocketbase'
import { useToast } from 'primevue/usetoast'

import { useRouter } from 'vue-router'
const router = useRouter()
const toast = useToast()

const breakpoints = useBreakpoints(breakpointsTailwind)
const hideProgressLabel = breakpoints.smallerOrEqual('lg')

const { api, user, session } = usePocketbaseStore()

const is_super = computed(() => {
  return user.role === 'super'
})

const { t } = useI18n()

const title = computed(() => {
  if (is_super.value) {
    return t('survey_title_3', { name: session.name })
  } else {
    return t('survey_title_2', { name: user.username })
  }
})

const props = defineProps({
  values: Object
})

const defaultVulnerabilityValue = {
  vulnerability: {
    value: '',
    label: '',
    note: ''
  },
  risk_management: {
    horizontal: {
      header: 'Severity',
      value: 'Negligible',
      index: 0
    },
    vertical: {
      header: 'Likelihood',
      value: 'Very Unlikely',
      index: 4
    }
  },
  capability_list: []
}
const values = ref(props.values)

const allowIncomplete = ref(false)
const vulnerability_number = ref(values.value?.form ? Object.keys(values.value.form).length - 4 : 1)
const l = ref('en')

async function gotoIncompleteStep() {
  return new Promise((resolve) => {
    useFormKitNodeById('multi-step', async (node) => {
      await node.settled
      // check step_1
      let targetStepIndex = 1
      if (node.value['step_1']) {
        const step1 = node.value['step_1']
        if (
          ((step1.job_function && step1.job_function !== '') ||
            (step1.job_function === 'Andet' && step1.job_function_other !== '')) &&
          step1.company &&
          step1.company !== ''
        ) {
          node.children[targetStepIndex - 1].props.valid = true
          node.children[targetStepIndex - 1].props.hasBeenVisited = true
          targetStepIndex = 2
        } else {
          resolve()
          return
        }
        while (node.value['step_' + targetStepIndex]) {
          const stepValue = node.value['step_' + targetStepIndex]
          if (
            stepValue.vulnerability &&
            stepValue.vulnerability.value &&
            stepValue.vulnerability.value !== '' &&
            stepValue.capability_list.length > 0
          ) {
            node.children[targetStepIndex - 1].props.valid = true
            node.children[targetStepIndex - 1].props.hasBeenVisited = true
            targetStepIndex++
          } else {
            break
          }
        }
        if (!node.value['step_' + targetStepIndex]) {
          node.props.activeStep = 'addmore'
        } else {
          node.props.activeStep = 'step_' + targetStepIndex
        }
      }
      resolve()
    })
  })
}

function appendIndex(str, joiner, index) {
  return str + joiner + index
}

const isDisabled = ref(false)

function increaseVulnerabilityNumber() {
  isDisabled.value = true
  vulnerability_number.value++
  const node = getNode('multi-step')
  node.settled.then(() => {
    node.goTo(`step_${vulnerability_number.value + 1}`)
    isDisabled.value = false
  })
}

async function removeElement(index) {
  const form = getNode('form')
  const multi = getNode('multi-step')
  const indexToRemove = index + 1
  const keyToGo = 'step_' + index
  form.settled.then(async () => {
    const res = { form: {} }
    const sortedKey = Object.keys(form.value.form).toSorted()
    for (const key of sortedKey) {
      if (key.startsWith('step')) {
        if (key.split('_')[1] < indexToRemove) {
          // lower, keep
          res.form[key] = form.value.form[key]
        } else {
          // greater, assign next
          // find next
          const nextStep = 'step_' + (Number(key.split('_')[1]) + 1)
          if (form.value.form[nextStep]) {
            res.form[key] = form.value.form[nextStep]
          }
        }
      } else {
        res.form[key] = form.value.form[key]
      }
    }
    await form.input(res)

    vulnerability_number.value--
    multi.goTo(keyToGo)
  })
}
const questions = ref({})
const isLoaded = ref(false)
const onlyVulnerabilities = ref([])

const ranking = (node) => {
  node.on('reorder', onReorder)
}

async function onReorder({ payload }) {
  let index = 2
  for (const vul of payload.value) {
    await new Promise((resolve) =>
      useFormKitNodeById(`step_${index}`, async (node) => {
        await node.input(vul)
        await node.settled
        resolve()
      })
    )
    index++
  }
}

onMounted(async () => {
  if (session && session.questions) {
    questions.value = await api.collection('questions').getOne(session.questions)
  } else {
    questions.value = await api.collection('questions').getFirstListItem()
  }
  changeLocale(questions.value.language)
  l.value = questions.value.language
  if (user.role === 'super') {
    values.value.form = {
      ...values.value.form,
      step_1: {
        job_function: 'Andet',
        company: session.name,
        job_function_other: 'Common'
      }
    }
  }
  isLoaded.value = true
  await gotoIncompleteStep()
  useFormKitNodeById('form', async (node) => {
    await node.settled

    let index = 2
    const ordered = []
    while (values.value.form[`step_${index}`]) {
      ordered.push(values.value.form[`step_${index}`])
      index++
    }
    onlyVulnerabilities.value = ordered

    node.on('commit', async ({ payload }) => {
      await node.settled

      let index = 2
      const ordered = []
      while (payload.form[`step_${index}`]) {
        ordered.push(payload.form[`step_${index}`])
        index++
      }
      onlyVulnerabilities.value = ordered
      values.value.form = payload.form
    })
  })
})

const removeNone = (formData) => {
  const f = structuredClone(toRaw(formData))
  delete f.form['addmore']
  delete f.form['summary']
  delete f.form['rank']
  return f
}

async function beforeLoad(value) {
  return null
}

const onSubmit = async (formData) => {
  const data = {
    session: session.id,
    user: user.id,
    is_complete: true,
    response: JSON.stringify(removeNone(formData))
  }

  let surveyId = null
  try {
    const res = await api
      .collection('surveys')
      .getFirstListItem(`user = "${user.id}" && session = "${session.id}"`)
    surveyId = res.id
  } catch (e) {
    const err = new ClientResponseError(e)
    if (err.status !== 404) {
      toast.add({
        severity: 'error',
        summary: 'error',
        detail: err.data.message + ' , ' + t('toast_contact_admin'),
        group: 'br',
        life: 2000
      })
    }
  }

  try {
    if (surveyId) {
      await api.collection('surveys').update(surveyId, data)
      toast.add({
        severity: 'success',
        summary: t('toast_survey_updated'),
        group: 'br',
        life: 2000
      })
    } else {
      await api.collection('surveys').create(data)
      await api.collection('app_users').update(user.id, { status: 'done' })
      toast.add({
        severity: 'success',
        summary: t('toast_survey_submitted'),
        group: 'br',
        life: 2000
      })
    }
    if (is_super.value) {
      await api.collection('sessions').update(session.id, { current_step: 'completed' })
    }
    await router.push({
      name: 'session-form-end',
      params: { sessionName: session.name, username: user.username },
      replace: true
    })
  } catch (e) {
    console.error(e)
    toast.add({
      severity: 'error',
      summary: 'error',
      detail: JSON.stringify(e) + ' , ' + t('toast_contact_admin'),
      group: 'br',
      life: 2000
    })
  }
}

const messages = ref([])

function showErrors(node) {
  const validations = getValidationMessages(node)
  messages.value = []
  validations.forEach((inputMessages) => {
    messages.value = messages.value.concat(inputMessages.map((message) => message.value))
  })
  const event = node.on('commit', () => {
    messages.value = []
    node.off(event)
  })
}
</script>

<template>
  <ProgressBar v-if="!isLoaded" />
  <h1 class="text-2xl font-bold">
    {{ title }}
  </h1>
  <FormKit
    v-if="isLoaded"
    type="form"
    name="f"
    id="form"
    key="f"
    :actions="false"
    v-model="values"
    :plugins="[
      createLocalStoragePlugin({
        key: 'n-u_' + user?.username,
        maxAge: 604800000, // 7 day
        debounce: 1000,
        clearOnSubmit: true
      })
    ]"
    use-local-storage
    @submit="onSubmit"
  >
    <FormKit
      type="multi-step"
      id="multi-step"
      tab-style="progress"
      name="form"
      key="ms"
      :allow-incomplete="allowIncomplete"
      :hide-progress-labels="hideProgressLabel"
    >
      <FormKit
        type="step"
        name="step_1"
        key="step_1"
        :label="t('form.new.step_1.step_label', l)"
        stepActionsClass="mt-4"
      >
        <JobFunction />
      </FormKit>

      <FormKit
        v-for="i in vulnerability_number"
        type="step"
        :id="appendIndex('step', '_', i + 1)"
        :name="appendIndex('step', '_', i + 1)"
        :key="appendIndex('step', '_', i + 1)"
        :label="appendIndex(t('form.new.step_2.step_label', l), ' ', i)"
        :value="{ ...defaultVulnerabilityValue }"
        stepActionsClass="mt-4"
      >
        <div v-if="i > 1">
          <Button
            class="float-right"
            label="Danger"
            severity="danger"
            rounded
            @click="removeElement(i)"
            >{{ t('form.new.step_2.remove_button', l) }}</Button
          >
        </div>
        <VulnerabilityStep
          :index="i"
          :vulnerabilityOptions="questions.vulnerability_schema"
          :capabilityOptions="questions.capability_schema"
        />

        <Message severity="error" v-if="messages.length > 0">
          <ul>
            <li v-for="message in messages" :key="message">
              {{ message }}
            </li>
          </ul>
        </Message>

        <template #stepNext="{ handlers, node }">
          <FormKit
            type="button"
            @click="
              () => {
                showErrors(node)
                handlers.incrementStep(1)()
              }
            "
            label="Next"
            data-next="true"
          />
        </template>
      </FormKit>

      <FormKit
        type="step"
        :label="t('form.new.step_addmore.step_label')"
        name="addmore"
        stepInner-class="m-4 p-4"
      >
        <div class="flex flex-col">
          <div class="flex justify-center items-center font-medium">
            <Button
              :label="t('form.new.step_addmore.add_button')"
              icon="pi pi-plus"
              size="large"
              :disabled="isDisabled"
              raised
              @click="increaseVulnerabilityNumber"
            />
          </div>
        </div>
      </FormKit>

      <FormKit type="step" :label="'Rank'" name="rank" stepInner-class="m-4 p-4">
        <FormKit
          type="reorderlist"
          v-model="onlyVulnerabilities"
          @node="ranking"
          @reorder="onReorder"
        >
          <template #stepNext="{ handlers, node }">
            <FormKit
              type="button"
              @click="handlers.incrementStep(1)()"
              :label="t('form.new.step_addmore.next_button', l)"
              data-next="true"
            />
          </template>
        </FormKit>
      </FormKit>
      <FormKit type="step" label="Summary" name="summary" key="summary">
        <FormResultSummaryNew :data="values" :hide-job-function="false" :is-show-summary="true" />
        <template #stepNext>
          <FormKit type="submit" />
        </template>
      </FormKit>
    </FormKit>
  </FormKit>
</template>

<style>
.formkit-outer[data-type='multi-step'] {
  max-width: none;
}

.formkit-outer[data-type='multi-step'] > .formkit-wrapper {
  max-width: none;
  box-shadow: var(--multistep-shadow);
  border-radius: var(--multistep-radius);
}
</style>

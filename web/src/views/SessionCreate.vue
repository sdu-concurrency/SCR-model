<script setup lang="ts">
import { ref, watch } from 'vue'
import { object, string, setLocale, defaultLocale } from 'yup'
import { useForm } from 'vee-validate'
import { usePocketbaseStore } from '@/stores/pb'
import { useToast } from 'primevue/usetoast'
import { ClientResponseError, type RecordModel } from 'pocketbase'
import { useI18n } from 'vue-i18n'
import yupLocaleDa from '@/yup_locales/da'
import Message from 'primevue/message'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'

const { t, locale } = useI18n()
const toast = useToast()
const { api } = usePocketbaseStore()

const schema = object().shape({
  name: string()
    .label('Session Name')
    .matches(/^\S.*$/, 'Whitespace is not allowed')
    .matches(
      /^[a-z0-9_-]+$/,
      'Session name must be lowercase and contain only letters, numbers, underscores, and hyphens'
    )
    .required(),
  description: string().label('Description').required(),
  survey: string().label('Survey').required(),
  contact_email: string().label('Contact Email').email().required(),
  contact_tel: string().label('Contact Tel.')
  // start_date: string().label('contact_email').email().required(),
  // end_date: string().label('contact_tel').required()
})

const { defineField, handleSubmit, errors, validate } = useForm({
  validationSchema: schema
})

const [name] = defineField('name')
const [description] = defineField('description')
const [survey] = defineField('survey')
const [contact_email] = defineField('contact_email')
const [contact_tel] = defineField('contact_tel')
const [start_date] = defineField('start_date')
const [end_date] = defineField('end_date')

start_date.value = new Date()
end_date.value = new Date(new Date().setDate(new Date().getDate() + 30))

const questions = await api.collection('questions').getFullList({
  sort: '-created',
  fields: 'id,name,version,language,description'
})

survey.value = questions[0].id

const insertedSession = ref<RecordModel | null>(null)
const createdSuperUser = ref<RecordModel | null>(null)
const isLoading = ref(false)
const isValidatingName = ref(false)
const isNameValidated = ref(false)
const nameValidationError = ref<string | null>(null)

// Watch for changes in the name field to revoke validation
watch(name, () => {
  if (isNameValidated.value) {
    isNameValidated.value = false
    nameValidationError.value = null
  }
})

const validateSessionName = async () => {
  if (!name.value || name.value.trim() === '') {
    nameValidationError.value = t('session_create.name_required')
    return
  }
  // Check format first
  if (!/^\S.*$/.test(name.value)) {
    nameValidationError.value = 'Whitespace is not allowed'
    return
  }

  if (!/^[a-z0-9_-]+$/.test(name.value)) {
    nameValidationError.value =
      'Session name must be lowercase and contain only letters, numbers, underscores, and hyphens'
    return
  }

  isValidatingName.value = true
  nameValidationError.value = null

  try {
    // Check if session name exists
    const existingSessions = await api.collection('sessions').getList(1, 1, {
      filter: `name="${name.value}"`
    })

    if (existingSessions.items.length > 0) {
      nameValidationError.value = t('session_create.name_exists') || 'Session name already exists'
      isNameValidated.value = false
    } else {
      isNameValidated.value = true
      toast.add({
        severity: 'success',
        summary: t('session_create.validation_success') || 'Valid',
        detail: t('session_create.name_available') || 'Session name is available',
        group: 'br',
        life: 2000
      })
    }
  } catch (e) {
    const err = new ClientResponseError(e)
    nameValidationError.value = err.message || 'Validation failed'
    isNameValidated.value = false
  } finally {
    isValidatingName.value = false
  }
}

const getSessionSuperUser = async () => {
  const sessionId = insertedSession.value?.id
  if (!sessionId) {
    console.error('session id is undefined')
    return null
  }

  const listRecord = await api
    .collection('app_users')
    .getFirstListItem(`session='${sessionId}' && role = 'super'`, { fields: 'id' })
  const record = await api.collection('app_users').getOne(listRecord.id)
  createdSuperUser.value = record
}

const onSubmit = handleSubmit(async (values) => {
  // Check if session name has been validated
  if (!isNameValidated.value) {
    toast.add({
      severity: 'warn',
      summary: t('session_create.validation_required') || 'Validation Required',
      detail:
        t('session_create.validate_name_first') ||
        'Please validate the session name before submitting',
      group: 'br',
      life: 3000
    })
    return
  }

  isLoading.value = true
  const now = new Date()
  const session = {
    name: values.name,
    description: values.description,
    current_step: 2,
    questions: values.survey,
    contact_email: values.contact_email,
    contact_tel: values.contact_tel,
    start_date: now,
    end_date: new Date(now.setMonth(now.getMonth() + 1))
  }
  try {
    const res = await api.collection('sessions').create(session)
    insertedSession.value = res
    await getSessionSuperUser()
    toast.add({
      severity: 'success',
      summary: t('session_create.success_text'),
      detail: res.record,
      group: 'br',
      life: 2000
    })
    setTimeout(() => {
      isLoading.value = false
    }, 2000)
  } catch (e) {
    const err = new ClientResponseError(e)
    toast.add({
      severity: 'error',
      summary: 'error',
      detail: err.data.message!,
      group: 'br',
      life: 2000
    })
    setTimeout(() => {
      isLoading.value = false
    }, 2000)
  }
})
watch(locale, async () => {
  if (locale.value === 'da') {
    setLocale(yupLocaleDa)
  } else {
    setLocale(defaultLocale)
  }
})
</script>
<template>
  <main class="container mx-auto">
    <Message v-if="createdSuperUser" class="mb-4">
      <section>
        <p class="text-xl md:text-2xl mb-2">{{ $t('session_create.success_text') }}:</p>
        <CopyableUserPassword
          :username="createdSuperUser.username"
          :password="createdSuperUser.pwordtext"
          link="/"
          :status="createdSuperUser.status"
        ></CopyableUserPassword>
      </section>
    </Message>
    <h1 class="text-xl md:text-2xl mb-4">{{ $t('session_create.header') }}</h1>
    <form @submit.prevent="onSubmit" class="space-y-4">
      <div class="flex flex-col gap-2">
        <label for="name">{{ $t('session_create.name_label') }} *</label>
        <div class="flex gap-2">
          <InputText
            id="name"
            v-model="name"
            aria-describedby="name-help"
            :invalid="errors.name !== undefined || nameValidationError !== null"
            :disabled="isValidatingName"
            class="flex-1"
          />
          <Button
            type="button"
            :label="$t('session_create.validate_button') || 'Validate'"
            severity="info"
            @click="validateSessionName"
            :disabled="isValidatingName || !name || name.trim() === ''"
            :loading="isValidatingName"
          />
        </div>
        <small v-if="errors.name" id="name-help" class="text-red-500">
          {{ errors.name }}
        </small>
        <small v-if="nameValidationError" class="text-red-500">
          {{ nameValidationError }}
        </small>
        <small v-if="isNameValidated" class="text-green-500">
          ✓ {{ $t('session_create.name_validated') || 'Session name validated successfully' }}
        </small>
        <small v-if="isValidatingName" class="text-blue-500">
          {{ $t('session_create.validating') || 'Validating session name...' }}
        </small>
      </div>
      <div class="flex flex-col gap-2">
        <label for="description" class="font-semibold text-sm md:text-base"
          >{{ $t('session_create.description_label') }} *</label
        >
        <InputText
          id="description"
          v-model="description"
          aria-describedby="description-help"
          :invalid="errors.description !== undefined"
          class="w-full"
        />
        <small id="description-help" class="text-red-500">
          {{ errors.description }}
        </small>
      </div>
      <div class="flex flex-col justify-center gap-2">
        <label for="questions" class="font-semibold text-sm md:text-base"
          >{{ $t('session_create.questions_label') }} *</label
        >
        <Listbox
          id="questions"
          v-model="survey"
          :options="questions"
          optionLabel="name"
          placeholder="Select Questions set"
          optionValue="id"
          class="w-full"
          :invalid="errors.survey !== undefined"
        >
          <template #option="slotProps">
            <span class="text-sm md:text-base">
              {{ slotProps.option.name }}
              <span class="italic font-light">{{
                slotProps.option.description
                  ? ` : ${slotProps.option.description}`
                  : ''
              }}</span>
            </span>
          </template>
        </Listbox>
        <small id="language-help" class="text-red-500">{{ errors.survey }}</small>
      </div>

      <div class="flex flex-col gap-2">
        <label for="contact-email" class="font-semibold text-sm md:text-base"
          >{{ $t('session_create.contact_email_label') }} *</label
        >
        <InputText
          id="contact-email"
          v-model="contact_email"
          aria-describedby="contact-email-help"
          :invalid="errors.contact_email !== undefined"
          class="w-full"
        />
        <small id="contact-email-help" class="text-red-500">{{ errors.contact_email }}</small>
      </div>

      <div class="flex flex-col gap-2">
        <label for="contact_tel" class="font-semibold text-sm md:text-base">{{
          $t('session_create.contact_tel_label')
        }}</label>
        <InputText
          id="contact_tel"
          v-model="contact_tel"
          aria-describedby="contact_tel-help"
          :invalid="errors.contact_tel !== undefined"
          class="w-full"
        />
        <small id="contact_tel-help" class="text-red-500">{{ errors.contact_tel }}</small>
      </div>

      <!-- <div class="flex flex-col gap-2">
        <label for="start_date" class="font-semibold text-sm md:text-base">{{
          $t('session_create.start_date_label')
        }}</label>
        <DatePicker
          id="start_date_calendar"
          v-model="start_date"
          :minDate="new Date()"
          :manualInput="false"
          showTime
          dateFormat="dd/mm/yy"
          hourFormat="24"
          aria-describedby="start_date-help"
          :invalid="errors.start_date !== undefined"
          class="w-full"
        />
      </div>
      <div class="flex flex-col gap-2">
        <label for="end_date" class="font-semibold text-sm md:text-base">{{
          $t('session_create.end_date_label')
        }}</label>
        <DatePicker
          id="end_date_calendar"
          v-model="end_date"
          :minDate="new Date()"
          :manualInput="false"
          showTime
          dateFormat="dd/mm/yy"
          hourFormat="24"
          aria-describedby="end_date-help"
          :invalid="errors.end_date !== undefined"
          class="w-full"
        />
      </div> -->
      <div class="my-2">
        <Button
          id="submit-button"
          :label="$t('session_create.submit_label')"
          severity="secondary"
          type="submit"
          :disabled="isLoading || !isNameValidated"
        />
        <small v-if="!isNameValidated" class="text-orange-500 ml-2">
          {{ $t('session_create.validate_name_hint') || 'Please validate the session name first' }}
        </small>
      </div>
    </form>
  </main>
</template>

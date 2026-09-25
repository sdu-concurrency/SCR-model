import { type ComposerTranslation } from 'vue-i18n'
import { toValue } from 'vue'
import i18n from '@/i18n'

export function useJobFunctionMapper(t: ComposerTranslation, vul_key: any) {
  if (!vul_key) {
    return ''
  }
  const value = toValue(vul_key)
  return t('form.step_1.job_function_options.' + value)
}

export function useVulnerabilityMapper(t: ComposerTranslation, vul_key: any) {
  if (!vul_key) {
    return ''
  }
  const value = toValue(vul_key)
  return t('form.step_2.vulnerability_options.' + value)
}

export function useCapabilityMapper(t: ComposerTranslation, cap_key: any) {
  if (!cap_key) {
    return ''
  }
  const value = toValue(cap_key)
  return t('form.step_4.capability_options.' + value)
}

export function useResponseMapper(t: ComposerTranslation, survey_name: string, res_key: any) {
  if (!res_key) {
    return ''
  }
  const value = toValue(res_key)
  const key = 'survey.' + survey_name + '.' + value
  return i18n.global.te(key) ? t(key) : String(value)
}

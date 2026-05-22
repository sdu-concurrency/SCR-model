// formkit.config.js
import { createI18nPlugin, en, da, type FormKitLocaleRegistry } from '@formkit/i18n'
import { createValidationPlugin } from '@formkit/validation'
import { extend } from '@formkit/utils'
import { required } from '@formkit/rules'
import { defaultConfig, createInput } from '@formkit/vue'
import { type FormKitNode, type FormKitOptions, type FormKitPlugin } from '@formkit/core'
import { createMultiStepPlugin } from '@formkit/addons'
import { rootClasses } from './formkit.theme'
import '@formkit/addons/css/multistep'
import {
  createLibraryPlugin,
  text,
  form,
  submit,
  select,
  hidden,
  list,
  group
} from '@formkit/inputs'
import RiskManagement from '@/components/forms/RiskManagement.vue'
import ReorderList from '@/components/forms/ReorderList.vue'
import Capabilities from '@/components/forms/Capabilities.vue'
import InplaceSelectInput from '@/components/forms/InplaceSelectInput.vue'

const library: FormKitPlugin = createLibraryPlugin({
  text,
  form,
  submit,
  hidden,
  select,
  list,
  group
})

const i18n: FormKitPlugin = createI18nPlugin(
  extend(
    { en, da },
    {
      en: {
        validation: {
          no_none(node: { name: any }) {
            return `Please select an item.`
          },
          distinct_vulnerability(node: { name: any }) {
            return `Please select distinct vulnerability.`
          },
          distinct_capability(node: { name: any }) {
            return `Please select distinct capability.`
          },
          required_note_if_other(node: { name: any }) {
            return `In case of 'Other', please elaborate.`
          },
          inplace_required(node: { name: any }) {
            return `Vulnerability is required.`
          }
        }
      },
      da: {
        validation: {
          no_none(node: { name: any }) {
            return `Vælg venligst et emne.`
          },
          distinct_vulnerability(node: { name: any }) {
            return `Vælg venligst særskilt sårbarhed.`
          },
          distinct_capability(node: { name: any }) {
            return `Vælg venligst en særskilt kapacitet.`
          },
          required_note_if_other(node: { name: any }) {
            return `I tilfælde af »Andet«, uddyb venligst.`
          },
          inplace_required(node: { name: any }) {
            return `Sårbarhed er påkrævet`
          }
        }
      }
    }
  ) as FormKitLocaleRegistry
)

const othersKeys = [
  'S1_5',
  'S2_16',
  'S3_16',
  'S4_9',
  'S5_9',
  'S6_6',
  'S7_14',
  'K1_10',
  'K2_16',
  'K3_21',
  'K4_7',
  'K5_12',
  'K6_8',
  'K7_23'
]

function required_note_if_other(node: {
  props: { id?: any }
  at: (arg0: string) => FormKitNode<unknown> | undefined
  value: unknown
}) {
  const note_node_name = 'note_' + node.props.id!.split('_')[0]
  const note_node: FormKitNode<string> = node.at(note_node_name) as FormKitNode<string>
  if (othersKeys.includes(node.value as string) && (!note_node.value || note_node.value === '')) {
    return false
  }
  return true
}

function inplace_required(node: { at: (arg0: string) => any; name: string; value: any }) {
  return node.value.value !== ''
}

function distinct_vulnerability(node: { at: (arg0: string) => any; name: string; value: any }) {
  const value = node.at('$parent').at('$parent')
  let index = 2
  const count = new Set<string>()
  while (value.value[`step_${index}`]) {
    if (count.has(value.value[`step_${index}`].vulnerability.value)) {
      return false
    }
    count.add(value.value[`step_${index}`].vulnerability.value)
    index++
  }
  return true
}

function no_none(node: { at: (arg0: string) => any; name: string; value: any }) {
  return !node.value || node.value === 'None' ? false : true
}

function distinct_capability(node: { at: (arg0: string) => any; name: string; value: any }) {
  if (!node.value && node.value === 'None') {
    return true
  }
  const capability_group = node.at('$parent').at('$parent')
  if (Array.isArray(capability_group.value)) {
    return (
      capability_group.value.filter((e: { capability: any }) => e.capability === node.value)
        .length <= 1
    )
  }
  return true
}

const validation: FormKitPlugin = createValidationPlugin({
  distinct_vulnerability,
  distinct_capability,
  no_none,
  required_note_if_other,
  required,
  inplace_required
})

const config: FormKitOptions = defaultConfig({
  plugins: [library, createMultiStepPlugin(), i18n, validation],
  config: {
    rootClasses
  },
  inputs: {
    risk_management: createInput(RiskManagement),
    inplace_select: createInput(InplaceSelectInput, {
      props: ['options']
    }),
    capabilities: createInput(Capabilities, {
      props: ['options']
    }),
    reorderlist: createInput(ReorderList)
  }
})
export default config

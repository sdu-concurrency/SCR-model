/**
 * useQuestionLabels
 *
 * Resolves a QuestionSchema (array of category groups) into a flat
 * { itemValue: label } map for the current locale. Used as FormKit :options
 * and for label lookup in summary views.
 *
 * Supports both label formats:
 *   Multilingual: { "en": "Sales", "da": "Salg" }
 *   Legacy string: "Sales"  (treated as the display value regardless of locale)
 *
 * Falls back to "en" when the requested locale key is absent, then to the
 * first available locale, then to the raw item value (avoids blank selects).
 */

import { computed, toValue, type MaybeRef } from 'vue'
import type {
  QuestionSchema,
  QuestionSchemaCategory,
  QuestionSchemaItem,
  QuestionLabelEntry
} from '@/types/question'

/**
 * Normalize a question schema stored in a slightly different shape into the
 * canonical array-of-categories format. Accepted inputs:
 *   - [{ label, items }]    canonical categories array
 *   - { label?, items }     single category object (unwrapped into an array)
 *   - [{ value, label }]    flat item list (wrapped into one category)
 * Returns [] for anything else.
 */
export function normalizeQuestionSchema(schema: unknown): QuestionSchema {
  if (!schema) return []

  // Single category object instead of an array of categories.
  if (
    typeof schema === 'object' &&
    !Array.isArray(schema) &&
    Array.isArray((schema as QuestionSchemaCategory).items)
  ) {
    return [schema as QuestionSchemaCategory]
  }

  if (Array.isArray(schema)) {
    // Flat item list (no nested "items" arrays) - wrap it in one category.
    const looksFlat =
      schema.length > 0 &&
      schema.every((entry) => {
        const c = entry as Partial<QuestionSchemaCategory> & Partial<QuestionSchemaItem>
        return c && !Array.isArray(c.items) && 'value' in c
      })
    if (looksFlat) {
      return [{ label: '', items: schema as QuestionSchemaItem[] }]
    }
    return schema as QuestionSchema
  }

  return []
}

/**
 * Resolve a single label entry to a string for the given locale.
 */
export function resolveLabel(
  entry: QuestionLabelEntry | undefined,
  locale: string,
  fallback = ''
): string {
  if (!entry) return fallback
  if (typeof entry === 'string') return entry
  return entry[locale] ?? entry['en'] ?? Object.values(entry)[0] ?? fallback
}

/**
 * Composable: returns a computed flat { value → label } map for the current locale.
 * Iterates all categories and their items; use the result as FormKit <select> :options.
 */
export function useQuestionLabels(
  schema: MaybeRef<QuestionSchema | null | undefined>,
  locale: MaybeRef<string>
) {
  return computed(() => {
    const categories = normalizeQuestionSchema(toValue(schema))
    const lang = toValue(locale) || 'en'

    if (!Array.isArray(categories)) return {} as Record<string, string>

    const result: Record<string, string> = {}
    for (const category of categories) {
      if (!Array.isArray(category.items)) continue
      for (const item of category.items) {
        result[item.value] = resolveLabel(item.label, lang, item.value)
      }
    }
    return result
  })
}

/**
 * Resolve a flat map imperatively (outside of a reactive context).
 * Useful in script-setup computed chains or migration helpers.
 */
export function resolveSchema(
  schema: QuestionSchema | null | undefined,
  locale: string
): Record<string, string> {
  const categories = normalizeQuestionSchema(schema)
  const result: Record<string, string> = {}
  for (const category of categories) {
    if (!Array.isArray(category.items)) continue
    for (const item of category.items) {
      result[item.value] = resolveLabel(item.label, locale, item.value)
    }
  }
  return result
}

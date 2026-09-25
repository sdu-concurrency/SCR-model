/**
 * Question Types
 *
 * Type definitions for question schemas used in surveys.
 *
 * Actual DB structure: schemas are arrays of category groups, each with items.
 * Labels are either plain strings (legacy, single-language) or multilingual objects.
 *
 * Legacy:       { value: "S1_1", label: "Asset turnover" }
 * Multilingual: { value: "S1_1", label: { "en": "Asset turnover", "da": "Aktivers omsætningshastighed" } }
 */

/** A label that is either a plain string or a locale->string map. */
export type QuestionLabelEntry = string | Record<string, string>

/** A single selectable item inside a category (e.g. S1_1, K2_3). */
export interface QuestionSchemaItem {
  value: string
  label: QuestionLabelEntry
}

/** A category group containing a set of items (e.g. "1. Finance"). */
export interface QuestionSchemaCategory {
  label: QuestionLabelEntry
  items: QuestionSchemaItem[]
}

/** The full schema for vulnerabilities, capabilities or job functions - an ordered array of categories. */
export type QuestionSchema = QuestionSchemaCategory[]

/**
 * Helper: flat map of { itemValue -> localised label string }.
 * Used by FormKit <select> :options and label-resolution in summary/analysis views.
 */
export type QuestionLabelMap = Record<string, string>

/**
 * Question record from PocketBase.
 */
export interface Question {
  id: string
  name: string
  description: string
  version: string // e.g. "v1"
  vulnerability_schema: QuestionSchema
  capability_schema: QuestionSchema
  job_function_schema?: QuestionSchema // optional - older records may not have it
  created: string
  updated: string
}

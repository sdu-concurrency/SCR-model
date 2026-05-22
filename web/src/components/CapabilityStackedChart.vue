<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import Chart from 'primevue/chart'
import { useI18n } from 'vue-i18n'
import { useResponseMapper } from '../composables/responseMapper.ts'

// ============================================================================
// Constants
// ============================================================================

const COLOR_PALETTE = [
  '#06b6d4',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#ec4899',
  '#3b82f6',
  '#84cc16',
  '#f97316',
  '#6366f1',
  '#14b8a6',
  '#a855f7',
  '#22c55e',
  '#f43f5e',
  '#0ea5e9',
  '#eab308',
  '#dc2626',
  '#9333ea',
  '#0891b2',
  '#059669',
  '#d97706',
  '#7c3aed',
  '#db2777',
  '#2563eb',
  '#65a30d',
  '#ea580c',
  '#4f46e5',
  '#0d9488',
  '#c026d3'
]

const CATEGORY_LABELS_KEYS = {
  K1: 'capability_stacked_chart.chart_title_k1',
  K2: 'capability_stacked_chart.chart_title_k2',
  K3: 'capability_stacked_chart.chart_title_k3',
  K4: 'capability_stacked_chart.chart_title_k4',
  K5: 'capability_stacked_chart.chart_title_k5',
  K6: 'capability_stacked_chart.chart_title_k6',
  K7: 'capability_stacked_chart.chart_title_k7'
}

const CATEGORY_ORDER = ['K1', 'K2', 'K3', 'K4', 'K5', 'K6', 'K7']

const COUNT_LABELS = [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]

const CHART_CONFIG = {
  printWidth: 600,
  printHeight: 400,
  defaultHeight: 800,
  mobileHeight: 600,
  barPercentage: 0.95,
  categoryPercentage: 0.95,
  borderWidth: 1
}

// ============================================================================
// Props
// ============================================================================

const props = defineProps({
  capabilitiesData: {
    type: Array,
    required: true
  },
  sessionName: {
    type: String,
    required: true
  }
})

// ============================================================================
// Composables
// ============================================================================

const { t } = useI18n()

// ============================================================================
// State
// ============================================================================

const chartRefs = ref({})

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Extracts the category key (K1, K2, etc.) from a capability string
 */
function extractCategoryKey(capability) {
  return capability.substring(0, 2)
}

/**
 * Gets the color for a capability based on its index
 */
function getCapabilityColor(capabilityIndex) {
  return COLOR_PALETTE[capabilityIndex % COLOR_PALETTE.length]
}

/**
 * Calculates the average value from a sum and count
 */
function calculateAverage(sum, count) {
  return count > 0 ? sum / count : 0
}

/**
 * Aggregates capabilities by category from vulnerability data
 */
function aggregateCapabilitiesByCategory(vulnerabilities) {
  const categories = {}
  const capabilityIndices = new Map()

  vulnerabilities.forEach((vulnerability) => {
    vulnerability.capability_list?.forEach((capability) => {
      const categoryKey = extractCategoryKey(capability.capability)
      const capabilityKey = capability.capability

      // Initialize category if not exists
      if (!categories[categoryKey]) {
        categories[categoryKey] = {}
      }

      // Initialize capability aggregation if not exists
      if (!categories[categoryKey][capabilityKey]) {
        categories[categoryKey][capabilityKey] = {
          count: 0,
          importanceSum: 0,
          currentAbilitySum: 0,
          jobFunctions: []
        }
      }

      // Aggregate capability data
      const aggregation = categories[categoryKey][capabilityKey]
      const jobFunctionsCount = capability.job_functions?.length || 0

      aggregation.count += jobFunctionsCount
      aggregation.importanceSum += capability.importance_avg || 0
      aggregation.currentAbilitySum += capability.current_ability_avg || 0
      aggregation.jobFunctions.push(...(capability.job_functions || []))

      // Track capability index for color assignment
      if (!capabilityIndices.has(capabilityKey)) {
        capabilityIndices.set(capabilityKey, capabilityIndices.size)
      }
    })
  })

  return { categories, capabilityIndices }
}

/**
 * Creates a dataset for a single capability
 */
function createCapabilityDataset(capabilityKey, capabilityData, capabilityIndex, countLabels) {
  const jobFunctionsCount = capabilityData.jobFunctions.length || 1
  const avgImportance = calculateAverage(capabilityData.importanceSum, jobFunctionsCount)
  const avgCurrentAbility = calculateAverage(capabilityData.currentAbilitySum, jobFunctionsCount)

  const minValue = Math.min(avgImportance, avgCurrentAbility)
  const maxValue = Math.max(avgImportance, avgCurrentAbility)

  // Create data array with range values only for the capability's count
  const data = countLabels.map((count) => {
    return count === capabilityData.count ? [minValue, maxValue] : null
  })

  const color = getCapabilityColor(capabilityIndex)

  return {
    label: useResponseMapper(t, 'V1', capabilityKey),
    data,
    backgroundColor: color,
    borderColor: color,
    borderWidth: CHART_CONFIG.borderWidth,
    barThickness: 'flex',
    maxBarThickness: 15,
    categoryPercentage: 1.0,
    barPercentage: 0.9,
    skipNull: true,
    capabilityKey,
    avgImportance: avgImportance.toFixed(2),
    avgCurrentAbility: avgCurrentAbility.toFixed(2),
    minValue: minValue.toFixed(2),
    maxValue: maxValue.toFixed(2),
    peopleCount: capabilityData.count,
    categoryData: capabilityData
  }
}

/**
 * Creates chart data for a specific category
 */
function createCategoryChartData(categoryKey, categoryCapabilities, capabilityIndices) {
  const datasets = []
  const capabilityKeys = Object.keys(categoryCapabilities)

  capabilityKeys.forEach((capabilityKey) => {
    const capabilityData = categoryCapabilities[capabilityKey]
    const capabilityIndex = capabilityIndices.get(capabilityKey) || 0

    const dataset = createCapabilityDataset(
      capabilityKey,
      capabilityData,
      capabilityIndex,
      COUNT_LABELS
    )

    // Skip capabilities where importance and current ability are equal (zero-width bar)
    if (dataset.avgImportance !== dataset.avgCurrentAbility) {
      datasets.push(dataset)
    }
  })

  return {
    labels: COUNT_LABELS.map(String),
    datasets,
    categoryData: categoryCapabilities,
    categoryTitle: t(CATEGORY_LABELS_KEYS[categoryKey]) || categoryKey
  }
}

// ============================================================================
// Computed Properties
// ============================================================================

/**
 * Processes capabilities data and generates chart data for each category
 */
const categoryCharts = computed(() => {
  const { categories, capabilityIndices } = aggregateCapabilitiesByCategory(props.capabilitiesData)
  const charts = {}

  CATEGORY_ORDER.forEach((categoryKey) => {
    if (categories[categoryKey]) {
      charts[categoryKey] = createCategoryChartData(
        categoryKey,
        categories[categoryKey],
        capabilityIndices
      )
    }
  })

  return charts
})

// ============================================================================
// Chart Configuration
// ============================================================================

/**
 * Generates chart options for a category chart
 */
const getChartOptions = (categoryTitle) => ({
  indexAxis: 'y',
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  plugins: {
    datalabels: {
      display: false
    },
    legend: {
      display: true,
      position: 'top',
      labels: {
        font: {
          size: 10
        }
      }
    },
    tooltip: {
      callbacks: {
        title: (context) => context[0].dataset.label,
        label: (context) => {
          const dataset = context.dataset
          const capabilityData = dataset.categoryData

          return [
            `${t('capability_stacked_chart.tooltip_people')}: ${dataset.peopleCount}`,
            `${t('capability_stacked_chart.tooltip_avg_importance')}: ${dataset.avgImportance}`,
            `${t('capability_stacked_chart.tooltip_avg_current_ability')}: ${dataset.avgCurrentAbility}`
          ]
        }
      }
    },
    title: {
      display: true,
      text: categoryTitle,
      font: {
        size: 14,
        weight: 'bold'
      },
      padding: 10
    }
  },
  scales: {
    x: {
      min: 1,
      max: 5,
      title: {
        display: true,
        text: t('capability_stacked_chart.x_axis_label'),
        font: {
          size: 11,
          weight: 'bold'
        }
      },
      ticks: {
        precision: 0,
        stepSize: 1,
        font: {
          size: 10
        }
      }
    },
    y: {
      title: {
        display: true,
        text: t('capability_stacked_chart.y_axis_label'),
        font: {
          size: 11,
          weight: 'bold'
        }
      },
      ticks: {
        autoSkip: false,
        font: {
          size: 10
        }
      }
    }
  }
})

// ============================================================================
// Print Handlers
// ============================================================================

/**
 * Resizes charts for print layout
 */
const handleBeforePrint = () => {
  Object.values(chartRefs.value).forEach((chartRef) => {
    if (chartRef?.getChart) {
      chartRef.getChart().resize(CHART_CONFIG.printWidth, CHART_CONFIG.printHeight)
    }
  })
}

/**
 * Restores chart sizes after printing
 */
const handleAfterPrint = () => {
  Object.values(chartRefs.value).forEach((chartRef) => {
    if (chartRef?.getChart) {
      chartRef.getChart().resize()
    }
  })
}

// ============================================================================
// Lifecycle Hooks
// ============================================================================

onMounted(() => {
  window.addEventListener('beforeprint', handleBeforePrint)
  window.addEventListener('afterprint', handleAfterPrint)
})

onUnmounted(() => {
  window.removeEventListener('beforeprint', handleBeforePrint)
  window.removeEventListener('afterprint', handleAfterPrint)
})
</script>

<template>
  <div class="capability-category-charts">
    <!-- Chart Description -->
    <section class="chart-description">
      <p class="description-primary">
        {{ t('capability_stacked_chart.description_primary') }}
      </p>
      <p class="description-secondary">
        {{ t('capability_stacked_chart.description_secondary') }}
      </p>
    </section>

    <!-- Charts Grid -->
    <div class="charts-grid">
      <article
        v-for="(chartData, categoryKey) in categoryCharts"
        :key="categoryKey"
        class="chart-card"
      >
        <div class="chart-container">
          <Chart
            :ref="(el) => (chartRefs[categoryKey] = el)"
            type="bar"
            :data="chartData"
            :options="getChartOptions(chartData.categoryTitle)"
            class="chart"
          />
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.capability-category-charts {
  width: 100%;
}

/* Description Section */
.chart-description {
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  color: #4b5563;
}

.description-primary {
  margin-bottom: 0.5rem;
}

.description-secondary {
  margin: 0;
}

/* Charts Grid */
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(500px, 100%), 1fr));
  gap: 2rem;
  align-items: start;
}

/* Chart Card */
.chart-card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  padding: 1rem;
  margin-bottom: 2rem;
  min-width: 0;
}

/* Chart Container */
.chart-container {
  height: 500px;
  position: relative;
}

.chart {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

/* Print Styles */
@media print {
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .chart-card {
    break-inside: avoid;
    page-break-inside: avoid;
    margin-bottom: 1rem;
    padding: 0.5rem;
  }

  .chart-container {
    height: 400px !important;
  }

  .chart-description {
    font-size: 0.75rem;
    margin-bottom: 1rem;
  }
}

/* Mobile Styles */
@media (max-width: 768px) {
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .chart-container {
    height: 400px;
  }
}
</style>

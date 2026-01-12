<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex flex-col gap-2">
      <h1 class="text-3xl font-bold">Research Dashboard</h1>
      <p class="text-[hsl(var(--maz-foreground))]/60">
        Public insights from the AV Ethics experiment
      </p>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        v-for="kpi in kpiCards"
        :key="kpi.label"
        class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm"
      >
        <div class="flex items-center gap-4">
          <div
            class="w-12 h-12 rounded-xl flex items-center justify-center"
            :class="kpi.bgClass"
          >
            <component :is="kpi.icon" class="w-6 h-6" :class="kpi.iconClass" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ kpi.value }}</p>
            <p class="text-sm text-[hsl(var(--maz-foreground))]/60">{{ kpi.label }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Row 1 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Least Harmful Outcome - Donut -->
      <div class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm">
        <h3 class="text-lg font-semibold mb-4">Least Harmful Outcome Choices</h3>
        <div class="h-64 flex items-center justify-center">
          <Doughnut :key="chartKey" :data="outcomeChartData" :options="doughnutOptions" />
        </div>
      </div>

      <!-- Decision Time Distribution - Line -->
      <div class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm">
        <h3 class="text-lg font-semibold mb-4">Decision Time Distribution</h3>
        <div class="h-64">
          <Line :key="chartKey" :data="timeChartData" :options="lineOptions" />
        </div>
      </div>
    </div>

    <!-- Charts Row 2 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Self Preservation Effect - Horizontal Bar -->
      <div class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm">
        <h3 class="text-lg font-semibold mb-4">Self-Preservation Effect</h3>
        <p class="text-sm text-[hsl(var(--maz-foreground))]/60 mb-4">
          % choosing to maintain course (self-preserving) based on tailgater presence
        </p>
        <div class="h-48">
          <Bar :key="chartKey" :data="preservationChartData" :options="horizontalBarOptions" />
        </div>
      </div>

      <!-- Entity Compliance Effect - Horizontal Bar -->
      <div class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm">
        <h3 class="text-lg font-semibold mb-4">Entity Compliance Effect</h3>
        <p class="text-sm text-[hsl(var(--maz-foreground))]/60 mb-4">
          % choosing to maintain course based on entity's rule compliance
        </p>
        <div class="h-48">
          <Bar :key="chartKey" :data="complianceChartData" :options="horizontalBarOptions" />
        </div>
      </div>
    </div>

    <!-- Archetype Distribution -->
    <div class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm">
      <h3 class="text-lg font-semibold mb-4">Decision-Making Archetypes</h3>
      <p class="text-sm text-[hsl(var(--maz-foreground))]/60 mb-6">
        Personality profiles based on response patterns
      </p>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="(archetype, index) in dashboardData.archetype_distribution"
          :key="archetype.archetype"
          class="p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))] flex items-center gap-4"
        >
          <div
            class="w-10 h-10 rounded-full flex items-center justify-center text-lg font-bold"
            :class="getArchetypeColor(index)"
          >
            {{ index + 1 }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="font-medium truncate">{{ archetype.archetype }}</p>
            <p class="text-sm text-[hsl(var(--maz-foreground))]/60">
              {{ archetype.count }} participant{{ archetype.count !== 1 ? 's' : '' }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { Doughnut, Line, Bar } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    LineElement,
    BarElement,
    PointElement,
    CategoryScale,
    LinearScale,
    Filler
} from 'chart.js'
import { MazUsers, MazGlobeAlt } from '@maz-ui/icons'

ChartJS.register(
    Title,
    Tooltip,
    Legend,
    ArcElement,
    LineElement,
    BarElement,
    PointElement,
    CategoryScale,
    LinearScale,
    Filler
)

// Dummy data - replace with API call later
const dashboardData = {
  completed_sessions: 3,
  countries_represented: 3,
  least_harmful_outcome: {
    maintain: 3,
    swerve_left: 1,
    swerve_right: 2,
    total: 6,
    maintain_pct: 50,
    swerve_left_pct: 16.67,
    swerve_right_pct: 33.33
  },
  self_preservation_effect: {
    with_tailgater: { maintain_count: 0, total_count: 1, percentage: 0 },
    without_tailgater: { maintain_count: 3, total_count: 5, percentage: 60 }
  },
  entity_compliance_effect: {
    compliant: { maintain_count: 3, total_count: 3, percentage: 100 },
    violation: { maintain_count: 0, total_count: 3, percentage: 0 }
  },
  decision_time_distribution: [
    { seconds: 1, count: 2 },
    { seconds: 2, count: 5 },
    { seconds: 3, count: 12 },
    { seconds: 4, count: 18 },
    { seconds: 5, count: 15 },
    { seconds: 6, count: 8 },
    { seconds: 7, count: 4 },
    { seconds: 8, count: 3 },
    { seconds: 9, count: 1 },
    { seconds: 10, count: 1 },
    { seconds: 11, count: 1 }
  ],
  archetype_distribution: [
    { archetype: 'The Hesitant Observer', count: 2 },
    { archetype: 'The Lawful Protector', count: 1 }
  ]
}

// Dynamic theme colors
const colors = ref({
  primary: '',
  accent: '',
  success: '',
  warning: '',
  danger: '',
  info: '',
  foreground: '',
  muted: ''
})

const chartKey = ref(0)

function hslToRgb(hslString: string): string {
  const parts = hslString.trim().split(/\s+/)
  if (parts.length < 3) return 'rgb(128, 128, 128)'
  
  const h = parseFloat(parts[0]) / 360
  const s = parseFloat(parts[1]) / 100
  const l = parseFloat(parts[2]) / 100

  let r, g, b
  if (s === 0) {
    r = g = b = l
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1
      if (t > 1) t -= 1
      if (t < 1/6) return p + (q - p) * 6 * t
      if (t < 1/2) return q
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6
      return p
    }
    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q
    r = hue2rgb(p, q, h + 1/3)
    g = hue2rgb(p, q, h)
    b = hue2rgb(p, q, h - 1/3)
  }
  return `rgb(${Math.round(r * 255)}, ${Math.round(g * 255)}, ${Math.round(b * 255)})`
}

function withAlpha(rgb: string, alpha: number): string {
  return rgb.replace('rgb', 'rgba').replace(')', `, ${alpha})`)
}

function updateColors() {
  const style = getComputedStyle(document.documentElement)
  
  colors.value = {
    primary: hslToRgb(style.getPropertyValue('--maz-primary').trim()),
    accent: hslToRgb(style.getPropertyValue('--maz-accent').trim()),
    success: hslToRgb(style.getPropertyValue('--maz-success').trim()),
    warning: hslToRgb(style.getPropertyValue('--maz-warning').trim()),
    danger: hslToRgb(style.getPropertyValue('--maz-danger').trim()),
    info: hslToRgb(style.getPropertyValue('--maz-info').trim()),
    foreground: hslToRgb(style.getPropertyValue('--maz-foreground').trim()),
    muted: hslToRgb(style.getPropertyValue('--maz-foreground').trim())
  }
  
  chartKey.value++
}

onMounted(() => {
  updateColors()
  
  const observer = new MutationObserver(() => {
    setTimeout(updateColors, 50)
  })
  
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  })
})

// KPI Cards
const kpiCards = computed(() => [
  {
    label: 'Completed Sessions',
    value: dashboardData.completed_sessions,
    icon: MazUsers,
    bgClass: 'bg-[hsl(var(--maz-primary))]/10',
    iconClass: 'text-[hsl(var(--maz-primary))]'
  },
  {
    label: 'Countries Represented',
    value: dashboardData.countries_represented,
    icon: MazGlobeAlt,
    bgClass: 'bg-[hsl(var(--maz-accent))]/10',
    iconClass: 'text-[hsl(var(--maz-accent))]'
  }
])

// Donut Chart
const outcomeChartData = computed(() => ({
  labels: ['Maintain Course', 'Swerve Left', 'Swerve Right'],
  datasets: [{
    data: [
      dashboardData.least_harmful_outcome.maintain,
      dashboardData.least_harmful_outcome.swerve_left,
      dashboardData.least_harmful_outcome.swerve_right
    ],
    backgroundColor: [colors.value.primary, colors.value.accent, colors.value.info],
    borderColor: 'transparent',
    hoverOffset: 4
  }]
}))

const doughnutOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: { color: colors.value.foreground, padding: 16 }
    }
  },
  cutout: '60%'
}))

// Line Chart
const timeChartData = computed(() => ({
  labels: dashboardData.decision_time_distribution.map(d => `${d.seconds}s`),
  datasets: [{
    label: 'Responses',
    data: dashboardData.decision_time_distribution.map(d => d.count),
    borderColor: colors.value.primary,
    backgroundColor: withAlpha(colors.value.primary, 0.1),
    fill: true,
    tension: 0.4,
    pointBackgroundColor: colors.value.primary,
    pointBorderColor: colors.value.primary,
    pointRadius: 4,
    pointHoverRadius: 6
  }]
}))

const lineOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: {
      grid: { color: withAlpha(colors.value.muted, 0.1) },
      ticks: { color: withAlpha(colors.value.muted, 0.7) }
    },
    y: {
      grid: { color: withAlpha(colors.value.muted, 0.1) },
      ticks: { color: withAlpha(colors.value.muted, 0.7) },
      beginAtZero: true
    }
  }
}))

// Horizontal Bars
const preservationChartData = computed(() => ({
  labels: ['With Tailgater', 'Without Tailgater'],
  datasets: [{
    label: 'Maintain Course %',
    data: [
      dashboardData.self_preservation_effect.with_tailgater.percentage,
      dashboardData.self_preservation_effect.without_tailgater.percentage
    ],
    backgroundColor: [colors.value.danger, colors.value.success],
    borderRadius: 8,
    barThickness: 32
  }]
}))

const complianceChartData = computed(() => ({
  labels: ['Compliant Entity', 'Violating Entity'],
  datasets: [{
    label: 'Maintain Course %',
    data: [
      dashboardData.entity_compliance_effect.compliant.percentage,
      dashboardData.entity_compliance_effect.violation.percentage
    ],
    backgroundColor: [colors.value.success, colors.value.warning],
    borderRadius: 8,
    barThickness: 32
  }]
}))

const horizontalBarOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  indexAxis: 'y' as const,
  plugins: { legend: { display: false } },
  scales: {
    x: {
      grid: { color: withAlpha(colors.value.muted, 0.1) },
      ticks: { color: withAlpha(colors.value.muted, 0.7) },
      max: 100,
      beginAtZero: true
    },
    y: {
      grid: { display: false },
      ticks: { color: colors.value.foreground }
    }
  }
}))

// Archetype colors
function getArchetypeColor(index: number): string {
  const colorClasses = [
    'bg-[hsl(var(--maz-primary))]/20 text-[hsl(var(--maz-primary))]',
    'bg-[hsl(var(--maz-accent))]/20 text-[hsl(var(--maz-accent))]',
    'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]',
    'bg-[hsl(var(--maz-info))]/20 text-[hsl(var(--maz-info))]',
    'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]'
  ]
  return colorClasses[index % colorClasses.length]
}
</script>

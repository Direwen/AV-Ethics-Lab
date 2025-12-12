import tailwindcss from "@tailwindcss/vite";

const mazUiConfig = {
  // General Settings
  general: {
    autoImportPrefix: 'Maz', // useMazToast instead of useToast
    defaultMazIconPath: '/icons', // Default path for <MazIcon />
    devtools: true, // Enable DevTools integration
  },

  // CSS & Styling
  css: {
    injectMainCss: true, // Auto-inject Maz-UI styles
  },

  // Theming System
  theme: {
    preset: 'maz-ui', // 'maz-ui' | 'dark' | 'ocean' | custom object
    strategy: 'hybrid', // 'runtime' | 'buildtime' | 'hybrid'
    darkModeStrategy: 'class', // 'class' | 'media' | 'auto'
    mode: 'both', // 'light' | 'dark' | 'both'
    colorMode: 'auto',
  },

  // Components
  components: {
    autoImport: true, // All components globally available
  },

  composables: {
      useToast: true,
  },

  // Plugins (not enabled by default)
  plugins: {
    aos: true,
    dialog: true,
    toast: true,
    wait: true,
  },

  // Directives (not enabled by default)
  directives: {
    vTooltip: true,
    vClickOutside: true,
  },
} as const

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  css: ['./app/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  modules: ['@maz-ui/nuxt', '@pinia/nuxt'],
  mazUi: mazUiConfig,
  runtimeConfig: {
    public: {
      apiBase: "http://localhost:8080"
    }
  }
} as any)
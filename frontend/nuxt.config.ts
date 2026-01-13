import tailwindcss from "@tailwindcss/vite";

const mazUiConfig = {
  general: {
    autoImportPrefix: 'Maz',
    defaultMazIconPath: '/icons',
    devtools: true,
  },

  css: {
    injectMainCss: true,
  },

  theme: {
    preset: 'maz-ui',
    strategy: 'hybrid',
    darkModeStrategy: 'class',
    mode: 'both',
    defaultMode: 'dark',
    colorMode: 'manual',
    overrides: {
      fontFamily: {
        base: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
        title: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
      },

      colors: {
        // Dark: Sensor/simulation dashboard aesthetic
        dark: {
          primary: '189 94% 43%',    // Cyan-500
          secondary: '215 28% 17%',  // Slate-800
          accent: '45 93% 58%',      // Amber-400 (warm contrast to cyan)
          background: '222 47% 5%',  // Slate-950
          'background-accent': '195 70% 15%', // Dark with stronger cyan glow
          foreground: '210 40% 96%', // Slate-100
          success: '158 64% 52%',    // Emerald-500
          warning: '45 93% 47%',     // Amber-500
          danger: '0 84% 60%',       // Red-500
          info: '217 91% 60%',       // Blue-500
        },

        // Light: Clean research/report aesthetic
        light: {
          primary: '221 83% 53%',    // Blue-600
          secondary: '210 40% 96%',  // Slate-100
          accent: '262 83% 58%',     // Violet-500 (energetic complement to blue)
          background: '0 0% 100%',   // White
          'background-accent': '221 80% 92%', // White with stronger blue tint
          foreground: '222 47% 11%', // Slate-900
          success: '161 94% 30%',    // Emerald-700
          warning: '37 91% 40%',     // Amber-600
          danger: '0 72% 51%',       // Red-600
          info: '221 83% 53%',       // Blue-600
        },
      },

      borderRadius: {
        DEFAULT: '0.375rem', // 6px
        lg: '0.5rem',        // 8px
        xl: '0.75rem',       // 12px
      },
    },
  },

  components: {
    autoImport: true,
  },

  composables: {
    useToast: true,
  },

  plugins: {
    aos: true,
    dialog: true,
    toast: true,
    wait: true,
  },

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
      apiBase: "http://localhost:8080",
      tridentZoneDistance: 1,
      timerDuration: 20
    }
  }
} as any)
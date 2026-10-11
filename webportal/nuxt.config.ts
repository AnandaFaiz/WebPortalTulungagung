// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ["@nuxt/ui", "@nuxt/fonts"],
  css: ["~/assets/css/main.css"],
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  nitro: {
    externals: {
      inline: [/[\\/]node_modules[\\/]nuxt[\\/]dist[\\/]/],
    },
  },
  fonts: {
    families: [
      {
        name: "Plus Jakarta Sans",
        provider: "google",
        weights: [400, 500, 600, 700, 800],
      },
    ],
  },
});

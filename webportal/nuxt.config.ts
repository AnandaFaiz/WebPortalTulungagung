// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ["@nuxt/ui"],
  css: ["~/assets/css/main.css"],
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  nitro: {
    externals: {
      inline: [/[\\/]node_modules[\\/]nuxt[\\/]dist[\\/]/],
    },
  },
});

<script setup>
definePageMeta({
    layout: false
})

const { login, loading, initialized } = useAuth()

const form = ref({
    email: '',
    password: ''
})

const showPassword = ref(false)
const errorMsg = ref('')

// ==== Quick Fill ====
const quickAccounts = [
    { label: 'Super Admin', email: 'superadmin@tulungagung.go.id' },
    { label: 'Admin OPD', email: 'admin.dinkes@tulungagung.go.id' },
    { label: 'Editor', email: 'editor@tulungagung.go.id' }
]

const fillAccount = (email) => {
    form.value.email = email
    form.value.password = ''
    errorMsg.value = ''
}

// ==== Submit ====
const handleLogin = async () => {
    errorMsg.value = ''
    if (!form.value.email || !form.value.password) {
        errorMsg.value = 'Email dan password wajib diisi.'
        return
    }

    const result = await login(form.value.email, form.value.password)

    if (result.success) {
        await navigateTo('/admin', { replace: true })
    } else {
        errorMsg.value = result.message || 'Login gagal.'
    }
}

// Tampilkan loading HANYA kalau belum login DAN state belum di-initialize
const showLoading = computed(() => !initialized.value)

onMounted(async () => {
    if (!initialized.value) {
        await fetchUser()
    }
})
</script>

<template>
    <!-- Loading screen sebelum auth selesai -->
    <div v-if="showLoading" class="min-h-screen flex items-center justify-center bg-sky-50">
        <div class="text-center">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat...</p>
        </div>
    </div>

    <!-- Konten login (hanya muncul kalau auth sudah dicek & user belum login) -->
    <div v-else class="min-h-screen relative flex items-center justify-center p-4">
        <!-- Background Image -->
        <div class="absolute inset-0 z-0">
            <img src="/login.webp" alt="Background" class="w-full h-full object-cover" />
            <!-- Overlay gradient -->
            <div class="absolute inset-0 bg-linear-to-br from-sky-900/70 via-sky-800/60 to-sky-950/80" />
        </div>

        <!-- Main Card -->
        <div
            class="relative z-10 w-full max-w-5xl bg-white/10 backdrop-blur-xl rounded-3xl border border-white/20 overflow-hidden shadow-2xl">
            <div class="grid grid-cols-1 md:grid-cols-2">

                <!-- ==== KIRI: Branding ==== -->
                <div class="relative p-10 md:p-12 flex flex-col justify-between text-white">
                    <!-- Logo -->
                    <div class="flex items-center gap-3">
                        <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo" class="h-14 w-auto drop-shadow-lg" />
                        <div class="leading-tight">
                            <p class="text-xs text-white/70">Pemerintah Kabupaten</p>
                            <p class="text-lg font-extrabold tracking-tight">TULUNGAGUNG</p>
                        </div>
                    </div>

                    <!-- Tagline -->
                    <div class="my-8 md:my-0">
                        <h1 class="text-4xl md:text-5xl font-extrabold leading-tight drop-shadow-lg">
                            Admin Panel
                            <br />
                            <span class="text-sky-300">Portal Pemkab</span>
                        </h1>
                        <p class="mt-4 text-white/80 leading-relaxed max-w-md">
                            Kelola konten, layanan publik, dan informasi portal Pemkab Tulungagung dalam satu dashboard
                            terpadu.
                        </p>
                    </div>

                    <!-- Footer kecil -->
                    <p class="text-xs text-white/50">
                        © {{ new Date().getFullYear() }} Pemkab Tulungagung
                    </p>
                </div>

                <!-- ==== KANAN: Form Login ==== -->
                <div class="bg-white/95 backdrop-blur-lg p-8 md:p-10">
                    <div class="mb-6">
                        <h2 class="text-2xl font-bold text-slate-900">Login Admin</h2>
                        <p class="text-sm text-slate-500 mt-1">Gunakan email dan password Anda untuk melanjutkan.</p>
                    </div>

                    <!-- Error Alert -->
                    <UAlert v-if="errorMsg" color="error" variant="subtle" icon="i-tabler-alert-circle"
                        :title="errorMsg" class="mb-5" />

                    <!-- Form -->
                    <form class="space-y-4" @submit.prevent="handleLogin">
                        <UFormField label="Email" required>
                            <UInput v-model="form.email" type="email" icon="i-tabler-mail"
                                placeholder="nama@tulungagung.go.id" size="lg" class="w-full" autocomplete="email" />
                        </UFormField>

                        <UFormField label="Password" required>
                            <UInput v-model="form.password" :type="showPassword ? 'text' : 'password'"
                                icon="i-tabler-lock" placeholder="••••••••" size="lg" class="w-full"
                                :ui="{ trailing: 'pe-1' }" autocomplete="current-password">
                                <template #trailing>
                                    <UButton :icon="showPassword ? 'i-tabler-eye-off' : 'i-tabler-eye'" color="neutral"
                                        variant="link" size="sm"
                                        :aria-label="showPassword ? 'Hide password' : 'Show password'"
                                        @click="showPassword = !showPassword" />
                                </template>
                            </UInput>
                        </UFormField>

                        <div class="flex items-center justify-between">
                            <UCheckbox label="Ingat saya" />
                            <a href="#" class="text-sm text-sky-600 hover:text-sky-700 font-medium">
                                Lupa password?
                            </a>
                        </div>

                        <UButton type="submit" label="Login" icon="i-tabler-login" color="primary" size="lg" block
                            :loading="loading" />
                    </form>

                    <!-- Divider -->
                    <div class="flex items-center gap-3 my-6">
                        <div class="flex-1 h-px bg-slate-200" />
                        <span class="text-xs text-slate-400">Akun Demo</span>
                        <div class="flex-1 h-px bg-slate-200" />
                    </div>

                    <!-- Quick Fill Buttons -->
                    <div class="grid grid-cols-3 gap-2">
                        <UButton v-for="acc in quickAccounts" :key="acc.email" :label="acc.label" color="neutral"
                            variant="soft" size="sm" block @click="fillAccount(acc.email)" />
                    </div>
                    <p class="text-xs text-slate-400 text-center mt-3">
                        Klik tombol di atas untuk auto-fill kredensial demo.
                    </p>

                    <!-- Back to Home -->
                    <div class="mt-6 text-center">
                        <NuxtLink to="/"
                            class="text-sm text-slate-500 hover:text-sky-600 inline-flex items-center gap-1">
                            <UIcon name="i-tabler-arrow-left" class="w-4 h-4" />
                            Kembali ke Portal
                        </NuxtLink>
                    </div>
                </div>

            </div>
        </div>
    </div>
</template>
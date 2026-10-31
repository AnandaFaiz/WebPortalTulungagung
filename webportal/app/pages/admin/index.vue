<script setup>
definePageMeta({
  layout: 'admin',
  title: 'Dashboard'
})

const { user } = useAuth()
const requestFetch = useRequestFetch()
const { data: dashboardData, pending: loading, error } = await useAsyncData('admin-dashboard', () =>
  requestFetch('http://localhost:8080/api/admin/dashboard', { credentials: 'include' })
)

const stats = computed(() => {
  const values = dashboardData.value?.stats || {}
  return [
    { label: 'Total Berita', value: values.news || 0, icon: 'i-tabler-news', color: 'sky' },
    { label: 'Tiket Pengaduan', value: values.complaints || 0, icon: 'i-tabler-message-2', color: 'amber' },
    { label: 'Total Tayangan Berita', value: values.news_views || 0, icon: 'i-tabler-eye', color: 'emerald' },
    { label: 'Aplikasi Aktif', value: values.active_apps || 0, icon: 'i-tabler-apps', color: 'violet' }
  ]
})

const aktivitasTerbaru = computed(() => dashboardData.value?.activities || [])
const pengaduanTerbaru = computed(() => dashboardData.value?.complaints || [])

const actionLabels = {
  create: 'Menambahkan',
  update: 'Mengubah',
  delete: 'Menghapus',
  login: 'Masuk ke sistem'
}

const actionText = action => actionLabels[action] || action

const getStatusColor = (status) => {
  if (status === 'Baru') return 'error'
  if (status === 'Diproses') return 'warning'
  if (status === 'Selesai') return 'success'
  return 'neutral'
}

const route = useRoute()
const toast = useToast()

watch(error, value => {
  if (value) toast.add({ title: 'Gagal memuat dashboard', description: 'Data dashboard tidak dapat diambil dari server.', color: 'error' })
})

// Cek query unauthorized
onMounted(() => {
  if (route.query.unauthorized === '1') {
    toast.add({
      title: 'Akses Ditolak',
      description: 'Anda tidak memiliki izin untuk mengakses halaman tersebut.',
      icon: 'i-tabler-alert-triangle',
      color: 'error',
      timeout: 5000
    })
    // Bersihkan query dari URL
    navigateTo('/admin', { replace: true })
  }
})
</script>

<template>
  <div class="space-y-6">

    <!-- Selamat Datang -->
    <UCard>
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Selamat Datang, {{ user?.nama || 'Admin' }}!
          </h2>
          <p class="text-slate-500 mt-1">Berikut ringkasan aktivitas portal hari ini.</p>
        </div>
      </div>
    </UCard>

    <!-- Statistik -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <UCard v-for="item in 4" :key="item">
        <USkeleton class="h-16 w-full" />
      </UCard>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <UCard v-for="stat in stats" :key="stat.label">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0"
            :class="`bg-${stat.color}-100 text-${stat.color}-600`">
            <UIcon :name="stat.icon" class="w-6 h-6" />
          </div>
          <div>
            <p class="text-2xl font-bold text-slate-900 dark:text-white">{{ stat.value }}</p>
            <p class="text-xs text-slate-500">{{ stat.label }}</p>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Grid 2 kolom -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

      <!-- Aktivitas Terbaru -->
      <UCard class="lg:col-span-2">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
              <UIcon name="i-tabler-history" class="text-sky-500" />
              Aktivitas Terbaru
            </h3>
            <UButton label="Lihat Semua" variant="link" size="xs" trailing-icon="i-tabler-arrow-right" />
          </div>
        </template>

        <div v-if="loading" class="space-y-3">
          <USkeleton v-for="item in 4" :key="item" class="h-14 w-full" />
        </div>
        <AdminEmptyState v-else-if="!aktivitasTerbaru.length" icon="i-tabler-history" title="Belum ada aktivitas"
          description="Aktivitas admin akan tampil di sini." />
        <div v-else class="space-y-3">
          <div v-for="item in aktivitasTerbaru" :key="item.id"
            class="flex items-start gap-3 p-3 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
            <div class="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-user" class="w-4 h-4 text-slate-500" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-slate-700 dark:text-slate-300">
                <span class="font-semibold">{{ item.user }}</span>
                {{ actionText(item.action) }}
                <span class="text-sky-600 dark:text-sky-400 font-medium">{{ item.target }}</span>
              </p>
              <p class="text-xs text-slate-400 mt-1">{{ item.time }}</p>
            </div>
          </div>
        </div>
      </UCard>

      <!-- Pengaduan Terbaru -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
              <UIcon name="i-tabler-message-2" class="text-amber-500" />
              Pengaduan Terbaru
            </h3>
            <UButton to="/admin/interaksi/tiket" label="Lihat" variant="link" size="xs"
              trailing-icon="i-tabler-arrow-right" />
          </div>
        </template>

        <div v-if="loading" class="space-y-3">
          <USkeleton v-for="item in 3" :key="item" class="h-20 w-full" />
        </div>
        <AdminEmptyState v-else-if="!pengaduanTerbaru.length" icon="i-tabler-message-off" title="Belum ada pengaduan"
          description="Pengaduan terbaru akan tampil di sini." />
        <div v-else class="space-y-3">
          <div v-for="item in pengaduanTerbaru" :key="item.id"
            class="p-3 rounded-lg border border-slate-200 dark:border-slate-700">
            <div class="flex items-start justify-between gap-2 mb-1">
              <p class="text-sm font-medium text-slate-900 dark:text-white line-clamp-2">{{ item.title }}</p>
            </div>
            <div class="flex items-center justify-between mt-2">
              <UBadge :color="getStatusColor(item.status)" variant="subtle" size="xs">
                {{ item.status }}
              </UBadge>
              <span class="text-xs text-slate-400">{{ item.date }}</span>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Aksi Cepat -->
    <UCard>
      <template #header>
        <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
          <UIcon name="i-tabler-bolt" class="text-sky-500" />
          Aksi Cepat
        </h3>
      </template>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <UButton to="/admin/cms/berita" icon="i-tabler-news" label="Tambah Berita" color="neutral" variant="soft" block
          size="lg" />
        <UButton to="/admin/cms/pengumuman" icon="i-tabler-speakerphone" label="Buat Pengumuman" color="neutral"
          variant="soft" block size="lg" />
        <UButton to="/admin/media" icon="i-tabler-photo" label="Upload Media" color="neutral" variant="soft" block
          size="lg" />
        <UButton to="/admin/pengguna/daftar" icon="i-tabler-user-plus" label="Tambah Pengguna" color="neutral"
          variant="soft" block size="lg" />
      </div>
    </UCard>

  </div>
</template>
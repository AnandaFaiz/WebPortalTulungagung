<script setup>
definePageMeta({
  layout: 'admin',
  title: 'Dashboard'
})

const stats = [
  { label: 'Total Berita', value: 42, icon: 'i-tabler-news', color: 'sky' },
  { label: 'Tiket Pengaduan', value: 18, icon: 'i-tabler-message-2', color: 'amber' },
  { label: 'Pengunjung Hari Ini', value: '1.234', icon: 'i-tabler-users', color: 'emerald' },
  { label: 'Aplikasi Aktif', value: 8, icon: 'i-tabler-apps', color: 'violet' }
]

const aktivitasTerbaru = [
  { waktu: '10:24', user: 'Admin Utama', aksi: 'Menambahkan berita baru', target: 'Festival Budaya 2026' },
  { waktu: '09:15', user: 'Operator OPD', aksi: 'Mengubah data OPD', target: 'Dinas Pendidikan' },
  { waktu: '08:40', user: 'Admin Utama', aksi: 'Menjawab pengaduan', target: 'Tiket #1234' },
  { waktu: '08:12', user: 'Editor', aksi: 'Menerbitkan pengumuman', target: 'Penerimaan CPNS' }
]

const pengaduanTerbaru = [
  { id: 1, judul: 'Jalan Rusak di Kecamatan Boyolangu', status: 'Baru', tanggal: '2026-10-10' },
  { id: 2, judul: 'Lampu Jalan Mati di Jl. A. Yani', status: 'Diproses', tanggal: '2026-10-09' },
  { id: 3, judul: 'Sampah Menumpuk di Pasar', status: 'Selesai', tanggal: '2026-10-08' }
]

const getStatusColor = (status) => {
  if (status === 'Baru') return 'error'
  if (status === 'Diproses') return 'warning'
  if (status === 'Selesai') return 'success'
  return 'neutral'
}
</script>

<template>
  <div class="space-y-6">

    <!-- Selamat Datang -->
    <UCard>
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Selamat Datang, Admin! 👋</h2>
          <p class="text-slate-500 mt-1">Berikut ringkasan aktivitas portal hari ini.</p>
        </div>
        <UButton icon="i-tabler-plus" label="Buat Berita Baru" color="primary" to="/admin/cms/berita" />
      </div>
    </UCard>

    <!-- Statistik -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
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

        <div class="space-y-3">
          <div v-for="(item, i) in aktivitasTerbaru" :key="i"
            class="flex items-start gap-3 p-3 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
            <div class="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-user" class="w-4 h-4 text-slate-500" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-slate-700 dark:text-slate-300">
                <span class="font-semibold">{{ item.user }}</span>
                {{ item.aksi }}
                <span class="text-sky-600 dark:text-sky-400 font-medium">{{ item.target }}</span>
              </p>
              <p class="text-xs text-slate-400 mt-1">{{ item.waktu }}</p>
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
            <UButton label="Lihat" variant="link" size="xs" trailing-icon="i-tabler-arrow-right" />
          </div>
        </template>

        <div class="space-y-3">
          <div v-for="item in pengaduanTerbaru" :key="item.id"
            class="p-3 rounded-lg border border-slate-200 dark:border-slate-700">
            <div class="flex items-start justify-between gap-2 mb-1">
              <p class="text-sm font-medium text-slate-900 dark:text-white line-clamp-2">{{ item.judul }}</p>
            </div>
            <div class="flex items-center justify-between mt-2">
              <UBadge :color="getStatusColor(item.status)" variant="subtle" size="xs">
                {{ item.status }}
              </UBadge>
              <span class="text-xs text-slate-400">{{ item.tanggal }}</span>
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
        <UButton icon="i-tabler-news" label="Tambah Berita" color="neutral" variant="soft" block size="lg" />
        <UButton icon="i-tabler-speakerphone" label="Buat Pengumuman" color="neutral" variant="soft" block size="lg" />
        <UButton icon="i-tabler-photo" label="Upload Media" color="neutral" variant="soft" block size="lg" />
        <UButton icon="i-tabler-user-plus" label="Tambah Pengguna" color="neutral" variant="soft" block size="lg" />
      </div>
    </UCard>

  </div>
</template>